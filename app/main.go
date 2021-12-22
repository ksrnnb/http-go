package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	socket, err := service.ListenSocket(port)
	if err != nil {
		fmt.Printf("Fail to listen socket: %v\n", err)
		os.Exit(1)
	}

	defer syscall.Close(socket)

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-ch:
		syscall.Close(socket)
		os.Exit(1)
	default:
		accept(socket, docroot)
	}
	// for {
	// 	nfd, _, err := syscall.Accept(socket)

	// 	if err != nil {
	// 		fmt.Printf("Fail to accept socket: %v\n", err)
	// 		os.Exit(1)
	// 	}

	// 	go startService(nfd, docroot)
	// }
}

func accept(socket int, docroot string) {
	for {
		nfd, _, err := syscall.Accept(socket)

		if err != nil {
			fmt.Printf("Fail to accept socket: %v\n", err)
			os.Exit(1)
		}

		go startService(nfd, docroot)
	}
}

func startService(nfd int, docroot string) {
	sock := os.NewFile(uintptr(nfd), "socket")
	service := service.NewService(sock, sock, docroot)
	err := service.Start()

	if err != nil {
		fmt.Printf("service error: %v\n", err)
		syscall.Close(nfd)
		os.Exit(1)
	}

	syscall.Close(nfd)
}
