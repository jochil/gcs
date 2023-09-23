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

var SupportedExt = map[string]*sitter.Language{
	".go":   golang.GetLanguage(),
	".java": java.GetLanguage(),
	".c":    c.GetLanguage(),
	".js":   javascript.GetLanguage(),
}

func GuessLanguage(path string) (string, *sitter.Language) {
	ext := filepath.Ext(path)
	slog.Info("guess language", "path", path, "ext", ext)

	if language, ok := SupportedExt[ext]; ok {
		return path, language
	} else {
		slog.Error("unable to guess language", "path", path)
		os.Exit(1)
		return "", nil
	}
}

//nolint:unused
func print(node *sitter.Node, ident int) {
	fmt.Printf("%s%s\n", strings.Repeat("\t", ident), node.Type())

	for i := 0; i < int(node.NamedChildCount()); i++ {
		print(node.NamedChild(i), ident+1)
	}
}
