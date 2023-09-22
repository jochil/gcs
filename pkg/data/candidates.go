package data

import (
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
	return fmt.Sprintf("%s", f.Name)
}

type Candidate struct {
	Path             string
	Function         *Function
	Class            string
	Package          string
	ControlFlowGraph graph.Graph[int, int]
}

func (c *Candidate) String() string {
	return fmt.Sprintf("%s (%s)", c.Function, c.Path)
}

func (c *Candidate) SaveGraph() {
	file, _ := os.Create(fmt.Sprintf(".draw/%s.gv", c.Function.Name))
	err := draw.DOT(c.ControlFlowGraph, file)
	if err != nil {
		panic(err)
	}
	slog.Info("saved cfg", "function", c, "file", file.Name())
}
