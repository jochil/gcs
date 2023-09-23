package parser_test

import (
	"testing"

	"github.com/jochil/test-helper/pkg/parser"
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
		"no_control": {path: "testdata/cyclo/a.go", wantEdges: 3, wantNodes: 4, edges: []edge{{0, 2}, {2, 3}, {3, 1}}},
		"simple_if": {
			path:      "testdata/cyclo/b.go",
			wantEdges: 6,
			wantNodes: 6,
			nodes:     []node{{2, "if_start"}, {3, "if_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {2, 3}},
		},
		"if_else": {
			path:      "testdata/cyclo/c.go",
			wantEdges: 15,
			wantNodes: 13,
			nodes:     []node{{2, "if_start"}, {5, "if_start"}, {8, "if_start"}, {9, "if_end"}, {6, "if_end"}, {3, "if_end"}},
			edges:     []edge{{2, 5}, {5, 8}, {2, 4}, {9, 6}, {6, 3}},
		},
		"switch_no_default": {
			path:      "testdata/cyclo/d.go",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 3}, {2, 4}, {4, 3}},
		},
		"switch_default": {
			path:      "testdata/cyclo/e.go",
			wantEdges: 8,
			wantNodes: 7,
			nodes:     []node{{2, "switch_start"}, {3, "switch_end"}},
			edges:     []edge{{2, 4}, {2, 5}, {2, 6}, {4, 3}, {5, 3}, {6, 3}},
		},
		"simple_for": {
			path:      "testdata/cyclo/f.go",
			wantEdges: 5,
			wantNodes: 5,
			nodes:     []node{{2, "for_start"}, {3, "for_end"}},
			edges:     []edge{{2, 4}, {4, 3}, {3, 2}},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			candidates := parser.NewParser(parser.GuessLanguage(tc.path)).Parse()
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
