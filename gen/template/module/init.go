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

	t := template.Must(template.New("template").Funcs(funcMap).Parse(Template))

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/init.go")
	if err != nil {
		panic(err)
	}

	t.Execute(f, d)
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
