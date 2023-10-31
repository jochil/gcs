package candidate

import (
	"strings"
)

func CountLines(sourceCode string) int {
	// TODO count actual lines.. no comments, no empty ones, ...
	lines := strings.Split(strings.ReplaceAll(sourceCode, "\r\n", "\n"), "\n")
	return len(lines)
}
