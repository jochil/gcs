package candidate

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/jochil/dlth/pkg/cfg"
	"github.com/jochil/dlth/pkg/metrics"
	"github.com/jochil/dlth/pkg/types"
	sitter "github.com/smacker/go-tree-sitter"
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

type Candidate struct {
	Path             string                `json:"path"`
	Function         *Function             `json:"function"`
	Class            string                `json:"class,omitempty"`
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
	if c.Class != "" {
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
	// calculate cfg + metrics for candidate
	if body := c.AST.ChildByFieldName("body"); body != nil {
		c.ControlFlowGraph = cfg.Create(body)
		c.Metrics.LinesOfCode = metrics.CountLines(c.Code)
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
