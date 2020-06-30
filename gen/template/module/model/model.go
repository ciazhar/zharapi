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
	flag.StringVar(&d.Name, "name", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
	flag.Parse()

	t := template.Must(template.New("t").Funcs(funcMap).Parse(Template))
	t.Execute(os.Stdout, d)
}

var Template = `
package model

import "time"

type {{.Name}} struct {
	tableName struct{}  
	Id        string
	CreatedAt time.Time 
	UpdatedAt time.Time 
	DeletedAt time.Time 
}
`
