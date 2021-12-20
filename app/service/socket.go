package service

import (
	"net"
	"syscall"
)

const backlog = 5

func ListenSocket(port int) (socket int, err error) {
	family := syscall.AF_INET
	sotype := syscall.SOCK_STREAM
	proto := 0
	socket, err = syscall.Socket(family, sotype, proto)
	if err != nil {
		return -1, err
	}

	defer syscall.Close(socket)

	sockAddr, err := sockAddr(port)

	err = syscall.Bind(socket, sockAddr)

	if err != nil {
		return -1, err
	}

	err = syscall.Listen(socket, backlog)

	if err != nil {
		return -1, err
	}

	return socket, nil
}

func sockAddr(port int) (syscall.Sockaddr, error) {
	ip := net.IPv4zero
	ip4 := ip.To4()

	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip4)
	return sa, nil
}
