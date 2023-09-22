package parser

import (
	"fmt"
	"log/slog"

	"github.com/dominikbraun/graph"
	sitter "github.com/smacker/go-tree-sitter"
)

type cfgParser struct {
	g        graph.Graph[int, int]
	counter  int
	startRef int
	endRef   int
}

func parseToCfg(node *sitter.Node) graph.Graph[int, int] {
	cp := &cfgParser{
		g:       graph.New(graph.IntHash, graph.Directed()),
		counter: -1,
	}

	cp.startRef = cp.addVertex("start", "lightgreen")
	cp.endRef = cp.addVertex("end", "crimson")

	prevRef := cp.blockToGraph(node, cp.startRef)

	cp.addEdge(prevRef, cp.endRef)

	return cp.g
}

func (cp *cfgParser) nodeToGraph(node *sitter.Node, prevRef int) int {
	if node == nil {
		return prevRef
	}
	switch node.Type() {
	case "if_statement":
		return cp.ifToGraph(node, prevRef)
	case "block":
		return cp.blockToGraph(node, prevRef)
	case "return_statement":
		return cp.returnToGraph(node, prevRef)
	default:
		slog.Info("graph: unknown node type", "type", node.Type())
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

func (cp *cfgParser) returnToGraph(node *sitter.Node, prevRef int) int {
	ref := cp.addVertex("return", "red")
	cp.addEdge(prevRef, ref)
	return ref
}

func (cp *cfgParser) unknownToGraph(node *sitter.Node, prevRef int) int {
	ref := cp.addVertex(node.Type(), "azure")
	cp.addEdge(prevRef, ref)
	return ref
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
