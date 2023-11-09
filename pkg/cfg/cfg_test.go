package cfg_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/helper"
	"github.com/jochil/dlth/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraph(t *testing.T) {
	type edge struct {
		s int
		e int
	}
	type node struct {
		h int
		l string
	}
	tests := map[string]struct {
		path      string
		wantNodes int
		wantEdges int
		edges     []edge
		nodes     []node
	}{
		"go_no_control": {path: "testdata/cyclo/golang/a.go", wantEdges: 3, wantNodes: 4, edges: []edge{{0, 2}, {2, 3}, {3, 1}}},
		"go_simple_if": {
			path:      "testdata/cyclo/golang/b.go",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{2, "if_start"}, {3, "if_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {2, 3}},
		},
		"go_if_else": {
			path:      "testdata/cyclo/golang/c.go",
			wantEdges: 15,
			wantNodes: 13,
			nodes:     []node{{2, "if_start"}, {5, "if_start"}, {8, "if_start"}, {9, "if_end"}, {6, "if_end"}, {3, "if_end"}},
			edges:     []edge{{2, 5}, {5, 8}, {2, 4}, {9, 6}, {6, 3}},
		},
		"go_switch_no_default": {
			path:      "testdata/cyclo/golang/d.go",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 3}, {2, 4}, {4, 3}},
		},
		"go_switch_default": {
			path:      "testdata/cyclo/golang/e.go",
			wantEdges: 8,
			wantNodes: 7,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 4}, {2, 5}, {2, 6}, {4, 3}, {5, 3}, {6, 3}},
		},
		"go_simple_for": {
			path:      "testdata/cyclo/golang/f.go",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "for_start"}, {3, "for_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {3, 2}},
		},
		"java_no_control": {path: "testdata/cyclo/java/NoControl.java", wantEdges: 3, wantNodes: 4, edges: []edge{{0, 2}, {2, 3}, {3, 1}}},
		"java_simple_if": {
			path:      "testdata/cyclo/java/If.java",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{2, "if_start"}, {3, "if_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {2, 3}},
		},
		"java_if_else": {
			path:      "testdata/cyclo/java/IfElse.java",
			wantEdges: 15,
			wantNodes: 13,
			nodes:     []node{{2, "if_start"}, {5, "if_start"}, {8, "if_start"}, {9, "if_end"}, {6, "if_end"}, {3, "if_end"}},
			edges:     []edge{{2, 5}, {5, 8}, {2, 4}, {9, 6}, {6, 3}},
		},
		"java_switch_no_default": {
			path:      "testdata/cyclo/java/Switch.java",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 3}, {2, 4}, {4, 3}},
		},
		"java_switch_default": {
			path:      "testdata/cyclo/java/SwitchDefault.java",
			wantEdges: 8,
			wantNodes: 7,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 4}, {2, 5}, {2, 6}, {4, 3}, {5, 3}, {6, 3}},
		},
		"java_simple_for": {
			path:      "testdata/cyclo/java/For.java",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "for_start"}, {3, "for_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {3, 2}},
		},
		"java_while": {
			path:      "testdata/cyclo/java/While.java",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{3, "while_start"}, {4, "while_end"}},
			edges:     []edge{{3, 4}, {3, 5}, {5, 3}, {4, 1}},
		},
		"java_do": {
			path:      "testdata/cyclo/java/Do.java",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{3, "do_start"}, {4, "do_end"}},
			edges:     []edge{{3, 5}, {5, 3}, {5, 4}, {4, 1}},
		},
		"javascript_no_control": {path: "testdata/cyclo/javascript/noControl.js", wantEdges: 3, wantNodes: 4, edges: []edge{{0, 2}, {2, 3}, {3, 1}}},
		"javascript_simple_if": {
			path:      "testdata/cyclo/javascript/if.js",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{2, "if_start"}, {3, "if_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {2, 3}},
		},
		"javascript_if_else": {
			path:      "testdata/cyclo/javascript/ifElse.js",
			wantEdges: 15,
			wantNodes: 13,
			nodes:     []node{{2, "if_start"}, {5, "if_start"}, {8, "if_start"}, {9, "if_end"}, {6, "if_end"}, {3, "if_end"}},
			edges:     []edge{{2, 5}, {5, 8}, {2, 4}, {9, 6}, {6, 3}},
		},
		"javascript_simple_for": {
			path:      "testdata/cyclo/javascript/for.js",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "for_start"}, {3, "for_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {3, 2}},
		},
		"javascript_switch_no_default": {
			path:      "testdata/cyclo/javascript/switch.js",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 3}, {2, 4}, {4, 3}},
		},
		"javascript_switch_default": {
			path:      "testdata/cyclo/javascript/switchDefault.js",
			wantEdges: 8,
			wantNodes: 7,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 4}, {2, 5}, {2, 6}, {4, 3}, {5, 3}, {6, 3}},
		},
		"javascript_while": {
			path:      "testdata/cyclo/javascript/while.js",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{3, "while_start"}, {4, "while_end"}},
			edges:     []edge{{3, 4}, {3, 5}, {5, 3}, {4, 1}},
		},
		"javascript_do": {
			path:      "testdata/cyclo/javascript/do.js",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{3, "do_start"}, {4, "do_end"}},
			edges:     []edge{{3, 5}, {5, 3}, {5, 4}, {4, 1}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			candidates := parser.NewParser(helper.GuessLanguage(tc.path)).Parse()
			candidate.CalcScore(candidates)
			cfg := candidates[0].ControlFlowGraph
			candidates[0].SaveGraph()

			edges, err := cfg.Size()
			require.NoError(t, err)
			assert.Equal(t, tc.wantEdges, edges, "wrong amount of edges")

			nodes, err := cfg.Order()
			require.NoError(t, err)
			assert.Equal(t, tc.wantNodes, nodes, "wrong amount of nodes")

			for _, e := range tc.edges {
				_, err := cfg.Edge(e.s, e.e)
				assert.NoError(t, err, "missing edge %d -> %d", e.s, e.e)
			}

			for _, n := range tc.nodes {
				_, props, err := cfg.VertexWithProperties(n.h)
				assert.NoError(t, err, "missing node %d %s", n.h, n.l)
				assert.Contains(t, props.Attributes["label"], n.l, "invalid label")
			}
		})
	}
}
