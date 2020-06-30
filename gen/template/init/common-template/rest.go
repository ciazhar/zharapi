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

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(RestTemplate))
	t.Execute(os.Stdout, d)
}

var RestTemplate = `
package rest

type Param struct {
	Query  map[string]interface{}
	Var    map[string]interface{}
	Offset int
	Limit  int
}

func NewParam() Param {
	return Param{
		Query:  map[string]interface{}{},
		Var:    map[string]interface{}{},
		Offset: 1,
		Limit:  10,
	}
}

func (it *Param) GetInclude() string {
	if it.Var["include"] != nil {
		return it.Var["include"].(string)
	}
	return ""
}

`
