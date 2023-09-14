package data

import "fmt"

type Candidate struct {
	Path     string
	Function string
	Class    string
	Package  string
}

func (c *Candidate) String() string {
	return fmt.Sprintf("%s %s %s (%s)", c.Package, c.Class, c.Function, c.Path)
}
