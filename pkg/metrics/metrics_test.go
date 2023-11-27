package metrics_test

import (
	"testing"

	"github.com/jochil/dlth/pkg/metrics"
	"github.com/stretchr/testify/assert"
)

func TestHasFuzzFriedlyName(t *testing.T) {
	tests := []string{
		"encode",
		"EnCode",
		"ENCODE",
		"enCODE",
		"EnCODE",
		"decode",
		"parse",
		"encrypt",
		"decrypt",
		"open",
		"load",
	}
	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, metrics.HasFuzzFriendlyName(name), "%s should be a fuzz friendly name", name)
		})
	}
}
