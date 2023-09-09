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
			Path: p.path,
		}

		switch child.Type() {
		case "method_declaration":
			candidate.Class = p.name(child.Parent().Parent())
			fallthrough
		case "function_declaration":
			// handle normal function declarations
			candidate.Function = p.name(child)
			SaveGraph(candidate.Function, ParseToCfg(child.ChildByFieldName("body")))

		case "function_definition":
			declarator := child.ChildByFieldName("declarator")
			candidate.Function = p.name(declarator)

		case "lexical_declaration":
			// get functions declared as variables
			declarator := child.NamedChild(0)
			value := declarator.ChildByFieldName("value")
			if value.Type() == "function" || value.Type() == "arrow_function" {
				candidate.Function = p.name(declarator)
			}

		case "class_declaration":
			methods := p.findFunctions(child.ChildByFieldName("body"))
			candidates = append(candidates, methods...)

		default:
			fmt.Println("not handled type:", child.Type())
		}

		if candidate.Function != "" {
			candidates = append(candidates, candidate)
			fmt.Println("\t Found candidate:", candidate)
		}

	}
	return candidates
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
