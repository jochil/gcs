package parser

import (
	"context"
	"os"

	"github.com/jochil/test-helper/pkg/data"
	sitter "github.com/smacker/go-tree-sitter"
)

func Parse(path string, language *sitter.Language) []*data.Candidate {
	parser := sitter.NewParser()
	parser.SetLanguage(language)

	sourceCode, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
	if err != nil {
		panic(err)
	}

	return findFunctions(tree.RootNode(), path, sourceCode)
}

func findFunctions(node *sitter.Node, path string, sourceCode []byte) []*data.Candidate {
	candidates := []*data.Candidate{}
	// TODO use treesitter predicates https://github.com/smacker/go-tree-sitter/#predicates

	// walking through the AST to get all function declarations
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)

		candidate := &data.Candidate{
			Path: path,
		}
		// handle normal function declarations
		if child.Type() == "function_declaration" {
			candidate.Name = name(child, sourceCode)
		} else if child.Type() == "lexical_declaration" {
			// get functions declared as variables
			declarator := child.NamedChild(0)
			value := declarator.ChildByFieldName("value")
			if value.Type() == "function" || value.Type() == "arrow_function" {
				candidate.Name = name(declarator, sourceCode)
			}
		}

		if candidate.Name != "" {
			candidates = append(candidates, candidate)
		}

	}
	return candidates
}

func name(node *sitter.Node, sourceCode []byte) string {
	return node.ChildByFieldName("name").Content(sourceCode)
}
