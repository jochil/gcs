package helper

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/jochil/gcs/pkg/types"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/c"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// Map of supported tree-sitter languages indexed by the
// file extension
var SupportedExt = map[string]types.Language{
	".go":   types.Go,
	".java": types.Java,
	".js":   types.JavaScript,
	".c":    types.C,
	".ts":   types.TypeScript,
}

var SitterLanguages = map[types.Language]*sitter.Language{
	types.Go:         golang.GetLanguage(),
	types.Java:       java.GetLanguage(),
	types.JavaScript: javascript.GetLanguage(),
	types.TypeScript: typescript.GetLanguage(),
	types.C:          c.GetLanguage(),
}

// GuessLanguage returns the tree-sitter language for
// supported languages (based on file extension)
func GuessLanguage(path string) (string, types.Language) {
	ext := filepath.Ext(path)
	slog.Info("guess language", "path", path, "ext", ext)

	if language, ok := SupportedExt[ext]; ok {
		return path, language
	} else {
		slog.Error("unable to guess language", "path", path)
		os.Exit(1)
		return "", 0
	}
}
