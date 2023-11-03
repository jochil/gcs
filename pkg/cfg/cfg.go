package cfg

import (
	"fmt"
	"log/slog"

	"github.com/dominikbraun/graph"
	"github.com/jochil/dlth/pkg/helper"
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
	case "switch_expression":
		switchBlock := helper.FirstChildByType(node, "switch_block")
		return cp.switchToGraph(switchBlock, prevRef)
	case "expression_switch_statement":
		return cp.switchToGraph(node, prevRef)
	case "do_statement":
		return cp.doToGraph(node, prevRef)
	case "while_statement":
		return cp.whileToGraph(node, prevRef)
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
		switch child.Type() {
		case "switch_label", "break_statement":
			// ignore these labels
		default:
			prevRef = cp.nodeToGraph(child, prevRef)
		}
	}
	return prevRef
}

func (cp *cfgParser) doToGraph(doStatement *sitter.Node, prevRef int) int {
	startRef := cp.addVertex("do_start", "cyan")
	cp.addEdge(prevRef, startRef)

	// create end node and connect it with the start node
	endRef := cp.addVertex("do_end", "cyan3")

	blockRef := cp.blockToGraph(doStatement.ChildByFieldName("body"), startRef)

	// connect the last node of the block with the start and end nodes
	cp.addEdge(blockRef, startRef)
	cp.addEdge(blockRef, endRef)

	return endRef
}

func (cp *cfgParser) whileToGraph(whileStatement *sitter.Node, prevRef int) int {
	startRef := cp.addVertex("while_start", "cyan")
	cp.addEdge(prevRef, startRef)

	// create end node and connect it with the start node
	endRef := cp.addVertex("while_end", "cyan3")
	cp.addEdge(startRef, endRef)

	blockRef := cp.blockToGraph(whileStatement.ChildByFieldName("body"), startRef)

	// connect the last node of the block with the start node
	cp.addEdge(blockRef, startRef)

	return endRef
}

// parses a for loop into the cfg
func (cp *cfgParser) forToGraph(forStatement *sitter.Node, prevRef int) int {
	// create start node
	startRef := cp.addVertex("for_start", "cyan")
	cp.addEdge(prevRef, startRef)

	// create end node and connect it with the start node
	endRef := cp.addVertex("for_end", "cyan3")
	cp.addEdge(endRef, startRef)

	blockRef := cp.blockToGraph(forStatement.ChildByFieldName("body"), startRef)

	// connect the last node of the block with the end node
	cp.addEdge(blockRef, endRef)

	return endRef
}

// parses switch statement into the cfg
func (cp *cfgParser) switchToGraph(switchStatement *sitter.Node, prevRef int) int {
	// create start node and connect it with the previous one
	startRef := cp.addVertex("switch_start", "cyan")
	cp.addEdge(prevRef, startRef)

	// create end node
	endRef := cp.addVertex("switch_end", "cyan3")

	defaultCase := false

	// iterate over the different cases
	for i := 0; i < int(switchStatement.NamedChildCount()); i++ {
		child := switchStatement.NamedChild(i)
		switch child.Type() {
		case "switch_block_statement_group":
			// TODO check explicitly for "default"
			// java: if the switch label has no child it has to be the default case
			if helper.FirstChildByType(child, "switch_label").NamedChildCount() == 0 {
				defaultCase = true
			}
			caseRef := cp.blockToGraph(child, startRef)
			cp.addEdge(caseRef, endRef)
		case "default_case":
			defaultCase = true
			fallthrough
		case "expression_case":
			caseRef := cp.blockToGraph(child, startRef)
			cp.addEdge(caseRef, endRef)
		}

	}

	// if there is no default case, connect the start node with the end node
	if !defaultCase {
		cp.addEdge(startRef, endRef)
	}

	return endRef
}

// parses if/elseif/else nodes into the cfg
func (cp *cfgParser) ifToGraph(ifStatement *sitter.Node, prevRef int) int {
	// create node for "if" start
	startRef := cp.addVertex("if_start", "cyan")
	cp.addEdge(prevRef, startRef)

	// add node to end if
	endRef := cp.addVertex("if_end", "cyan3")

	// parse the "if" path
	prevRef = cp.nodeToGraph(ifStatement.ChildByFieldName("consequence"), startRef)
	cp.addEdge(prevRef, endRef)
	prevRef = cp.nodeToGraph(ifStatement.ChildByFieldName("alternative"), startRef)
	cp.addEdge(prevRef, endRef)

	return endRef
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
