package candidate

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/jochil/gcs/pkg/cfg"
	"github.com/jochil/gcs/pkg/metrics"
	"github.com/jochil/gcs/pkg/types"
	sitter "github.com/smacker/go-tree-sitter"
)

type Parameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (p *Parameter) String() string {
	return fmt.Sprintf("%s:%s", p.Name, p.Type)
}

type Parameters []*Parameter

func (p Parameters) Types() []string {
	types := []string{}
	for _, param := range p {
		types = append(types, param.Type)
	}
	return types
}

func (p Parameters) Names() []string {
	names := []string{}
	for _, param := range p {
		names = append(names, param.Name)
	}
	return names
}

type Function struct {
	Name         string       `json:"name"`
	Parameters   Parameters   `json:"parameters"`
	ReturnValues []*Parameter `json:"return_values"`
	Visibility   string       `json:"visibility"`
	Static       bool         `json:"static"`
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

	mods := ""
	if f.Static {
		mods += "static "
	}
	if f.Visibility != "" {
		mods += f.Visibility
	}
	return fmt.Sprintf("%s(%s)(%s) [%s]", f.Name, params, returnValues, mods)
}

type Class struct {
	Name         string      `json:"name"`
	Constructors []*Function `json:"constructors"`
}

func (c *Class) String() string {
	return c.Name
}

type Candidate struct {
	Path             string                `json:"path"`
	Function         *Function             `json:"function"`
	Class            *Class                `json:"class,omitempty"`
	Package          string                `json:"package,omitempty"`
	ControlFlowGraph graph.Graph[int, int] `json:"-"`
	Score            float64               `json:"score"`
	Metrics          *metrics.Metrics      `json:"metrics"`
	Code             string                `json:"code"`
	AST              *sitter.Node          `json:"-"`
	Language         types.Language        `json:"language"`
}

func (c *Candidate) String() string {
	out := c.Function.String()
	if c.Class != nil {
		out = fmt.Sprintf("%s:%s", c.Class, out)
	}
	if c.Package != "" {
		out = fmt.Sprintf("%s.%s", c.Package, out)
	}
	return out
}

func (c *Candidate) CalculateMetrics() {
	c.Metrics = &metrics.Metrics{}
	slog.Debug("calculating metrics", "func", c.Function.Name)

	c.Metrics.FuzzFriendlyName = metrics.HasFuzzFriendlyName(c.Function.Name)
	c.Metrics.PrimitiveParametersOnly = metrics.HasPrimitiveParametersOnly(c.Function.Parameters.Types(), c.Language)

	// calculate cfg + metrics for candidate
	if c.AST != nil {
		if body := c.AST.ChildByFieldName("body"); body != nil {
			c.ControlFlowGraph = cfg.Create(body)
			c.Metrics.LinesOfCode = metrics.CountLines(c.Code)
		}
	}

	if c.ControlFlowGraph != nil {
		cc, err := metrics.CalcCyclomaticComplexity(c.ControlFlowGraph)
		if err != nil {
			cc = -1
			slog.Warn("unable to calc cyclomatic complexity", "func", c.Function.Name)
		}
		c.Metrics.CyclomaticComplexity = cc
	}

}

func (c *Candidate) SaveGraph() {
	filename := fmt.Sprintf("%s_%s_%s.gv", c.Package, c.Class, c.Function.Name)
	file, _ := os.Create(filepath.Join("..", "..", ".draw", filename))
	err := draw.DOT(c.ControlFlowGraph, file)
	if err != nil {
		panic(err)
	}
	slog.Info("saved cfg", "function", c, "file", file.Name())
	cmd := exec.Command("dot", "-Tsvg", "-O", file.Name())
	slog.Debug("run command", "cmd", cmd.String())
	err = cmd.Run()
	if err != nil {
		slog.Error("unable to generate svg from gv file", "error", err.Error(), "file", file.Name())
	}
}
