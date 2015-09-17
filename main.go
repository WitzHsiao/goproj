package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"gopkg.in/yaml.v2"
)

const package_yml = "package.yml"

type Package struct {
	Deps []string
}

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
		err = here()
	default:
		usage()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "goproj: ", err)
		os.Exit(1)
	}

	p := Package{}

	data, err := ioutil.ReadFile(package_yml)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &p)
	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}

func initial() error {
	err := genPackageYml()
	if err != nil {
		return err
	}
	err = here()
	if err != nil {
		return err
	}
	return nil
}

func genPackageYml() error {
	_, err := os.Stat(package_yml)
	if err == nil {
		return errors.New("package.yml already exists")
	}
	f, err := os.Create(package_yml)
	if err != nil {
		return err
	}
	defer f.Close()
	yml_btyes, _ := yaml.Marshal(Package{})
	f.WriteString(string(yml_btyes))
	return nil
}

func setEnv() error {
	// Get the current user.
	me, err := user.Current()
	if err != nil {
		return err
	}

	// Get the current working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Set an environment variable.
	pwd, _ := os.Getwd()
	os.Setenv("GOPATH", pwd)

	// Transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   cwd,
	}

	// Start up a new shell.
	// Note that we supply "login" twice.
	// -fpl means "don't prompt for PW and pass through environment."
	fmt.Print(">> Starting a new interactive shell and set GOPATH")
	proc, err := os.StartProcess("/usr/bin/login", []string{"login", "-fpl", me.Username}, &pa)
	if err != nil {
		return err
	}

	// Wait until user exits the shell
	state, err := proc.Wait()
	if err != nil {
		return err
	}

	// Keep on keepin' on.
	fmt.Printf("<< Exited shell: %s\n", state.String())
	return nil
}

func get() error {
	return nil
}

func here() error {
	err := setEnv()
	if err != nil {
		return err
	}
	return nil
}

func usage() {
	fmt.Printf(`Usage of %s:
Tasks:
	goproj init : Initial a go project
	goproj get  : get all dependencies
	goproj here : Set GOPATH to this project
`, os.Args[0])
	os.Exit(1)
}
