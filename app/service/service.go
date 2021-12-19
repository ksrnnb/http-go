package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
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
	req, err := s.readRequestLine()

	if err != nil {
		return nil, err
	}

	err = s.readHeaderField(req)

	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", req.header)

	return req, nil
}

// リクエストライン（1行目）を読み込む
func (s *Service) readRequestLine() (*HTTPRequest, error) {
	// 全てがスペース区切りだと他の読み込みで都合が悪くなる
	// s.scanner.Split(bufio.ScanWords)
	s.scanner.Scan()
	reqLine := s.scanner.Text()
	reqLineSplitted := strings.Split(reqLine, " ")

	if len(reqLineSplitted) != 3 {
		return nil, errors.New("error while parsing request line")
	}

	return NewHTTPRequest(reqLineSplitted[0], reqLineSplitted[1], reqLineSplitted[2])
}
