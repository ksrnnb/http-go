package service

type HTTPRequest struct {
	protocolMinorVersion int
	method               string
	path                 string
	header               *HTTPHeaderField
	body                 string
	length               int
}

func NewHTTPRequest() *HTTPRequest {
	return &HTTPRequest{}
}
