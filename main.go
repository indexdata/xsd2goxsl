package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

//go:embed xsd2go.xsl
var xslFile []byte

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Specify XSD file")
	}
	params := make([]string, 0, 2+(len(os.Args)-2)*3)
	if len(os.Args) > 2 {
		for _, p := range os.Args[2:] {
			kv := strings.Split(p, "=")
			if len(kv) == 2 {
				params = append(params, "--stringparam", kv[0], kv[1])
			}
		}
	}
	params = append(params, "-")
	params = append(params, os.Args[1])
	cmd := exec.Command("xsltproc", params...)
	cmd.Stdin = bytes.NewReader(xslFile)
	goFile, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Fatal(string(exitErr.Stderr))
		} else {
			log.Fatal(err)
		}
	}
	fmt.Print(string(goFile))
}
