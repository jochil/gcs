package types

type Language int64

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
