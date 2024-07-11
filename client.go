//go:build ignore

package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func run(host string, cmdline string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("cannot connect", err)
	}
	defer conn.Close()

	conn.Write([]byte(cmdline))

	for {
		buf := make([]byte, 65535)
		n, err := conn.Read(buf)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Print(err)
			}
			break
		}

		fmt.Print(string(buf[:n]))
	}
}

func main() {
	// cmdline := "i=1; while [ $i -le 1000 ]; do openssl rand -hex 10; i=$(expr $i + 1); done"
	if len(os.Args) >= 3 {
		run(os.Args[1], strings.Join(os.Args[2:], " "))
	}
}
