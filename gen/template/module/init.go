package module

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func InitModuleInitializer(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init module initializer")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(Template))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name)); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name))
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/init.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var Template = `
package {{.Name | toLower}}

import (
	"{{.Package }}/app"
	"{{.Package }}/src/{{.Name | toLower}}/controller/grpc"
	"{{.Package }}/src/{{.Name | toLower}}/controller/rest"
	"{{.Package }}/src/{{.Name | toLower}}/repository/postgres"
	"{{.Package }}/src/{{.Name | toLower}}/usecase"
	postgres2 "{{.Package }}/src/{{.Name | toLower}}/validator/postgres"
	"github.com/gin-gonic/gin"
	grpc2 "google.golang.org/grpc"
)

func InitHTTP(engine *gin.Engine, routes string, app *app.Application) {
	repo := postgres.New{{.Name }}PostgresRepository(app)
	uc := usecase.New{{.Name }}UseCase(repo)
	controller := rest.New{{.Name }}Controller(uc)
	postgres2.New{{.Name }}PostgresValidator(repo)

	r := engine.Group(routes)
	{
		r.GET("/", controller.Fetch)
		r.GET("/:id", controller.GetByID)
		r.POST("/", controller.Store)
		r.PUT("/", controller.Update)
		r.DELETE("/:id", controller.Delete)
	}
}

func InitGRPC(server *grpc2.Server, app *app.Application) {
	repo := postgres.New{{.Name }}PostgresRepository(app)
	uc := usecase.New{{.Name }}UseCase(repo)
	grpc.New{{.Name }}GRPCController(server, uc)
}
`
