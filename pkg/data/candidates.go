package data

import "fmt"

type Candidate struct {
	Path     string
	Function string
	Class    string
}

func (c *Candidate) String() string {
	return fmt.Sprintf("%s %s (%s)", c.Class, c.Function, c.Path)
}
