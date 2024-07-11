package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"os/exec"
)

func handleStdouterr(done chan struct{}, conn net.Conn, reader io.Reader) {
	defer func() { done <- struct{}{} }()

	b := make([]byte, 1024)
	for {
		n, err := reader.Read(b)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				if ok := err.(*fs.PathError); ok == nil {
					fmt.Printf("%T\n", err)
					fmt.Println(errors.Unwrap(err))
					log.Print("handleStdouterr: Read: ", err)
				}
			}
			return
		}

		_, err = conn.Write(b[:n])
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Print("handleStdouterr: Write: ", err)
			}
			return
		}
	}
}

func proc(conn net.Conn) (err error) {
	defer conn.Close()

	buf := make([]byte, 65535)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("cannot read", err)
	}

	cmdline := string(buf[:n])
	cmdName, cmdParams, err := generateCommandParams(cmdline)
	if err != nil {
		fmt.Println("command params generation was failed: ", err)
		return err
	}

	cmd := exec.Command(cmdName, cmdParams...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Print("StdoutPipe: ", err)
		return err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Print("StderrPipe: ", err)
		return err
	}
	defer stderr.Close()

	stdouterr := io.MultiReader(stdout, stderr)

	done := make(chan struct{})

	err = cmd.Start()
	if err != nil {
		log.Print("Start: ", err)
		return err
	}

	go handleStdouterr(done, conn, stdouterr)

	cmd.Wait()

	<-done

	code := fmt.Sprintf("exit status: %d", cmd.ProcessState.ExitCode())
	conn.Write([]byte(code))

	return nil
}

func main() {
	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("cannot listen", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("cannot accept", err)
			continue
		}

		go proc(conn)
	}
}
