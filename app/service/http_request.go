package service

import (
	"errors"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	protocolMinorVersion int
	method               string
	path                 string
	header               *HTTPHeaderField
	body                 string
	length               int
}

func NewHTTPRequest(method string, path string, protocol string) (*HTTPRequest, error) {
	method = strings.ToUpper(method)
	switch method {
	case "GET":
		break
	default:
		return nil, errors.New("method is invalid")
	}

	protocolSplitted := strings.Split(protocol, "HTTP/1.")

	if len(protocolSplitted) != 2 {
		return nil, errors.New("protocol is invalid")
	}

	minorVersionInt, err := strconv.Atoi(protocolSplitted[1])

	if err != nil {
		return nil, err
	}

	return &HTTPRequest{method: method, path: path, protocolMinorVersion: minorVersionInt}, nil
}
