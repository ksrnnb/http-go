package service

import (
	"errors"
	"fmt"
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

func (req *HTTPRequest) ContentLength() (length int, err error) {
	value, err := req.FindHeaderFieldValue("Content-Length")

	if err != nil {
		return 0, err
	}
	trimemdValue := strings.Trim(value, " ")
	return strconv.Atoi(trimemdValue)
}

func (req *HTTPRequest) FindHeaderFieldValue(field string) (value string, err error) {
	header := req.header
	for header != nil {
		if header.name == field {
			return header.value, nil
		}

		header = header.next
	}

	return "", fmt.Errorf("%s field is not found", field)
}
