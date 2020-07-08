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

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(StringTemplate))
	t.Execute(os.Stdout, d)
}

var StringTemplate = `
package string

func Contains(s []interface{}, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
`
