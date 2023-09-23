package parser

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type Function struct {
	Name       string
	Parameters []*Parameter
}

type Parameter struct {
	Name string
	Type string
}

func (f *Function) String() string {
	return f.Name
}

type Candidate struct {
	Path             string
	Function         *Function
	Class            string
	Package          string
	ControlFlowGraph graph.Graph[int, int]
	Lines            int
	Score            float64
}

func (c *Candidate) String() string {
	return fmt.Sprintf("%s (%s)", c.Function, c.Path)
}

func (c *Candidate) SaveGraph() {
	file, _ := os.Create(fmt.Sprintf("../../.draw/%s.gv", c.Function.Name))
	err := draw.DOT(c.ControlFlowGraph, file)
	if err != nil {
		panic(err)
	}
	slog.Info("saved cfg", "function", c, "file", file.Name())
}

func (c *Candidate) CyclomaticComplexity() (cc int, err error) {
	if c.ControlFlowGraph == nil {
		err = errors.New("no graph found")
		return
	}

	edges, err := c.ControlFlowGraph.Size()
	if err != nil {
		return
	}
	nodes, err := c.ControlFlowGraph.Order()
	if err != nil {
		return
	}
	cc = edges - nodes + 2
	return
}
