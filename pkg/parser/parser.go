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

	return p.findFunctions(tree.RootNode())
}

func (p *Parser) findFunctions(node *sitter.Node) []*Candidate {
	candidates := []*Candidate{}
	// TODO use treesitter predicates https://github.com/smacker/go-tree-sitter/#predicates

	// walking through the AST to get all function declarations
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		candidate := &Candidate{
			Path:     p.path,
			Function: &Function{},
			Metrics:  &Metrics{},
		}

		slog.Info("parsing child", "type", child.Type())
		switch child.Type() {
		case "method_declaration":
			// TODO check go methods
			// TODO can this move to "class_declaration"?
			// find class, if there is one (eg. java)
			if child.Parent() != nil && child.Parent().Parent() != nil {
				candidate.Class = p.name(child.Parent().Parent())
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
			methods := p.findFunctions(child.ChildByFieldName("body"))
			candidates = append(candidates, methods...)

		// TODO figure out how to handle global stuff like this
		case "package_clause":
			candidate.Package = child.NamedChild(0).Content(p.sourceCode)

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
	params := p.findByType(node, "parameter_list")
	if params != nil {
		for i := 0; i < int(params.NamedChildCount()); i++ {
			param := params.NamedChild(i)
			f.Parameters = append(f.Parameters, &Parameter{
				Name: p.name(param),
				Type: param.ChildByFieldName("type").Content(p.sourceCode),
			},
			)
		}
	}

	return f
}

// searches inside a node for a child having the given type
func (p *Parser) findByType(node *sitter.Node, nodeType string) *sitter.Node {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == nodeType {
			return child
		}
	}
	return nil
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
