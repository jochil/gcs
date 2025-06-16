package candidate_test

import (
	"testing"

	"github.com/jochil/gcs/pkg/candidate"
	"github.com/jochil/gcs/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {

	candidates := candidate.Candidates{
		{Function: &candidate.Function{Name: "A"}},
		{Function: &candidate.Function{Name: "B"}, Language: types.Go},
		// This should be the only remaining candidate
		{Function: &candidate.Function{Name: "C", Parameters: candidate.Parameters{&candidate.Parameter{Type: "String"}}}, Language: types.Java},
		{Function: &candidate.Function{Name: "D", Parameters: candidate.Parameters{&candidate.Parameter{Type: "MyClass"}}}, Language: types.Java},
	}

	candidates.CalcScore()

	filtered := candidates.Filter(func(c *candidate.Candidate) bool {
		if c.Function.Name == "A" || c.Language == types.Go || !c.Metrics.PrimitiveParametersOnly {
			return false
		}
		return true
	})

	assert.Len(t, filtered, 1)
	assert.Equal(t, "C", filtered[0].Function.Name)
}
