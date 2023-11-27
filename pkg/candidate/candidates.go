package candidate

import "log/slog"

type Candidates []*Candidate

// CalcScore calculates the scores for a list of candidates
// All metrics are getting normalized based against the min/max values
// in the list
func (candidates Candidates) CalcScore() {
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
		"prim": 0, // not weighting it by now, as this is more of a filter
	}

	normBool := func(val bool) float64 {
		if val {
			return 1
		}
		return 0
	}

	for _, c := range candidates {
		normCC := float64(c.Metrics.CyclomaticComplexity) / float64(maxCC)
		normLines := float64(c.Metrics.LinesOfCode) / float64(maxLines)
		normName := normBool(c.Metrics.FuzzFriendlyName)
		normPrim := normBool(c.Metrics.PrimitiveParametersOnly)

		// applying different weights for the single metrics
		c.Score =
			(normCC * w["cc"]) +
				(normLines * w["loc"]) +
				(normName * w["name"]) +
				(normPrim * w["prim"])
	}
}

func (candidates Candidates) Filter(filter func(*Candidate) bool) Candidates {
	if filter == nil {
		return candidates
	}
	filtered := Candidates{}
	for _, c := range candidates {
		if filter(c) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
