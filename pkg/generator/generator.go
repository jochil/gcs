package generator

import (
	"github.com/jochil/dlth/pkg/candidate"
	"github.com/jochil/dlth/pkg/types"
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
