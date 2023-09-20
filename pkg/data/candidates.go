package data

import "fmt"

type Candidate struct {
	Path     string
	Function *Function
	Class    string
	Package  string
}

func (c *Candidate) String() string {
	return fmt.Sprintf("%s %s %s (%s)", c.Package, c.Class, c.Function, c.Path)
}

type Function struct {
	Name       string
	Parameters []*Parameter
}

type Parameter struct {
	Name string
	Type string
}

func (f *Function) String() string {
	return fmt.Sprintf("%s", f.Name)
}
