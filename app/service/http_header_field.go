package service

type HTTPHeaderField struct {
	name  string
	value string
	next  *HTTPHeaderField
}

func NewHTTPHeaderField(name string, value string) *HTTPHeaderField {
	return &HTTPHeaderField{name: name, value: value}
}
