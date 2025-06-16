package filter

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/jochil/gcs/pkg/helper"
)

func Valid(path string, includedExtensions []string) bool {
	ext := filepath.Ext(path)

	if len(includedExtensions) > 0 && !slices.Contains(includedExtensions, ext) {
		return false
	}

	// unsupported extension
	if _, ok := helper.SupportedExt[ext]; !ok {
		return false
	}

	// filter go tests
	if strings.HasSuffix(path, "_test.go") {
		return false
	}

	return true
}
