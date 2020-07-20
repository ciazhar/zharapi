package validator

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func InitPostgresValidator(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init postgres validator")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(PgValidatorTemplate))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name) + "/validator/postgres/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name)+"/validator/postgres/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/validator/postgres/" + strings.ToLower(d.Name) + "_pg_validator.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
