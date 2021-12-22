package service

import (
	"errors"
	"strconv"
	"strings"
)

const MAX_CONTENT_LENGTH = 10000

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

	// 見つからない場合は0とする
	if err != nil {
		return 0, nil
	}

	trimemdValue := strings.Trim(value, " ")
	length, err = strconv.Atoi(trimemdValue)

	if err != nil {
		return 0, err
	}

	if length < 0 || length > MAX_CONTENT_LENGTH {
		return 0, errors.New("content length is invalid")
	}

	return length, nil
}

func (req *HTTPRequest) FindHeaderFieldValue(field string) (value string, err error) {
	header := req.header
	for header != nil {
		if header.name == field {
			return header.value, nil
		}

		header = header.next
	}

	return "0", nil
}
