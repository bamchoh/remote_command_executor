//go:build linux

package main

func generateCommandParams(cmdline string) (cmdName string, cmdParams []string, err error) {
	cmdName = "sh"
	cmdParams = []string{"-c"}
	cmdParams = append(cmdParams, cmdline)
	return
}
