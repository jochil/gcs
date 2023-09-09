package parser

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func print(node *sitter.Node, ident int) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", ident), node.Type())

	for i := 0; i < int(node.NamedChildCount()); i++ {
		print(node.NamedChild(i), ident+1)
	}
}
