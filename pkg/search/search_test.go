package search_test

import (
	"testing"

	"github.com/jochil/gcs/pkg/candidate"
	"github.com/jochil/gcs/pkg/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	candidates, err := search.Search([]string{"testdata"})
	require.NoError(t, err)
	assert.Len(t, candidates, 6, "wrong number of candidates")
}

func TestSearch_InvalidPath(t *testing.T) {
	_, err := search.Search([]string{"testdata_foo"})
	require.Error(t, err)
}

func TestSearchOptions_Extension(t *testing.T) {
	candidates, err := search.SearchWithOptions([]string{"testdata"}, search.Options{
		Extensions: []string{".java"},
	})
	require.NoError(t, err)
	assert.Len(t, candidates, 5, "wrong number of candidates")
}

func TestSearchOptions_Limit(t *testing.T) {
	candidates, err := search.SearchWithOptions([]string{"testdata"}, search.Options{
		Limit: 2,
	})
	require.NoError(t, err)
	assert.Len(t, candidates, 2, "wrong number of candidates")
}

func TestSearchOptions_LimitBounds(t *testing.T) {
	candidates, err := search.SearchWithOptions([]string{"testdata"}, search.Options{
		Limit: 10,
	})
	require.NoError(t, err)
	assert.Len(t, candidates, 6, "wrong number of candidates")
}

func TestSearchOptions_Filter(t *testing.T) {
	candidates, err := search.SearchWithOptions([]string{"testdata"}, search.Options{
		Filter: func(c *candidate.Candidate) bool {
			return c.Function.Name == "A"
		},
	})
	require.NoError(t, err)
	assert.Len(t, candidates, 1, "wrong number of candidates")
	assert.Equal(t, "A", candidates[0].Function.Name)
}
