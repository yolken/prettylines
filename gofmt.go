package main

import (
	"bytes"
	"os/exec"
)

func formatFile(contents []byte) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})

	cmd := exec.Cmd{
		Path:   "/usr/local/bin/gofmt",
		Stdin:  bytes.NewReader(contents),
		Stdout: output,
	}

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
