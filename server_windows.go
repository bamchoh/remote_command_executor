//go:build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func generateCommandParams(cmdline string) (cmdName string, cmdParams []string, err error) {
	cmdName = "cmd"
	cmdParams = []string{"/c"}

	var cmdargs []string
	cmdargs, err = windows.DecomposeCommandLine(cmdline)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdParams = append(cmdParams, cmdargs...)
	return
}
