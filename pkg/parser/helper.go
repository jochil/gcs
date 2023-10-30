package parser

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/c"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
)

// Map of supported tree-sitter languages indexed by the
// file extension
var SupportedExt = map[string]Language{
	".go":   Go,
	".java": Java,
	".js":   JavaScript,
	".c":    C,
}

type Language int64

const (
	Go Language = iota
	Java
	JavaScript
	C
)

var sitterLanguages = map[Language]*sitter.Language{
	Go:         golang.GetLanguage(),
	Java:       java.GetLanguage(),
	JavaScript: javascript.GetLanguage(),
	C:          c.GetLanguage(),
}

// GuessLanguage returns the tree-sitter language for
// supported languages (based on file extension)
func GuessLanguage(path string) (string, Language) {
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

// Prints a tree-sitter node nicely
//
//nolint:unused
func print(node *sitter.Node, ident int) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", ident), node.Type())

	for i := 0; i < int(node.NamedChildCount()); i++ {
		print(node.NamedChild(i), ident+1)
	}
}
