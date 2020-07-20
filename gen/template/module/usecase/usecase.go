package usecase

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func InitUseCase(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init use case")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(UsecaseTemplate))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name) + "/usecase/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name)+"/usecase/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/usecase/" + strings.ToLower(d.Name) + "_usecase.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
