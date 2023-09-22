package parser

import (
	"log/slog"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func GuessLanguage(path string) (string, *sitter.Language) {
	var language *sitter.Language
	ext := filepath.Ext(path)
	slog.Info("guess language", "path", path, "ext", ext)
	switch ext {
	case ".go":
		language = golang.GetLanguage()
  default:
    slog.Error("unable to guess language", "path", path)
    os.Exit(1)
	}
	return path, language
}
