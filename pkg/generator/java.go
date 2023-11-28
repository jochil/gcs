package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/jochil/dlth/pkg/candidate"
)

//go:embed tmpl/java.tmpl
var javaTemplate []byte

func renderJavaFuzzTest(c *candidate.Candidate) string {
	tmpl, err := template.New("java").Funcs(template.FuncMap{
		"renderParamsAsVar": renderParamsAsVar,
		"renderClassInit":   renderClassInit,
		"renderMethodCall":  renderMethodCall,
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

func renderObjVar(class string) string {
	return strings.ToLower(class[:1]) + class[1:] + "Obj"
}

func renderMethodCall(c *candidate.Candidate) string {
	return fmt.Sprintf("\t\t%s.%s(%s);", renderObjVar(c.Class.Name), c.Function.Name, strings.Join(c.Function.Parameters.Names(), ", "))
}

func renderParamsAsVar(params candidate.Parameters) string {
	out := ""
	for _, p := range params {
		out += fmt.Sprintf("\t\t%s %s = fuzzData.%s();\n", p.Type, p.Name, consumeFunc(p.Type))
	}
	return out
}

func renderClassInit(class *candidate.Class) string {
	params := ""
	out := ""
	if len(class.Constructors) >= 1 {
		// TODO find a better approach as just taking the first one
		con := class.Constructors[0]
		out += renderParamsAsVar(con.Parameters)
		params = strings.Join(con.Parameters.Names(), ", ")
	}
	out += fmt.Sprintf("\t\t%s %s = new %s(%s);", class.Name, renderObjVar(class.Name), class.Name, params)
	return out
}

func consumeFunc(typeName string) string {
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
		return "consumeRemainingAsString"
	}
	return "???"
}
