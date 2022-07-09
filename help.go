package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

//go:embed help.txt
var HelpText []byte

type HelpTextStr struct {
	Name string
}

func FullHelpHandler() {
	tmpl, err := template.New("help").Parse(string(HelpText))
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer

	if err = tmpl.Execute(&out, HelpTextStr{os.Args[0]}); err != nil {
		panic(err)
	}

	fmt.Printf("\n%s\n", out.String())
}
