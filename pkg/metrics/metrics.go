package metrics

import (
	"errors"
	"log/slog"
	"regexp"
	"strings"

	"github.com/CodeIntelligenceTesting/dlth/pkg/types"
	"github.com/dominikbraun/graph"
)

type Metrics struct {
	LinesOfCode             int
	CyclomaticComplexity    int
	FuzzFriendlyName        bool
	PrimitiveParametersOnly bool
}

func CountLines(sourceCode string) int {
	// TODO count actual lines.. no comments, no empty ones, ...
	lines := strings.Split(strings.ReplaceAll(sourceCode, "\r\n", "\n"), "\n")
	return len(lines)
}

func CalcCyclomaticComplexity(cfg graph.Graph[int, int]) (cc int, err error) {
	if cfg == nil {
		err = errors.New("no graph found")
		return
	}

	edges, err := cfg.Size()
	if err != nil {
		return
	}
	nodes, err := cfg.Order()
	if err != nil {
		return
	}
	cc = edges - nodes + 2
	return
}

func HasFuzzFriendlyName(name string) bool {
	lcSubstrings := []string{
		"encode",
		"decode",
		"parse",
		"encrypt",
		"decrypt",
		"open",
		"load",
	}
	lcName := strings.ToLower(name)

	for _, lcSubstring := range lcSubstrings {
		if strings.Contains(lcName, lcSubstring) {
			return true
		}
	}
	return false
}

func HasPrimitiveParametersOnly(paramTypes []string, lang types.Language) bool {

	// TODO Java handle generic data types like List<String> or Map<String,String>
	primitives := map[types.Language]*regexp.Regexp{
		types.Java: regexp.MustCompile(`^(int|Integer|[Bb]yte|[Ss]hort|[Ll]ong|[Ff]loat|[Dd]ouble|char|Character|[Bb]oolean|String|AtomicBoolean|AtomicLong|AtomicInteger)(\[\]|\.\.\.)?$`),
	}
	re, ok := primitives[lang]
	if !ok {
		slog.Warn("Unsupported language for primitive parameter check", "lang", lang.String())
		return false
	}
	for _, t := range paramTypes {
		if !re.MatchString(t) {
			return false
		}
	}

	return true
}
