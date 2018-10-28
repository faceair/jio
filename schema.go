package jio

type Schema interface {
	Validate(*Context) error
}

func boolPtr(value bool) *bool {
	return &value
}
