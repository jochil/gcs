package data

import "fmt"

type Candidate struct {
	Path string
	Name string
}

func (c *Candidate) String() string {
	return fmt.Sprintf("%s (%s)", c.Name, c.Path)
}
