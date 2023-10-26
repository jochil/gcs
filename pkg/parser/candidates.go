package parser

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (p *Parameter) String() string {
	return fmt.Sprintf("%s:%s", p.Name, p.Type)
}

type Function struct {
	Name         string       `json:"name"`
	Parameters   []*Parameter `json:"parameters"`
	ReturnValues []*Parameter `json:"return_values"`
}

func (f *Function) String() string {
	params := ""
	for _, p := range f.Parameters {
		params += fmt.Sprintf(" %s:%s ", p.Name, p.Type)
	}
	returnValues := ""
	for _, rv := range f.ReturnValues {
		returnValues += fmt.Sprintf(" %s:%s ", rv.Name, rv.Type)
	}
	return fmt.Sprintf("%s(%s)(%s)", f.Name, params, returnValues)
}

type Metrics struct {
	CyclomaticComplexity int
	LinesOfCode          int
}

type Candidate struct {
	Path             string                `json:"path"`
	Function         *Function             `json:"function"`
	Class            string                `json:"class,omitempty"`
	Package          string                `json:"package,omitempty"`
	ControlFlowGraph graph.Graph[int, int] `json:"-"`
	Score            float64               `json:"score"`
	Metrics          *Metrics              `json:"metrics"`
	Code             string                `json:"code"`
}

func (c *Candidate) String() string {
	out := c.Function.String()
	if c.Class != "" {
		out = fmt.Sprintf("%s:%s", c.Class, out)
	}
	if c.Package != "" {
		out = fmt.Sprintf("%s.%s", c.Package, out)
	}
	return out
}

func (c *Candidate) SaveGraph() {
	file, _ := os.Create(fmt.Sprintf("../../.draw/%s.gv", c.Function.Name))
	err := draw.DOT(c.ControlFlowGraph, file)
	if err != nil {
		panic(err)
	}
	slog.Info("saved cfg", "function", c, "file", file.Name())
}

func (c *Candidate) CalcCyclomaticComplexity() (cc int, err error) {
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
