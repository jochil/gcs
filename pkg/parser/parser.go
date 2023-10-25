package parser

import (
	"context"
	"log/slog"
	"os"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// Parser encapsulates a parser for a given source code file
type Parser struct {
	*sitter.Parser
	path       string
	sourceCode []byte
}

func NewParser(path string, language *sitter.Language) *Parser {
	parser := &Parser{
		Parser: sitter.NewParser(),
		path:   path,
	}
	parser.SetLanguage(language)

	return parser
}

// Parse returns a list of candidates for a given source code file
func (p *Parser) Parse() []*Candidate {
	slog.Info("Start parsing", "file", p.path)

	var err error
	p.sourceCode, err = os.ReadFile(p.path)
	if err != nil {
		panic(err)
	}

	tree, err := p.ParseCtx(context.Background(), nil, p.sourceCode)
	if err != nil {
		panic(err)
	}

	root := tree.RootNode()
	packageName := p.findPackage(root)
	return p.findFunctions(root, packageName)
}

func (p *Parser) findFunctions(node *sitter.Node, packageName string) []*Candidate {
	// TODO use treesitter predicates https://github.com/smacker/go-tree-sitter/#predicates
	candidates := []*Candidate{}

	// walking through the AST to get all function declarations
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		candidate := &Candidate{
			Path:     p.path,
			Function: &Function{},
			Metrics:  &Metrics{},
			Package:  packageName,
		}

		slog.Info("parsing child", "type", child.Type())
		switch child.Type() {
		case "method_declaration":
			// TODO can this move to "class_declaration"?
			// find class, if there is one (eg. java)
			if child.Parent() != nil && child.Parent().Parent() != nil {
				candidate.Class = p.name(child.Parent().Parent())
			}

			// handle go receiver
			if child.NamedChild(0).Type() == "parameter_list" {
				// parameter_list -> parameter_declaration
				param := p.parseParameter(child.NamedChild(0).NamedChild(0))
				candidate.Class = strings.Replace(param.Type, "*", "", 1)
			}

			fallthrough
		case "function_declaration":
			// handle normal function declarations
			candidate.Function = p.function(child)
			// generate control flow graph
			body := child.ChildByFieldName("body")
			candidate.ControlFlowGraph = parseToCfg(body)

			candidate.Metrics.LinesOfCode = p.countLines(body)
			candidate.Code = child.Content(p.sourceCode)

		case "function_definition":
			declarator := child.ChildByFieldName("declarator")
			candidate.Function.Name = p.name(declarator)

		case "lexical_declaration":
			// get functions declared as variables
			declarator := child.NamedChild(0)
			value := declarator.ChildByFieldName("value")
			if value.Type() == "function" || value.Type() == "arrow_function" {
				candidate.Function.Name = p.name(declarator)
			}

		case "class_declaration":
			methods := p.findFunctions(child.ChildByFieldName("body"), packageName)
			candidates = append(candidates, methods...)

		case "package_clause", "package_declaration":
			// ignored types

		default:
			slog.Warn("not handled type", "type", child.Type())
		}

		if candidate.Function.Name != "" {
			candidates = append(candidates, candidate)

			if candidate.ControlFlowGraph != nil {
				cc, err := candidate.CalcCyclomaticComplexity()
				if err != nil {
					cc = -1
					slog.Warn("unable to calc cyclomatic complexity", "func", candidate.Function.Name)
				}
				candidate.Metrics.CyclomaticComplexity = cc
			}

			slog.Info("Found candidate", "function", candidate)
		}

	}
	return candidates
}

// initializes a Function struct from a given tree-sitter node
func (p *Parser) function(node *sitter.Node) *Function {
	f := &Function{
		Name:       p.name(node),
		Parameters: []*Parameter{},
	}

	// getting all the parameters
	param_lists := p.findByType(node, "parameter_list")
	var params *sitter.Node

	switch len(param_lists) {
	case 1:
		params = param_lists[0]
	case 2:
		params = param_lists[1]
	case 0:
		params = nil
	default:
		slog.Warn("more parameter_list nodes than expected", "function", f.Name)
	}

	if params != nil {
		for i := 0; i < int(params.NamedChildCount()); i++ {
			f.Parameters = append(f.Parameters, p.parseParameter(params.NamedChild(i)))
		}
	}

	return f
}

// searches inside a node for a child having the given type
func (p *Parser) findByType(node *sitter.Node, nodeType string) []*sitter.Node {
	nodes := []*sitter.Node{}
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == nodeType {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (p *Parser) findPackage(node *sitter.Node) string {
	packageDefs := p.findByType(node, "package_clause")
	if len(packageDefs) == 0 {
		packageDefs = p.findByType(node, "package_declaration")

	}

	if len(packageDefs) > 0 {
		// if there are more than one node log a warning and use the first one
		if len(packageDefs) > 1 {
			slog.Warn("found multiple package_clause|_declaration nodes")
		}
		// package_clause -> package_identifier
		// package_declaration -> scoped_identifier
		return packageDefs[0].NamedChild(0).Content(p.sourceCode)
	}

	return ""
}

func (p *Parser) parseParameter(param *sitter.Node) *Parameter {
	return &Parameter{
		Name: p.name(param),
		Type: param.ChildByFieldName("type").Content(p.sourceCode),
	}
}

// returns the name/identifier of a tree-sitter node (eg. function/variable name)
func (p *Parser) name(node *sitter.Node) string {
	child := node.ChildByFieldName("name")
	// sometimes the function name is stored in the declarator field
	// for example in the "function_definition" type
	if child == nil {
		child = node.ChildByFieldName("declarator")
	}
	if child == nil {
		slog.Warn("unable to get name", "type", node.Type())
		return "???"
	}

	return child.Content(p.sourceCode)
}

func (p *Parser) countLines(node *sitter.Node) int {
	// TODO count actual lines.. no comments, no empty ones, ...
	code := node.Content(p.sourceCode)
	lines := strings.Split(strings.ReplaceAll(code, "\r\n", "\n"), "\n")
	return len(lines)
}
