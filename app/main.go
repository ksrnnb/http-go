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
	env := cfg.Section("").Key("SERVER_ENV").String()
	port, err := cfg.Section("").Key("PORT").Int()

	if err != nil {
		fmt.Printf("Port number is invalid: %v\n", err)
		os.Exit(1)
	}

	if env == "socket" {
		doSocketService(port, docroot)
	} else {
		doFileService(docroot)
	}
}

func doFileService(docroot string) {
	f, err := os.Open("test.txt")

	if err != nil {
		fmt.Printf("Fail to read file: %v\n", err)
		os.Exit(1)
	}

	service := service.NewService(f, os.Stdout, docroot)
	err = service.Start()

	if err != nil {
		fmt.Printf("stdin, stdout error: %v\n", err)
		os.Exit(1)
	}

	defer f.Close()
}

func doSocketService(port int, docroot string) {
	socket, err := service.ListenSocket(port)
	if err != nil {
		fmt.Printf("Fail to listen socket: %v\n", err)
		os.Exit(1)
	}

	defer syscall.Close(socket)

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)

	// ctrl + c で中断した場合にsocketをcloseする
	go func() {
		for sig := range ch {
			fmt.Println(sig)
			close(ch)

			syscall.Close(socket)
			os.Exit(1)
		}
	}()

	accept(socket, docroot)
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
	defer sock.Close()

	service := service.NewService(sock, sock, docroot)
	err := service.Start()

	if err != nil {
		fmt.Printf("service error: %v\n", err)
		syscall.Close(nfd)
		os.Exit(1)
	}

	syscall.Close(nfd)
}
