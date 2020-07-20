package repository

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func InitPostgresRepository(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init postgres repository")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(Template))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name) + "/repository/postgres/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name)+"/repository/postgres/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/repository/postgres/" + strings.ToLower(d.Name) + "_pg_repo.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
