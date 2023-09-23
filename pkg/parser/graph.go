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
	case "expression_switch_statement":
		return cp.switchToGraph(node, prevRef)
	case "for_statement":
		return cp.forToGraph(node, prevRef)
	case "block":
		return cp.blockToGraph(node, prevRef)
	case "return_statement":
		return cp.returnToGraph(node, prevRef)

	// ignore these nodes
	case "expression_list":
		return prevRef
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

func (cp *cfgParser) forToGraph(forStatement *sitter.Node, prevRef int) int {
	forStartRef := cp.addVertex("for_start", "cyan")
	cp.addEdge(prevRef, forStartRef)

	forEndRef := cp.addVertex("for_end", "cyan3")
	cp.addEdge(forEndRef, forStartRef)

	blockRef := cp.blockToGraph(forStatement.ChildByFieldName("body"), forStartRef)
	cp.addEdge(blockRef, forEndRef)

	return forEndRef
}

func (cp *cfgParser) switchToGraph(switchStatement *sitter.Node, prevRef int) int {
	switchStartRef := cp.addVertex("switch_start", "cyan")
	cp.addEdge(prevRef, switchStartRef)

	switchEndRef := cp.addVertex("switch_end", "cyan3")

	defaultCase := false

	for i := 0; i < int(switchStatement.NamedChildCount()); i++ {
		child := switchStatement.NamedChild(i)
		if child.Type() == "expression_case" || child.Type() == "default_case" {
			caseRef := cp.blockToGraph(child, switchStartRef)
			cp.addEdge(caseRef, switchEndRef)
		}

		if child.Type() == "default_case" {
			defaultCase = true
		}
	}

	if !defaultCase {
		cp.addEdge(switchStartRef, switchEndRef)
	}

	return switchEndRef
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
	err := cp.g.AddEdge(start, end)
	if err != nil {
		slog.Warn("unable to add edge to graph", "start", start, "end", end)
	}
}

func (cp *cfgParser) addVertex(label string, color string) int {
	cp.counter++
	err := cp.g.AddVertex(cp.counter, graph.VertexAttributes(map[string]string{
		"label":     fmt.Sprintf("%d: %s", cp.counter, label),
		"style":     "filled, solid",
		"color":     "black",
		"fillcolor": color,
	}))
	if err != nil {
		slog.Warn("unable to add node to graph", "label", label)
	}
	return cp.counter
}
