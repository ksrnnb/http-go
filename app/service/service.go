package service

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	httpMinorVersion = 0
	serverVersion    = "1.0"
	serverName       = "MyHTTPServer"
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
	req, err := s.readRequest()

	if err != nil {
		return err
	}

	return s.writeResponse(req)
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

	req.length, err = req.ContentLength()

	if err != nil {
		return nil, err
	}

	err = s.readRequestBody(req)

	if err != nil {
		return nil, err
	}

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

// リクエストヘッダーを読み込む
func (s *Service) readHeaderField(req *HTTPRequest) error {
	header := new(HTTPHeaderField)
	for s.scanner.Scan() {
		line := s.scanner.Text()

		if line == "" {
			break
		}

		h := strings.Split(line, ":")
		if len(h) != 2 {
			return errors.New("header field is invalid")
		}

		header.name = h[0]
		header.value = h[1]
		header.next = req.header
		req.header = header
	}

	return nil
}

// リクエストボディの読み込み
// 最大でContent-Lengthまで
func (s *Service) readRequestBody(req *HTTPRequest) error {
	// ヘッダとボディの間の空行
	// なければ何もしない
	if req.length == 0 {
		return nil
	}

	var body string

	for s.scanner.Scan() {
		body += s.scanner.Text()
	}

	if len(body) < req.length {
		return nil
	}

	req.body = body[:req.length]
	return nil
}

// レスポンスを出力に書き込む
func (s *Service) writeResponse(req *HTTPRequest) error {
	fileInfo, err := NewFileInfo(s.docroot, req.path)

	if err != nil {
		// todo: log
		return err
	}

	if !fileInfo.isOk {
		return s.writeNotFound(req, fileInfo)
	}
	if req.method == "GET" {
		return s.writeOkResponse(req, fileInfo)
	}

	if req.method == "POST" {
		return s.writeNotAllowedResponse(req)
	}

	return s.writeNotImplementedResponse(req)
}

func (s *Service) writeStatusLine(status string) {
	fmt.Fprintf(s.out, "HTTP/1.%d %s\r\n", httpMinorVersion, status)
}

func (s *Service) writeCommonHeaderFields(fileInfo *FileInfo) {
	t := time.Now().Format("Mon, 2 Jan 2006 15:04:05 GMT")
	fmt.Fprintf(s.out, "Date: %s\r\n", t)
	fmt.Fprintf(s.out, "Server: %s/%s\r\n", serverName, serverVersion)
	fmt.Fprintf(s.out, "Connection: close\r\n")
}

func (s *Service) writeOkResponse(req *HTTPRequest, fileInfo *FileInfo) error {
	s.writeStatusLine("200 OK")
	s.writeCommonHeaderFields(fileInfo)
	fmt.Fprintf(s.out, "Content-Length: %d\r\n", fileInfo.size)
	fmt.Fprintf(s.out, "Content-Type: %s\r\n", fileInfo.guessContentType())
	fmt.Fprintf(s.out, "\r\n")

	if req.method == "GET" {
		s.writeResponseBody(fileInfo)
	}

	return nil
}

func (s *Service) writeNotFound(req *HTTPRequest, fileInfo *FileInfo) error {
	s.writeStatusLine("404 Not Found")
	s.writeCommonHeaderFields(fileInfo)
	fmt.Fprintf(s.out, "\r\n")

	return nil
}

func (s *Service) writeNotAllowedResponse(req *HTTPRequest) error {
	fileInfo, err := NewFileInfo(s.docroot, "405.html")

	if err != nil {
		return err
	}

	if !fileInfo.isOk {
		return errors.New("file not found")
	}

	s.writeStatusLine("405 Method Not Allowed")
	s.writeCommonHeaderFields(fileInfo)
	fmt.Fprintf(s.out, "Content-Length: %d\r\n", fileInfo.size)
	fmt.Fprintf(s.out, "Content-Type: %s\r\n", fileInfo.guessContentType())
	fmt.Fprintf(s.out, "\r\n")

	s.writeResponseBody(fileInfo)
	return nil
}

func (s *Service) writeNotImplementedResponse(req *HTTPRequest) error {
	fileInfo, err := NewFileInfo(s.docroot, "501.html")

	if err != nil {
		return err
	}

	if !fileInfo.isOk {
		return errors.New("file not found")
	}

	s.writeStatusLine("501 Not Implemeted")
	s.writeCommonHeaderFields(fileInfo)
	fmt.Fprintf(s.out, "Content-Length: %d\r\n", fileInfo.size)
	fmt.Fprintf(s.out, "Content-Type: %s\r\n", fileInfo.guessContentType())
	fmt.Fprintf(s.out, "\r\n")

	s.writeResponseBody(fileInfo)

	return nil
}

func (s *Service) writeResponseBody(fileInfo *FileInfo) error {
	file, err := os.Open(fileInfo.path)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(s.out, file)
	fmt.Fprintf(s.out, "\r\n")

	return err
}
