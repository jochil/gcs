package generator

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/jochil/dlth/pkg/candidate"
)

// CreateGoTest generates the test source code for a given candidate
func renderGoUnitTest(c *candidate.Candidate) string {
	// TODO add package to candidate
	goPackage := "foo"

	// statement for the function body
	block := jen.Statement{}

	// declare variables for every parameter used for calling
	// the function under test
	callParams := jen.Statement{}
	for _, param := range c.Function.Parameters {
		id := jen.Id(param.Name)
		v := jen.Var().Add(id)

		// handling slices
		bType := param.Type
		if strings.HasPrefix(bType, "[]") {
			v.Index()
			bType = bType[2:]
		}

		// find parameter type for jennifer
		switch bType {
		case "string":
			v.String()
		case "bool":
			v.Bool()
		case "byte":
			v.Byte()
		case "rune":
			v.Rune()
		case "uintptr":
			v.Uintptr()
		case "int":
			v.Int()
		case "int8":
			v.Int8()
		case "int16":
			v.Int16()
		case "int32":
			v.Int32()
		case "int64":
			v.Int64()
		case "uint":
			v.Uint()
		case "uint8":
			v.Uint8()
		case "uint16":
			v.Uint16()
		case "uint32":
			v.Uint32()
		case "uint64":
			v.Uint64()
		case "complex64":
			v.Complex64()
		case "complex128":
			v.Complex128()
		case "float32":
			v.Float32()
		case "float64":
			v.Float64()
		default:
			v.Interface()
		}

		callParams = append(callParams, id)
		block = append(block, v)
	}

	// create the function call for the function under test
	call := jen.Qual(goPackage, c.Function.Name).Call(callParams...)
	block = append(block, call)

	// new source code file
	f := jen.NewFile(fmt.Sprintf("%s_test", goPackage))

	// create the test function
	f.Func().Id(fmt.Sprintf("Test%s", c.Function)).
		Params(jen.Id("t").Op("*").Qual("testing", "T")).
		Block(block...)

	buf := &bytes.Buffer{}
	err := f.Render(buf)
	if err != nil {
		slog.Error("unable to render function", "func", c.Function.Name)
	}
	return buf.String()
}
