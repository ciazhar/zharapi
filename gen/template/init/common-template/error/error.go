package main

import (
	"flag"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"strings"
	"text/template"
)

func main() {
	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
	}

	var d data.Data
	flag.StringVar(&d.Package, "package", "github.com/ciazhar/example", "The package used for the queue being generated")
	flag.Parse()

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ErrorTemplate))
	t.Execute(os.Stdout, d)
}

var ErrorTemplate = `
package error

type Status string

const (
	WrongInput Status = "01"
	NotExist   Status = "02"
)

type Error struct {
	Error  string 
	Status Status 
}

func New(err error) Error {
	return Error{Error: err.Error()}
}

func NewF(err string) Error {
	return Error{
		Error:  err,
	}
}

func NewS(err error, status Status) Error {
	return Error{Error: err.Error(), Status: status}
}
`
