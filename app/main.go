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
	port, err := cfg.Section("").Key("PORT").Int()

	if err != nil {
		fmt.Printf("Port number is invalid: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Open("test.txt")

	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	_, err = service.ListenSocket(port)
	if err != nil {
		fmt.Printf("Fail to listen socket: %v\n", err)
		os.Exit(1)
	}

	service := service.NewService(f, os.Stdout, docroot)
	err = service.Start()

	if err != nil {
		fmt.Printf("service error: %v\n", err)
		os.Exit(1)
	}
}
