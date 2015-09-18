package main

import (
	"flag"
	"fmt"
	"os"
)

const package_yml = "package.yml"

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	var err error
	switch flag.Arg(0) {
	case "init":
		err = initial()
	case "get":
		err = get()
	case "here":
		here()
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "goproj: ", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(`Usage of %s:
Tasks:
	goproj init : Initial a go project and set GOPATH to it
	goproj get  : get all dependencies
	goproj here : Set GOPATH to this project
`, os.Args[0])
	os.Exit(1)
}
