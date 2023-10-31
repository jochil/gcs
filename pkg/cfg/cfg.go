package cfg

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

// generates a control flow graph based on a tree-sitter node (usually a function body)
func Create(node *sitter.Node) graph.Graph[int, int] {
	cp := &cfgParser{
		g:       graph.New(graph.IntHash, graph.Directed()),
		counter: -1,
	}

	// start and endpoint
	cp.startRef = cp.addVertex("start", "lightgreen")
	cp.endRef = cp.addVertex("end", "crimson")

	// handle the (function) body
	prevRef := cp.blockToGraph(node, cp.startRef)

	cp.addEdge(prevRef, cp.endRef)

	return cp.g
}

// handles a single node
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

// iterates over all childs of a given block (eg. function body, if/else body, ...)
func (cp *cfgParser) blockToGraph(block *sitter.Node, prevRef int) int {
	for i := 0; i < int(block.NamedChildCount()); i++ {
		child := block.NamedChild(i)
		prevRef = cp.nodeToGraph(child, prevRef)
	}
	return prevRef
}

// parses a for loop into the cfg
func (cp *cfgParser) forToGraph(forStatement *sitter.Node, prevRef int) int {
	// create start node
	forStartRef := cp.addVertex("for_start", "cyan")
	cp.addEdge(prevRef, forStartRef)

	// create end node and connect it with the start node
	forEndRef := cp.addVertex("for_end", "cyan3")
	cp.addEdge(forEndRef, forStartRef)

	blockRef := cp.blockToGraph(forStatement.ChildByFieldName("body"), forStartRef)

	// connect the last node of the block with the end node
	cp.addEdge(blockRef, forEndRef)

	return forEndRef
}

// parses switch statement into the cfg
func (cp *cfgParser) switchToGraph(switchStatement *sitter.Node, prevRef int) int {
	// create start node and connect it with the previous one
	switchStartRef := cp.addVertex("switch_start", "cyan")
	cp.addEdge(prevRef, switchStartRef)

	// create end node
	switchEndRef := cp.addVertex("switch_end", "cyan3")

	defaultCase := false

	// iterate over the different cases
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

	// if there is no default case, connect the start node with the end node
	if !defaultCase {
		cp.addEdge(switchStartRef, switchEndRef)
	}

	return switchEndRef
}

// parses if/elseif/else nodes into the cfg
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

// handle return statement
func (cp *cfgParser) returnToGraph(node *sitter.Node, prevRef int) int {
	ref := cp.addVertex("return", "red")
	cp.addEdge(prevRef, ref)
	return ref
}

// handles unknown nodes
func (cp *cfgParser) unknownToGraph(node *sitter.Node, prevRef int) int {
	ref := cp.addVertex(node.Type(), "azure")
	cp.addEdge(prevRef, ref)
	return ref
}

// wrapper for adding edges
func (cp *cfgParser) addEdge(start, end int) {
	err := cp.g.AddEdge(start, end)
	if err != nil {
		slog.Warn("unable to add edge to graph", "start", start, "end", end)
	}
}

// wrapper for adding nodes
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
