package main

import (
	"flag"
	"github.com/ciazhar/generate/gen/template/data"
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

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ValidatorTemplate))
	t.Execute(os.Stdout, d)
}

var ValidatorTemplate = `
package validator

import "github.com/asaskevich/govalidator"

var MustCheck = false

func Init() {
	MustCheck = true
}

func Struct(payload interface{}) error {
	if MustCheck {
		//validate valid tag
		if _, err := govalidator.ValidateStruct(payload); err != nil {
			return err
		}
	}
	return nil
}
`