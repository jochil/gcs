package metric

import "github.com/jochil/dlth/pkg/parser"

func CalcScore(candidates []*parser.Candidate) {
	maxCC := 0
	maxLines := 0
	//find max values for normalization
	for _, c := range candidates {
		if c.Lines > maxLines {
			maxLines = c.Lines
		}

		if cc, err := c.CyclomaticComplexity(); err == nil && cc > maxCC {
			maxCC = cc
		}
	}

	for _, c := range candidates {
		cc, _ := c.CyclomaticComplexity()
		normCC := float64(cc) / float64(maxCC)
		normLines := float64(c.Lines) / float64(maxLines)
		c.Score = (normCC * 5) + (normLines * 1)
	}
}
