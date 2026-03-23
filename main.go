package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

//go:embed xsd2go.xsl
var xslFile []byte

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) < 1 {
		fmt.Fprintln(stderr, "Specify schema file")
		return 1
	}
	if len(args) < 2 {
		fmt.Fprintln(stderr, "Specify output file")
		return 1
	}
	cmd := exec.Command("xsltproc", buildParams(args[0], args[1], args[2:])...)
	cmd.Stdin = bytes.NewReader(xslFile)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			_, _ = stderr.Write(exitErr.Stderr)
			return 1
		}
		fmt.Fprintln(stderr, err)
		return 1
	}
	_, _ = stdout.Write(out)
	return 0
}

func buildParams(inputFile, outputFile string, rawParams []string) []string {
	params := make([]string, 0, 3+len(rawParams)*3)
	for _, p := range rawParams {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			params = append(params, "--stringparam", kv[0], kv[1])
		}
	}
	params = append(params, "-o", outputFile, "-", inputFile)
	return params
}
