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
package rest

import (
	"{{.Package }}/common/logger"
	"{{.Package }}/common/rest"
	"{{.Package }}/src/{{.Name | toLower }}/model"
	"{{.Package }}/src/{{.Name | toLower }}/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type {{.Name  }}Controller interface {
	Fetch(c *gin.Context)
	GetByID(c *gin.Context)
	Store(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
}

type {{.Name | toLower }}Controller struct {
	{{.Name  }}UseCase usecase.{{.Name  }}UseCase
}

func (it {{.Name | toLower }}Controller) Fetch(c *gin.Context) {
	param := rest.NewParamGin(c)
	param.Paging()

	payload, err := it.{{.Name  }}UseCase.Fetch(param.Param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, logger.ErrorS(err))
		return
	}
	c.JSON(http.StatusOK, payload)
}

func (it {{.Name | toLower }}Controller) GetByID(c *gin.Context) {
	id := rest.RequestPathVariableString(c, "id")

	payload, err := it.{{.Name  }}UseCase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, logger.ErrorS(err))
		return
	}
	c.JSON(http.StatusOK, payload)
}

func (it {{.Name | toLower }}Controller) Store(ctx *gin.Context) {
	var payload model.{{.Name  }}
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, logger.WarnS(err))
		return
	}

	err := it.{{.Name  }}UseCase.Store(&payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, logger.ErrorS(err))
		return
	}

	ctx.JSON(http.StatusOK, payload)
}

func (it {{.Name | toLower }}Controller) Update(ctx *gin.Context) {
	var payload model.{{.Name  }}
	if err := ctx.BindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, logger.WarnS(err))
		return
	}

	err := it.{{.Name  }}UseCase.Update(&payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, logger.ErrorS(err))
		return
	}

	ctx.JSON(http.StatusOK, payload)
}

func (it {{.Name | toLower }}Controller) Delete(c *gin.Context) {
	id := rest.RequestPathVariableString(c, "id")

	if err := it.{{.Name  }}UseCase.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, logger.ErrorS(err))
		return
	}
	c.JSON(http.StatusOK, nil)
}

func New{{.Name  }}Controller({{.Name  }}UseCase usecase.{{.Name  }}UseCase) {{.Name  }}Controller {
	return {{.Name | toLower }}Controller{
		{{.Name  }}UseCase: {{.Name  }}UseCase,
	}
}

`
