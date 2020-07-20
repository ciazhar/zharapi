package controller

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func InitRestController(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init rest")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(Template))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name) + "/controller/rest/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name)+"/controller/rest/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/controller/rest/" + strings.ToLower(d.Name) + "_rest_controller.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
