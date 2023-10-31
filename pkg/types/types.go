package types

import "encoding/json"

type Language int64

func (l Language) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.String())
}

func (l Language) String() string {
	switch l {
	case Go:
		return "Go"
	case Java:
		return "Java"
	case JavaScript:
		return "JavaScript"
	case C:
		return "C"
	}

	return NoName
}

const (
	Go Language = iota
	Java
	JavaScript
	C
)

const (
	NoName              string = "?"
	VisibilityPublic    string = "public"
	VisibilityPrivate   string = "private"
	VisibilityProtected string = "protected"
)
