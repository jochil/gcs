package generator

import (
	"fmt"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/jochil/test-helper/pkg/parser"
)

func CreateGoTest(candidate *parser.Candidate) {
	// TODO add package to candidate
	goPackage := "foo"

	block := jen.Statement{}

	callParams := jen.Statement{}
	for _, param := range candidate.Function.Parameters {
		id := jen.Id(param.Name)
		v := jen.Var().Add(id)

		bType := param.Type
		if strings.HasPrefix(bType, "[]") {
			v.Index()
			bType = bType[2:]
		}

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

	call := jen.Qual(goPackage, candidate.Function.Name).Call(callParams...)
	block = append(block, call)

	f := jen.NewFile(fmt.Sprintf("%s_test", goPackage))
	f.Func().Id(fmt.Sprintf("Test%s", candidate.Function)).
		Params(jen.Id("t").Op("*").Qual("testing", "T")).
		Block(block...)
	fmt.Printf("\n----------\n%#v", f)
}
