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
	flag.StringVar(&d.Name, "name", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
	flag.Parse()

	t := template.Must(template.New("t").Funcs(funcMap).Parse(PgValidatorTemplate))
	t.Execute(os.Stdout, d)
}

var PgValidatorTemplate = `
package postgres

import (
	"github.com/asaskevich/govalidator"
	"{{.Package }}/src/{{.Name | toLower }}/repository/postgres"
)

type {{.Name }}PostgresValidator interface {
	{{.Name }}MustExist()
}

type {{.Name | toLower }}PostgresValidator struct {
	{{.Name }}Repository postgres.{{.Name }}PostgresRepository
}

func (r {{.Name | toLower }}PostgresValidator) {{.Name }}MustExist() {
	govalidator.TagMap["{{.Name | toLower }}MustExist"] = func(str string) bool {
		return r.validateId(str)
	}
}

func (r {{.Name | toLower }}PostgresValidator) validateId(postId string) bool {
	if postId != "" {
		if _, err := r.{{.Name }}Repository.GetByID(postId); err != nil {
			return false
		}
	}
	return true
}

func New{{.Name }}PostgresValidator({{.Name }}Repository postgres.{{.Name }}PostgresRepository) {
	validator := {{.Name | toLower }}PostgresValidator{
		{{.Name }}Repository: {{.Name }}Repository,
	}
	validator.{{.Name }}MustExist()
}

`
