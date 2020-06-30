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

	t := template.Must(template.New("t").Funcs(funcMap).Parse(UsecaseTemplate))
	t.Execute(os.Stdout, d)
}

var UsecaseTemplate = `
package usecase

import (
	"errors"
	"{{.Package }}/common/logger"
	"{{.Package }}/common/rest"
	"{{.Package }}/common/validator"
	"{{.Package }}/src/{{.Name | toLower }}/model"
	"{{.Package }}/src/{{.Name | toLower }}/repository/postgres"
	"github.com/imdario/mergo"
	"time"
)

type {{.Name  }}UseCase interface {
	Fetch(param rest.Param) ([]*model.{{.Name  }}, error)
	GetByID(id string) (model.{{.Name  }}, error)
	Store(req *model.{{.Name  }}) error
	Update(req *model.{{.Name  }}) error
	Delete(id string) error
}

type {{.Name | toLower }}UseCase struct {
	{{.Name  }}Repository postgres.{{.Name  }}PostgresRepository
}

func (c {{.Name | toLower }}UseCase) GetByID(id string) (model.{{.Name  }}, error) {
	return c.{{.Name  }}Repository.GetByID(id)
}

func (c {{.Name | toLower }}UseCase) Update(req *model.{{.Name  }}) error {
	oldReq, err := c.{{.Name  }}Repository.GetByID(req.Id)
	if err != nil {
		return logger.WithError(err)
	}

	if err := mergo.Merge(req, oldReq); err != nil {
		return logger.WithError(err)
	}
	if err := validator.Struct(req); err != nil {
		return logger.WithError(err)
	}

	req.CreatedAt = oldReq.CreatedAt
	req.UpdatedAt = time.Now()
	req.DeletedAt = oldReq.DeletedAt

	return c.{{.Name  }}Repository.Update(req)
}

func (c {{.Name | toLower }}UseCase) Delete(id string) error {
	payload, err := c.GetByID(id)
	if err != nil {
		return logger.WithError(err)
	}
	if !payload.DeletedAt.IsZero() {
		return logger.WithError(errors.New("not found"))
	}
	payload.DeletedAt = time.Now()
	return c.{{.Name  }}Repository.Update(&payload)
}

func (c {{.Name | toLower }}UseCase) Fetch(param rest.Param) ([]*model.{{.Name  }}, error) {
	return c.{{.Name  }}Repository.Fetch(param)
}

func (c {{.Name | toLower }}UseCase) Store(req *model.{{.Name  }}) error {
	if err := validator.Struct(req); err != nil {
		return logger.WithError(err)
	}
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	return c.{{.Name  }}Repository.Store(req)
}

func New{{.Name  }}UseCase({{.Name  }}Repository postgres.{{.Name  }}PostgresRepository) {{.Name  }}UseCase {
	return {{.Name | toLower }}UseCase{
		{{.Name  }}Repository: {{.Name  }}Repository,
	}
}
`
