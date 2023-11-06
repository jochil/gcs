package generator

import (
	"bytes"
	_ "embed"
	"log/slog"
	"text/template"

	"github.com/jochil/dlth/pkg/candidate"
)

//go:embed tmpl/java.tmpl
var javaTemplate []byte

func renderJavaFuzzTest(c *candidate.Candidate) string {
	tmpl, err := template.New("java").Parse(string(javaTemplate))
	if err != nil {
		slog.Error("unable to load template", "err", err.Error())
		panic(err)
	}
	var out bytes.Buffer
	err = tmpl.Execute(&out, c)
	if err != nil {
		slog.Error("unable to render template", "err", err.Error())
		panic(err)
	}
	return out.String()
}
