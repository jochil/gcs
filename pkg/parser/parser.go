package parser

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"strings"
	"unicode"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/helper"
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
	parser.SetLanguage(helper.SitterLanguages[language])

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
			Package:  packageName,
			Language: p.language,
		}

		slog.Info("parsing child", "type", child.Type())
		switch child.Type() {

		case "function_definition":
			declarator := child.ChildByFieldName("declarator")
			c.Function.Name = p.name(declarator)
			p.parseFunction(child, c)

		case "function_declaration", "method_declaration", "method_definition":
			p.parseFunction(child, c)

		case "lexical_declaration":
			// get functions declared as variables
			declarator := child.NamedChild(0)
			functionNode := helper.FirstChildByTypes(declarator, []string{"function", "arrow_function"})
			if functionNode != nil {
				c.Function.Name = p.name(declarator)
				p.parseFunction(functionNode, c)
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
	p.parseSignature(node, c)
	p.parseVisibility(node, c.Function)
}

func (p *Parser) parseSignature(node *sitter.Node, c *candidate.Candidate) {
	f := c.Function
	if f.Name == "" {
		f.Name = p.name(node.ChildByFieldName("name"))
	}

	returnFieldName := "return_type"
	switch p.language {
	case types.Java:
		returnFieldName = "type"
	case types.Go:
		returnFieldName = "result"
		if receiver := node.ChildByFieldName("receiver"); receiver != nil {
			c.Class = p.parseParameters(receiver)[0].Type
		}
	}

	f.Parameters = p.parseParameters(node.ChildByFieldName("parameters"))
	f.ReturnValues = p.parseParameters(node.ChildByFieldName(returnFieldName))

}

func (p *Parser) parseParameters(node *sitter.Node) []*candidate.Parameter {
	params := []*candidate.Parameter{}
	if node == nil {
		return params
	}

	switch node.Type() {
	case "parameter_list", "formal_parameters":
		for i := 0; i < int(node.NamedChildCount()); i++ {
			child := node.NamedChild(i)
			params = append(params, p.parseParameter(child))
		}

	case "void_type":
		// do nothing

	default:
		params = append(params, p.parseParameter(node))
	}

	return params
}

func (p *Parser) parseParameter(param *sitter.Node) *candidate.Parameter {
	var typeName string
	var name string

	switch param.Type() {
	case "spread_parameter":
		name = p.name(param.NamedChild(1))
		typeName = p.typeName(param)

	case "required_parameter", "optional_parameter":
		name = p.name(param.NamedChild(0))
		typeName = p.typeName(param.ChildByFieldName("type"))

	case "parameter_declaration", "formal_parameter":
		name = p.name(param.ChildByFieldName("name"))
		typeName = p.typeName(param.ChildByFieldName("type"))

	default:
		name = types.NoName
		typeName = p.typeName(param)
	}

	return &candidate.Parameter{
		Name: name,
		Type: typeName,
	}
}

func (p *Parser) parseVisibility(node *sitter.Node, f *candidate.Function) {
	// setting default visibility
	f.Visibility = types.VisibilityPublic

	// TODO handle more than 1 modifiers node... is this even possible?
	if p.language == types.Java || p.language == types.TypeScript || p.language == types.JavaScript {

		mod := helper.FirstChildByTypes(node, []string{"modifiers", "accessibility_modifier"})
		if mod != nil {
			modifiers := mod.Content(p.sourceCode)
			if strings.Contains(modifiers, "private") {
				f.Visibility = types.VisibilityPrivate
			} else if strings.Contains(modifiers, "protected") {
				f.Visibility = types.VisibilityProtected
			}

			if strings.Contains(modifiers, "static") {
				f.Static = true
			}
		}

		// handle js private fields
		// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes/Private_class_fields
		if child := helper.FirstChildByType(node, "private_property_identifier"); child != nil {
			f.Visibility = types.VisibilityPrivate
		}

	}

	if p.language == types.Go {
		runes := []rune(f.Name)
		if unicode.IsUpper(runes[0]) {
			f.Visibility = types.VisibilityPublic
		} else {
			f.Visibility = types.VisibilityPrivate
		}
	}
}

func (p *Parser) findPackage(node *sitter.Node) string {
	packageDefs := helper.ChildrenByType(node, "package_clause")
	if len(packageDefs) == 0 {
		packageDefs = helper.ChildrenByType(node, "package_declaration")

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

func (p *Parser) typeName(node *sitter.Node) string {
	switch node.Type() {
	case "type_annotation":
		return node.NamedChild(0).Content(p.sourceCode)
	case "integral_type",
		"floating_point_type",
		"boolean_type",
		"array_type",
		"generic_type",
		"predefined_type",
		"union_type":
		return node.Content(p.sourceCode)

	case "spread_parameter":
		return node.NamedChild(0).Content(p.sourceCode) + "..."
	default:
		return p.name(node)
	}
}

// returns the name/identifier of a tree-sitter node (eg. function/variable name)
func (p *Parser) name(node *sitter.Node) string {

	if node == nil {
		return types.NoName
	}

	nameTypes := []string{
		"name",
		"identifier",
		"field_identifier",
		"property_identifier",
		"type_identifier",
		"private_property_identifier",
		"type_annotation",
	}

	if slices.Contains(nameTypes, node.Type()) {
		name := node.Content(p.sourceCode)
		if node.Type() == "private_property_identifier" {
			// handle js private fields
			// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes/Private_class_fields
			name = strings.TrimPrefix(name, "#")
		}
		return name
	}

	// try different types that can contain the name
	child := helper.FirstChildByTypes(node, nameTypes)
	if child == nil {
		slog.Warn("unable to get name", "type", node.Type())
		return types.NoName
	}

	return child.Content(p.sourceCode)
}
