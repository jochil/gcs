package generator

import (
	"github.com/jochil/gcs/pkg/candidate"
	"github.com/jochil/gcs/pkg/types"
)

func Render(c *candidate.Candidate) string {
	switch c.Language {
	case types.Go:
		return renderGoUnitTest(c)
	case types.Java:
		return renderJavaFuzzTest(c)
	}

	return ""
}
