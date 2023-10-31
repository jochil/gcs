package parser

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"unicode"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/types"
	sitter "github.com/smacker/go-tree-sitter"
)

// Parser encapsulates a parser for a given source code file
type Parser struct {
	*sitter.Parser
	path       string
	sourceCode []byte
	language   types.Language
}

func NewParser(path string, language types.Language) *Parser {
	parser := &Parser{
		Parser:   sitter.NewParser(),
		path:     path,
		language: language,
	}
	parser.SetLanguage(sitterLanguages[language])

	return parser
}

// Parse returns a list of candidates for a given source code file
func (p *Parser) Parse() []*candidate.Candidate {
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

func (p *Parser) findFunctions(node *sitter.Node, packageName string) []*candidate.Candidate {
	// TODO use treesitter predicates https://github.com/smacker/go-tree-sitter/#predicates
	candidates := []*candidate.Candidate{}

	// walking through the AST to get all function declarations
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		c := &candidate.Candidate{
			Path:     p.path,
			Function: &candidate.Function{},
			Metrics:  &candidate.Metrics{},
			Package:  packageName,
		}

		slog.Info("parsing child", "type", child.Type())
		switch child.Type() {
		case "function_declaration", "method_declaration":
			p.parseFunction(child, c)

		case "function_definition":
			declarator := child.ChildByFieldName("declarator")
			c.Function.Name = p.name(declarator)

		case "lexical_declaration":
			// get functions declared as variables
			declarator := child.NamedChild(0)
			value := declarator.ChildByFieldName("value")
			if value.Type() == "function" || value.Type() == "arrow_function" {
				c.Function.Name = p.name(declarator)
			}

		case "class_declaration":
			methods := p.findFunctions(child.ChildByFieldName("body"), packageName)
			candidates = append(candidates, methods...)
			for _, c := range candidates {
				c.Class = p.name(child)
			}

		case "package_clause", "package_declaration":
			// ignored types

		default:
			slog.Warn("not handled type", "type", child.Type())
		}

		if c.Function.Name != "" {

			c.AST = child
			c.Code = child.Content(p.sourceCode)

			slog.Info("Found candidate", "function", c)
			candidates = append(candidates, c)
		}

	}
	return candidates
}

// initializes a Function struct from a given tree-sitter node
func (p *Parser) parseFunction(node *sitter.Node, c *candidate.Candidate) {
	f := &candidate.Function{
		Name:         p.name(node),
		Parameters:   []*candidate.Parameter{},
		ReturnValues: []*candidate.Parameter{},
	}

	p.parseVisibility(node, f)

	// getting all the parameter_list nodes
	paramLists := p.findByType(node, "parameter_list")
	if len(paramLists) == 0 {
		// used by java
		paramLists = p.findByType(node, "formal_parameters")
	}

	goReceiverType := func(paramList *sitter.Node) string {
		goReceiverParams := p.parseParameters(paramList)
		return goReceiverParams[0].Type
	}

	switch len(paramLists) {
	case 1:
		f.Parameters = p.parseParameters(paramLists[0])
	case 2:
		// handle golang case with a receiver and no/one unnamed return value
		if node.Type() == "method_declaration" && node.NamedChild(0).Type() == "parameter_list" {
			c.Class = goReceiverType(paramLists[0])
			f.Parameters = p.parseParameters(paramLists[1])
		} else {
			f.Parameters = p.parseParameters(paramLists[0])
			f.ReturnValues = p.parseParameters(paramLists[1])
		}
	case 3:
		// three param lists has to be a go method with a receiver and multiple return values
		c.Class = goReceiverType(paramLists[0])
		f.Parameters = p.parseParameters(paramLists[1])
		f.ReturnValues = p.parseParameters(paramLists[2])
	case 0:
	default:
		slog.Warn("more parameter_list nodes than expected", "function", f.Name)
	}

	p.parseReturnType(node, f)

	c.Function = f
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

func (p *Parser) parseReturnType(node *sitter.Node, f *candidate.Function) {
	knownTypes := []string{
		"type_identifier",
		"integral_type",
		"floating_point_type",
		"boolean_type",
		"array_type",
	}
	for _, t := range knownTypes {
		nodes := p.findByType(node, t)
		if len(nodes) == 1 {
			f.ReturnValues = []*candidate.Parameter{{Name: NoName, Type: nodes[0].Content(p.sourceCode)}}
			return
		}
	}
}

func (p *Parser) parseVisibility(node *sitter.Node, f *candidate.Function) {
	// TODO handle more than 1 modifiers node... is this even possible?
	if p.language == types.Java {
		// setting default visibility
		f.Visibility = VisibilityPublic

		if mods := p.findByType(node, "modifiers"); len(mods) >= 1 {
			modifiers := mods[0].Content(p.sourceCode)
			if strings.Contains(modifiers, "private") {
				f.Visibility = VisibilityPrivate
			} else if strings.Contains(modifiers, "protected") {
				f.Visibility = VisibilityProtected
			}

			if strings.Contains(modifiers, "static") {
				f.Static = true
			}
		}
	}

	if p.language == types.Go {
		runes := []rune(f.Name)
		if unicode.IsUpper(runes[0]) {
			f.Visibility = VisibilityPublic
		} else {
			f.Visibility = VisibilityPrivate
		}
	}
}

func (p *Parser) parseParameters(node *sitter.Node) []*candidate.Parameter {
	params := []*candidate.Parameter{}
	for i := 0; i < int(node.NamedChildCount()); i++ {
		params = append(params, p.parseParameter(node.NamedChild(i)))
	}
	return params
}

func (p *Parser) parseParameter(param *sitter.Node) *candidate.Parameter {
	return &candidate.Parameter{
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
		return NoName
	}

	return child.Content(p.sourceCode)
}
