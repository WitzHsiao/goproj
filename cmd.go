package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"

	"gopkg.in/yaml.v2"
)

func get() error {
	p := Package{}

	data, err := ioutil.ReadFile(package_yml)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return err
	}
	for _, dep := range p.Deps {
		cmd := exec.Command("go", "get", "-v", dep)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func here() {
	setEnv()
}

func setEnv() {
	// Get the current user.
	me, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Set an environment variable.
	pwd, _ := os.Getwd()
	os.Setenv("GOPATH", pwd)

	// Transfer stdin, stdout, and stderr to the new process
	// and also set target directory for the shell to start in.
	pa := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   pwd,
	}

	// Start up a new shell.
	// Note that we supply "login" twice.
	// -fpl means "don't prompt for PW and pass through environment."
	fmt.Print(">> Starting a new interactive shell and set GOPATH\n")

	proc, err := os.StartProcess("/usr/bin/login", []string{"login", "-fpl", me.Username}, &pa)
	if err != nil {
		panic(err)
	}

	//// Wait until user exits the shell
	state, err := proc.Wait()
	if err != nil {
		panic(err)
	}

	//// Keep on keepin' on.
	fmt.Printf("<< Exited shell: %s\n", state.String())
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

func initial() error {
	err := genPackageYml()
	if err != nil {
		return err
	}
	here()
	return nil
}
