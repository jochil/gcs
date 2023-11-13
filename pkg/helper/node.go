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
	fmt.Println(node.String())
	printNode(node, 0)
}

func printNode(node *sitter.Node, depth int) {
	ident := "    "
	fmt.Printf("%s%s\n", strings.Repeat(ident, depth), node.Type())

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if !child.IsNamed() {
			continue
		}
		if fieldName := node.FieldNameForChild(i); fieldName != "" {
			fmt.Printf("\n%s:%s:\n", strings.Repeat(ident, depth+1), fieldName)
		}
		printNode(child, depth+1)
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
