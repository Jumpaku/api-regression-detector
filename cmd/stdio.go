package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
)

var Stdio = &cli.Stdio{
	Stdout: os.Stdout,
	Stderr: os.Stderr,
	Stdin:  os.Stdin,
}

func PrintError(stderr io.Writer, err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(stderr, "Error: %v\n", err)
}
