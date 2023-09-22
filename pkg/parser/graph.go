package parser

import (
	"fmt"

	"github.com/dominikbraun/graph"
	sitter "github.com/smacker/go-tree-sitter"
)

type cfgParser struct {
	g       graph.Graph[int, int]
	counter int
}

func parseToCfg(node *sitter.Node) graph.Graph[int, int] {
	cp := &cfgParser{
		g:       graph.New(graph.IntHash, graph.Directed()),
		counter: -1,
	}

	startRef := cp.addVertex("start", "lightgreen")

	prevRef := cp.blockToGraph(node, startRef)

	endRef := cp.addVertex("end", "crimson")
	cp.addEdge(prevRef, endRef)

	return cp.g

}

func (cp *cfgParser) addEdge(start, end int) {
	cp.g.AddEdge(start, end)
}

func (cp *cfgParser) addVertex(label string, color string) int {
	cp.counter++
	cp.g.AddVertex(cp.counter, graph.VertexAttributes(map[string]string{
		"label":     fmt.Sprintf("%d: %s", cp.counter, label),
		"style":     "filled, solid",
		"color":     "black",
		"fillcolor": color,
	}))
	return cp.counter
}

func (cp *cfgParser) nodeToGraph(node *sitter.Node, prevRef int) int {
	switch node.Type() {
	case "if_statement":
		return cp.ifToGraph(node, prevRef)
	case "block":
		return cp.blockToGraph(node, prevRef)
	default:
		return cp.unknownToGraph(node, prevRef)
	}
}

func (cp *cfgParser) blockToGraph(block *sitter.Node, prevRef int) int {
	for i := 0; i < int(block.NamedChildCount()); i++ {
		child := block.NamedChild(i)
		prevRef = cp.nodeToGraph(child, prevRef)
	}
	return prevRef
}

func (cp *cfgParser) ifToGraph(ifStatement *sitter.Node, prevRef int) int {
	// create node for "if" start
	ifStartRef := cp.addVertex("if_start", "cyan")
	cp.addEdge(prevRef, ifStartRef)

	// add node to end if
	ifEndRef := cp.addVertex("if_end", "cyan3")

	// parse the "if" path
	prevRef = cp.nodeToGraph(ifStatement.ChildByFieldName("consequence"), ifStartRef)
	cp.addEdge(prevRef, ifEndRef)
	prevRef = cp.nodeToGraph(ifStatement.ChildByFieldName("alternative"), ifStartRef)
	cp.addEdge(prevRef, ifEndRef)

	return ifEndRef
}

func (cp *cfgParser) unknownToGraph(node *sitter.Node, prevRef int) int {
	ref := cp.addVertex(node.Type(), "azure")
	cp.addEdge(prevRef, ref)
	return ref
}
