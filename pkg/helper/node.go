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
	s := node.String()
	ident := "    "
	level := 0
	word := ""
	for _, c := range s {
		if string(c) == " " || string(c) == ")" || string(c) == "(" {
			if word != "" {
				// mark field names with F:
				if strings.HasSuffix(word, ":") {
					word = "F: " + word[:len(word)-1]
				}
				fmt.Printf("%s%s\n", strings.Repeat(ident, level), word)
				word = ""
			}

			if string(c) == "(" {
				level++
				// First element after ( is a type, mark it with T:
				word = "T: " + word
			} else if string(c) == ")" {
				level--
			}
		} else {
			word += string(c)
		}

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
