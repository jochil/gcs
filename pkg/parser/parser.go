package parser

import (
	"context"
	"fmt"
	"os"

	"github.com/jochil/test-helper/pkg/data"
	sitter "github.com/smacker/go-tree-sitter"
)

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

func (p *Parser) Parse() []*data.Candidate {
	fmt.Println("Parsing:", p.path)

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

func (p *Parser) findFunctions(node *sitter.Node) []*data.Candidate {
	candidates := []*data.Candidate{}
	// TODO use treesitter predicates https://github.com/smacker/go-tree-sitter/#predicates

	// walking through the AST to get all function declarations
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		candidate := &data.Candidate{
			Path:     p.path,
			Function: &data.Function{},
		}

		switch child.Type() {
		case "method_declaration":
			// TODO can this move to "class_declaration"?
			candidate.Class = p.name(child.Parent().Parent())
			fallthrough
		case "function_declaration":
			// handle normal function declarations
			candidate.Function = p.function(child)

			// generate control flow graph
			// TODO: move this somewhere else?
			SaveGraph(candidate.Function.Name, ParseToCfg(child.ChildByFieldName("body")))

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
			fmt.Println("not handled type:", child.Type())
		}

		if candidate.Function.Name != "" {
			candidates = append(candidates, candidate)
			fmt.Println("\t Found candidate:", candidate)
		}

	}
	return candidates
}

func (p *Parser) function(node *sitter.Node) *data.Function {
	f := &data.Function{
		Name:       p.name(node),
		Parameters: []*data.Parameter{},
	}

	params := p.findByType(node, "parameter_list")
	if params != nil {
		for i := 0; i < int(params.NamedChildCount()); i++ {
			param := params.NamedChild(i)
			f.Parameters = append(f.Parameters, &data.Parameter{
				Name: p.name(param),
				Type: param.ChildByFieldName("type").Content(p.sourceCode),
			},
			)
		}
	}

	return f
}

func (p *Parser) findByType(node *sitter.Node, nodeType string) *sitter.Node {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == nodeType {
			return child
		}
	}
	return nil
}

func (p *Parser) name(node *sitter.Node) string {
	child := node.ChildByFieldName("name")
	// sometimes the function name is stored in the declarator field
	// for example in the "function_definition" type
	if child == nil {
		child = node.ChildByFieldName("declarator")
	}
	return child.Content(p.sourceCode)
}
