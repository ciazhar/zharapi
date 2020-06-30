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

	t := template.Must(template.New("t").Funcs(funcMap).Parse(Template))
	t.Execute(os.Stdout, d)
}

var Template = `
package postgres

import (
	"{{.Package  }}/app"
	"{{.Package  }}/common/logger"
	"{{.Package  }}/common/rest"
	"{{.Package  }}/src/{{.Name | toLower }}/model"
	"github.com/go-pg/pg/v9/orm"
	uuid "github.com/satori/go.uuid"
)

type {{.Name }}PostgresRepository interface {
	Fetch(param rest.Param) ([]*model.{{.Name }}, error)
	GetByID(id string) (model.{{.Name }}, error)
	Store(req *model.{{.Name }}) error
	Update(req *model.{{.Name }}) error
}

type repository struct {
	app *app.Application
}

func (r repository) Fetch(param rest.Param) ([]*model.{{.Name }}, error) {
	{{.Name | toLower }}s := make([]*model.{{.Name }}, 0)
	query := r.app.Postgres.Model(&{{.Name | toLower }}s).
		Where("deleted_at is null").
		Order("created_at DESC").
		Offset(param.Offset).
		Limit(param.Limit)
	if err := query.Select(); err != nil {
		return {{.Name | toLower }}s, logger.WithError(err)
	}
	return {{.Name | toLower }}s, nil
}

func (r repository) GetByID(id string) (model.{{.Name }}, error) {
	{{.Name | toLower }} := model.{{.Name }}{Id: id}
	if err := r.app.Postgres.Select(&{{.Name | toLower }}); err != nil {
		return {{.Name | toLower }}, logger.WithError(err)
	}
	return {{.Name | toLower }}, nil
}

func (r repository) Store(req *model.{{.Name }}) error {
	id := uuid.Must(uuid.NewV4(), nil)
	req.Id = id.String()
	return r.app.Postgres.Insert(req)
}

func (r repository) Update(req *model.{{.Name }}) error {
	return r.app.Postgres.Update(req)
}

func New{{.Name }}PostgresRepository(app *app.Application) {{.Name }}PostgresRepository {
	r := repository{
		app: app,
	}

	if err := r.app.Postgres.CreateTable((*model.{{.Name }})(nil), &orm.CreateTableOptions{
		IfNotExists: true,
		Temp:        false,
	}); err != nil {
		panic(err)
	}

	return r
}
`
