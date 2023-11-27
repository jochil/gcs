package candidate

import "log/slog"

// CalcScore calculates the scores for a list of candidates
// All metrics are getting normalized based against the min/max values
// in the list
func CalcScore(candidates []*Candidate) {
	slog.Info("calculating score for candidates")
	maxCC := 0
	maxLines := 0
	// find max values for normalization
	for _, c := range candidates {
		c.CalculateMetrics()

		if c.Metrics.LinesOfCode > maxLines {
			maxLines = c.Metrics.LinesOfCode
		}

		if c.Metrics.CyclomaticComplexity > maxCC {
			maxCC = c.Metrics.CyclomaticComplexity
		}
	}

	// weights for the different metrics
	w := map[string]float64{
		"cc":   4,
		"loc":  1,
		"name": 5,
	}

	for _, c := range candidates {
		normCC := float64(c.Metrics.CyclomaticComplexity) / float64(maxCC)
		normLines := float64(c.Metrics.LinesOfCode) / float64(maxLines)
		var normName float64 = 0
		if c.Metrics.FuzzFriendlyName {
			normName = 1
		}

		// applying different weights for the single metrics
		c.Score =
			(normCC * w["cc"]) +
				(normLines * w["loc"]) +
				(normName * w["name"])
	}
}
