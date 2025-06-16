package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/jochil/gcs/pkg/candidate"
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
	params := strings.Join(c.Function.Parameters.Names(), ", ")
	if c.Function.Static {
		return fmt.Sprintf("\t\t%s.%s(%s);", c.Class.Name, c.Function.Name, params)
	} else {
		return fmt.Sprintf("\t\t%s.%s(%s);", renderObjVar(c.Class.Name), c.Function.Name, params)
	}
}

func renderParamsAsVar(params candidate.Parameters) string {
	out := ""
	for _, p := range params {
		out += fmt.Sprintf("\t\t%s %s = %s;\n", p.Type, p.Name, consumeFunc(p.Type))
	}
	return out
}

func renderClassInit(c *candidate.Candidate) string {
	if c.Function.Static {
		return ""
	}
	params := ""
	out := ""
	if len(c.Class.Constructors) >= 1 {
		// TODO find a better approach as just taking the first one
		con := c.Class.Constructors[0]
		out += renderParamsAsVar(con.Parameters)
		params = strings.Join(con.Parameters.Names(), ", ")
	}
	out += fmt.Sprintf("\t\t%s %s = new %s(%s);", c.Class.Name, renderObjVar(c.Class.Name), c.Class.Name, params)
	return out
}

func consumeFunc(typeName string) string {
	// TODO handle non primitive data types
	obj := "fuzzData"
	var funcName string
	switch typeName {
	case "int", "Integer":
		funcName = "consumeInt"
	case "int[]", "Integer[]":
		funcName = "consumeInts"
	case "long", "Long":
		funcName = "consumeLong"
	case "long[]", "Long[]":
		funcName = "consumeLongs"
	case "short", "Short":
		funcName = "consumeShort"
	case "short[]", "Short[]":
		funcName = "consumeShorts"
	case "byte", "Byte":
		funcName = "consumeByte"
	case "byte[]", "Byte[]":
		funcName = "consumeBytes"
	case "boolean", "Boolean":
		funcName = "consumeBoolean"
	case "boolean[]", "Boolean[]":
		funcName = "consumeBooleans"
	case "char", "Character":
		funcName = "consumeChar"
	case "double", "Double":
		funcName = "consumeDouble"
	case "float", "Float":
		funcName = "consumeFloat"
	case "String":
		funcName = "consumeRemainingAsString"
	case "AtomicBoolean":
		return fmt.Sprintf("new AtomicBoolean(%s.consumeBoolean())", obj)
	case "AtomicLong":
		return fmt.Sprintf("new AtomicLong(%s.consumeLong())", obj)
	case "AtomicInteger":
		return fmt.Sprintf("new AtomicInteger(%s.consumeLong())", obj)
	default:
		return "{}"
	}
	return fmt.Sprintf("%s.%s()", obj, funcName)
}
