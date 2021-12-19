package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Service struct {
	in      *os.File
	out     *os.File
	docroot string
	scanner *bufio.Scanner
}

func NewService(in *os.File, out *os.File, docroot string) *Service {
	scanner := bufio.NewScanner(in)

	return &Service{in, out, docroot, scanner}
}

func (s *Service) Start() error {
	_, err := s.readRequest()

	if err != nil {
		return err
	}

	return nil

	// s.respondTo(req)
}

func (s *Service) readRequest() (*HTTPRequest, error) {
	req := NewHTTPRequest()
	err := s.readRequestLine(req)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", req)

	return req, nil
}

// リクエストライン（1行目）を読み込む
func (s *Service) readRequestLine(req *HTTPRequest) error {
	s.scanner.Split(bufio.ScanWords)

	method, err := s.scanHTTPMethod()

	if err != nil {
		return err
	}

	path, err := s.scanPath()
	if err != nil {
		return err
	}

	protocolMinorVersion, err := s.scanMinorProtocolVersion()

	if err != nil {
		return err
	}

	req.method = method
	req.path = path
	req.protocolMinorVersion = protocolMinorVersion
	return nil
}

func (s *Service) scanHTTPMethod() (method string, err error) {
	if !s.scanner.Scan() {
		return "", errors.New("erorr while read request line")
	}

	method = strings.ToUpper(s.scanner.Text())

	switch method {
	case "GET":
		return method, nil
	default:
		return "", errors.New("method is invalid")
	}
}

func (s *Service) scanPath() (path string, err error) {
	if !s.scanner.Scan() {
		return "", errors.New("path is invalid")
	}

	path = s.scanner.Text()
	return path, nil
}

func (s *Service) scanMinorProtocolVersion() (version int, err error) {
	if !s.scanner.Scan() {
		return 0, errors.New("protocol is invalid")
	}

	proto := s.scanner.Text()

	protoSplitted := strings.Split(proto, "HTTP/1.")

	if len(protoSplitted) != 2 {
		return 0, errors.New("protocol is invalid")
	}

	minorVersion, err := strconv.Atoi(protoSplitted[1])

	if err != nil {
		return 0, err
	}

	return minorVersion, nil
}
