# goproj

##Install
1. ```go get -u github.com/WitzHsiao/goproj```
2. Remove ```GOPATH``` in ```.bashrc``` or ```.bash_profil```

##Usage
```
Usage of goproj:
Tasks:
	goproj init : Initial a go project and set GOPATH to it
	goproj get  : get all dependencies
	goproj here : Set GOPATH to this project
```
##Package management format
1. Using yaml
2. Sample:
``` yml
deps: [
	"github.com/WitzHsiao/goproj", 
	"gopkg.in/yaml.v2"
	]
```
##Reference
Inspired by [mattn/gom](https://github.com/mattn/gom)
