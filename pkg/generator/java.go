package generator

import (
	"bytes"
	_ "embed"
	"log/slog"
	"strings"
	"text/template"

	"github.com/jochil/dlth/pkg/candidate"
)

//go:embed tmpl/java.tmpl
var javaTemplate []byte

func renderJavaFuzzTest(c *candidate.Candidate) string {
	tmpl, err := template.New("java").Funcs(template.FuncMap{
		"consumeFunc": mapTypeToFuzzedDataProviderFunc,
		"join":        strings.Join,
		"classVar": func(class string) string {
			return strings.ToLower(class[:1]) + class[1:] + "Obj"
		},
	}).Parse(string(javaTemplate))
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

func mapTypeToFuzzedDataProviderFunc(typeName string) string {
	// TODO handle non primitive data types
	switch typeName {
	case "int", "Integer", "AtomicInteger":
		return "consumeInt"
	case "int[]", "Integer[]", "AtomicInteger[]":
		return "consumeInts"
	case "long", "Long", "AtomicLong":
		return "consumeLong"
	case "long[]", "Long[]", "AtomicLong[]":
		return "consumeLongs"
	case "short", "Short":
		return "consumeShort"
	case "short[]", "Short[]":
		return "consumeShorts"
	case "byte", "Byte":
		return "consumeByte"
	case "byte[]", "Byte[]":
		return "consumeBytes"
	case "boolean", "Boolean", "AtomicBoolean":
		return "consumeBoolean"
	case "boolean[]", "Boolean[]", "AtomicBoolean[]":
		return "consumeBooleans"
	case "char", "Character":
		return "consumeChar"
	case "double", "Double":
		return "consumeDouble"
	case "float", "Float":
		return "consumeFloat"
	case "String":
		return "consumeString"
	}
	return "test"
}
