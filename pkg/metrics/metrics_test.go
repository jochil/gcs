package metrics_test

import (
	"testing"

	"github.com/jochil/gcs/pkg/metrics"
	"github.com/jochil/gcs/pkg/types"
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

func TestPrimitiveParametersOnly(t *testing.T) {
	tests := map[string]struct {
		types    []string
		lang     types.Language
		expected bool
	}{
		"java_prim":           {types: []string{"int", "byte", "short", "long", "float", "double", "char", "boolean"}, lang: types.Java, expected: true},
		"java_wrapper":        {types: []string{"Integer", "Byte", "Short", "Long", "Float", "Double", "Character", "Boolean", "String"}, lang: types.Java, expected: true},
		"java_atomic_wrapper": {types: []string{"AtomicInteger", "AtomicLong", "AtomicBoolean"}, lang: types.Java, expected: true},
		"java_class":          {types: []string{"MyClass"}, lang: types.Java, expected: false},
		"java_mixed":          {types: []string{"int", "MyClass"}, lang: types.Java, expected: false},
		"java_mixed_wrapper":  {types: []string{"String", "MyClass"}, lang: types.Java, expected: false},
		"java_Case":           {types: []string{"InTeGER"}, lang: types.Java, expected: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, metrics.HasPrimitiveParametersOnly(tc.types, tc.lang))
		})
	}
}
