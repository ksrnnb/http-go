package service

import (
	"fmt"
	"os"
)

type FileInfo struct {
	path string
	size int64
	isOk bool
}

func NewFileInfo(docroot string, path string) (*FileInfo, error) {
	fileInfo := &FileInfo{
		path: buildFSPath(docroot, path),
		isOk: false,
	}

	osFileInfo, err := os.Lstat(fileInfo.path)

	if err != nil {
		return fileInfo, nil
	}

	if osFileInfo.IsDir() {
		return fileInfo, nil
	}

	fileInfo.isOk = true
	fileInfo.size = osFileInfo.Size()

	return fileInfo, nil
}

func buildFSPath(docroot string, path string) string {
	return fmt.Sprintf("%s/%s", docroot, path)
}

// 今回はContent-Typeをtext/htmlで固定
func (f FileInfo) guessContentType() string {
	return "text/html"
}
