package metric

import "github.com/jochil/dlth/pkg/parser"

// CalcScore calculates the scores for a list of candidates
// All metrics are getting normalized based against the min/max values
// in the list
func CalcScore(candidates []*parser.Candidate) {
	maxCC := 0
	maxLines := 0
	// find max values for normalization
	for _, c := range candidates {
		if c.Metrics.LinesOfCode > maxLines {
			maxLines = c.Metrics.LinesOfCode
		}

		if c.Metrics.CyclomaticComplexity > maxCC {
			maxCC = c.Metrics.CyclomaticComplexity
		}
	}

	for _, c := range candidates {
		normCC := float64(c.Metrics.CyclomaticComplexity) / float64(maxCC)
		normLines := float64(c.Metrics.LinesOfCode) / float64(maxLines)

		// applying different weights for the single metrics
		c.Score = (normCC * 5) + (normLines * 1)
	}
}
