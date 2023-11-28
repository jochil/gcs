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
		"var":         createVar,
		"constructor": createConstructorCall,
		"call":        call,
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

func classVar(class string) string {
	return strings.ToLower(class[:1]) + class[1:] + "Obj"
}

func call(c *candidate.Candidate) string {
	return fmt.Sprintf("\t%s.%s(%s)", classVar(c.Class.Name), c.Function.Name, strings.Join(c.Function.Parameters.Names(), ", "))
}

func createConstructorCall(class *candidate.Class) string {
	if len(class.Constructors) >= 1 {
		// TODO find a better approach as just taking the first one
		con := class.Constructors[0]
		out := ""
		for _, p := range con.Parameters {
			out += createVar(p)
		}
		out += fmt.Sprintf("\t%s %s = new %s(%s):", class.Name, classVar(class.Name), class.Name, strings.Join(con.Parameters.Names(), ", "))
		return out
	}
	return fmt.Sprintf("\t%s %s = new %s(%s):", class.Name, classVar(class.Name), class.Name, "")
}

func createVar(p *candidate.Parameter) string {
	return fmt.Sprintf("\t%s %s = fuzzData.%s();\n", p.Type, p.Name, consumeFunc(p.Type))
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
		return "consumeString"
	}
	return "???"
}
