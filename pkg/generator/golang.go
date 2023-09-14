package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/jochil/test-helper/pkg/data"
)

func CreateGoTest(candidate *data.Candidate) {
	// TODO add package to candidate
	goPackage := "foo"
	f := jen.NewFile(fmt.Sprintf("%s_test", goPackage))
	f.Func().Id(fmt.Sprintf("Test%s", candidate.Function)).Params(jen.Id("t").Op("*").Qual("testing", "T")).Block(

    jen.Qual(goPackage, candidate.Function).Call(),
  )
	fmt.Printf("\n----------\n%#v", f)
}
