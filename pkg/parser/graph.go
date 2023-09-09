package parser

import (
	"fmt"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	sitter "github.com/smacker/go-tree-sitter"
)

func (p *Parser) parseFunctionBody(name string, body *sitter.Node) {
	g := graph.New(graph.IntHash, graph.Directed())

	counter := 0
	g.AddVertex(counter, graph.VertexAttributes(map[string]string{
		"label":"start",
		"color":"lightgreen",
		"style": "filled",
	}))
	counter++

	prevRef, counter := p.blockToGraph(g, body, 0, counter)

	counter++
	g.AddVertex(counter, graph.VertexAttributes(map[string]string{
		"label":"end",
		"color":"crimson",
		"style": "filled",
	}))
	g.AddEdge(prevRef, counter)

	file, _ := os.Create(fmt.Sprintf(".draw/%s.gv", name))
	err := draw.DOT(g, file)
	if err != nil {
		panic(err)
	}
}

func (p *Parser) nodeToGraph(g graph.Graph[int, int], node *sitter.Node, prevRef, counter int) (int, int) {
	switch node.Type() {
	case "if_statement":
		return p.ifToGraph(g, node, prevRef, counter)
	case "block":
		return p.blockToGraph(g, node, prevRef, counter)
	default:
		return p.unknownToGraph(g, node, prevRef, counter)
	}
}

func (p *Parser) blockToGraph(g graph.Graph[int, int], block *sitter.Node, prevRef, counter int) (int, int) {
	for i := 0; i < int(block.NamedChildCount()); i++ {
		child := block.NamedChild(i)
		prevRef, counter = p.nodeToGraph(g, child, prevRef, counter)
	}
	return prevRef, counter
}

func (p *Parser) ifToGraph(g graph.Graph[int, int], ifStatement *sitter.Node, prevRef, counter int) (int, int) {
	// create node for "if" start
	counter++
	ifStartRef := counter
	label := fmt.Sprintf("%d %s", counter, "if_start")
	g.AddVertex(counter, graph.VertexAttributes(map[string]string{
		"label": label,
		"color": "cyan3",
		"style": "filled",
	}))
	g.AddEdge(prevRef, ifStartRef)

	// add node to end if
	counter++
	ifEndRef := counter
	label = fmt.Sprintf("%d %s", counter, "if_end")
	g.AddVertex(counter, graph.VertexAttributes(map[string]string{
		"label": label,
		"color": "cyan",
		"style": "filled",
	}))

	// parse the "if" path
	prevRef, counter = p.nodeToGraph(g, ifStatement.ChildByFieldName("consequence"), ifStartRef, counter)
	g.AddEdge(prevRef, ifEndRef)
	prevRef, counter = p.nodeToGraph(g, ifStatement.ChildByFieldName("alternative"), ifStartRef, counter)
	g.AddEdge(prevRef, ifEndRef)

	return ifEndRef, counter
}

func (p *Parser) unknownToGraph(g graph.Graph[int, int], node *sitter.Node, prevRef, counter int) (int, int) {
	counter++
	label := fmt.Sprintf("%d %s", counter, node.Type())
	g.AddVertex(counter, graph.VertexAttribute("label", label))
	g.AddEdge(prevRef, counter)
	return counter, counter
}
