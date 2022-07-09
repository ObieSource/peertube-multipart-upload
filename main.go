package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	TESTSERVERHOSTNAME = "https://peertube.cpy.re"
)

func main() {
	if len(os.Args) == 2 {
		//if strings.ToLower(os.Args[1]) == "list" {
		switch strings.ToLower(os.Args[1]) {
		case "list":
			fallthrough
		case "channels":
			ListUserChansHandler()
		case "help":
			FullHelpHandler()
		default:
			log.Printf("unknown arg %s\n", os.Args[1])
		}
		os.Exit(2)
	}

	/*
		Run the multipart upload
	*/
	input, err, failtext := ReadEnvironmentVars()
	if err != nil {
		fmt.Println(strings.Join(failtext, "\n"))
		os.Exit(1)
	}
	if len(input.Password) == 0 {
		// don't do password TUI if
		// already specified with environment
		// variable
		input.Password, err = PasswordSecret()
	}

	if err != nil {
		panic(err)
	}
	if err = MultipartUploadHandler(input); err != nil {
		log.Printf("multipart error %+v\n", err)
		panic(err)
	}

}
