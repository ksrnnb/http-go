package service

type HTTPHeaderField struct {
	name  string
	value string
	next  *HTTPHeaderField
}
