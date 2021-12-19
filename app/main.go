package main

import (
	"fmt"
	"os"

	"github.com/ksrnnb/http-go/service"
	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("config.ini")

	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		os.Exit(1)
	}

	docroot := cfg.Section("").Key("DOCUMENT_ROOT").String()

	f, err := os.Open("test.txt")

	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	service := service.NewService(f, os.Stdout, docroot)
	err = service.Start()

	if err != nil {
		fmt.Printf("service error: %v\n", err)
		os.Exit(1)
	}
}
