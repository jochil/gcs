package helper

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// Prints a tree-sitter node nicely
//
//nolint:unused
func PrintNode(node *sitter.Node) {
	printNode(node, 0)
}

func printNode(node *sitter.Node, ident int) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", ident), node.Type())

	for i := 0; i < int(node.NamedChildCount()); i++ {
		printNode(node.NamedChild(i), ident+1)
	}
}

func FirstChildByType(node *sitter.Node, typeName string) *sitter.Node {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == typeName {
			return child
		}
	}
	return nil
}

func FirstChildByTypes(node *sitter.Node, typeNames []string) *sitter.Node {
	for _, typeName := range typeNames {
		child := FirstChildByType(node, typeName)
		if child != nil {
			return child
		}
	}
	return nil
}

// searches inside a node for a child having the given type
func ChildrenByType(node *sitter.Node, nodeType string) []*sitter.Node {
	nodes := []*sitter.Node{}
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(i)
		if child.Type() == nodeType {
			nodes = append(nodes, child)
		}
	}
	return nodes
}
