package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

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
	pwd, _ := os.Getwd()

	pa := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Dir:   pwd,
		Env:   append(os.Environ(), fmt.Sprintf("GOPATH=%s", pwd)),
	}

	fmt.Print(">> Starting a new interactive shell and set GOPATH\n")

	proc, err := os.StartProcess("/usr/bin/env", []string{"/usr/bin/env", "zsh"}, pa)
	if err != nil {
		proc, err = os.StartProcess("/usr/bin/env", []string{"/usr/bin/env", "bash"}, pa)
		if err != nil {
			panic(err)
		}
	}
	_, err = proc.Wait()
	if err != nil {
		panic(err)
	}
	fmt.Print(">> Exit shell\n")
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
