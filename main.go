package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed xsd2go.xsl
var xslFile []byte

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Specify schema file")
		os.Exit(1)
	}
	if len(os.Args) < 3 {
		fmt.Println("Specify output file")
		os.Exit(1)
	}
	params := make([]string, 0, 3+(len(os.Args)-3)*3)
	if len(os.Args) > 3 {
		for _, p := range os.Args[2:] {
			kv := strings.Split(p, "=")
			if len(kv) == 2 {
				params = append(params, "--stringparam", kv[0], kv[1])
			}
		}
	}
	params = append(params, "-o")
	params = append(params, os.Args[2])
	params = append(params, "-")
	params = append(params, os.Args[1])
	cmd := exec.Command("xsltproc", params...)
	cmd.Stdin = bytes.NewReader(xslFile)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Print(string(exitErr.Stderr))
			os.Exit(1)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Print(string(out))
}
