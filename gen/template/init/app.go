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
	flag.Parse()

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(AppTemplate))
	t.Execute(os.Stdout, d)
}

var AppTemplate = `
package app

import (
	pg2 "{{.Package}}/common/db/pg"
	"{{.Package}}/common/env"
	"{{.Package}}/common/logger"
	"{{.Package}}/common/validator"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"os"
)

type Application struct {
	Env      *env.Environtment
	Postgres *pg.DB
}

func SetupApp() (*Application, error) {

	//env
	environment := env.InitEnv()

	//set default timezone
	if err := os.Setenv("TZ", "Asia/Jakarta"); err != nil {
		panic(err.Error())
	}

	//profile
	gin.SetMode(environment.Get("profile"))

	//logger
	logger.InitLogger()

	//postgres
	pgConn := pg2.InitPG(environment)

	//validator
	validator.Init()

	return &Application{
		Env:      environment,
		Postgres: pgConn,
	}, nil
}

func SetupAppWithPath(path string) (*Application, error) {

	//env
	environment := env.InitPath(path)

	//set default timezone
	if err := os.Setenv("TZ", "Asia/Jakarta"); err != nil {
		panic(err.Error())
	}

	//profile
	gin.SetMode(environment.Get("profile"))

	//logger
	logger.InitLogger()

	//postgres
	pgConn := pg2.InitPG(environment)

	//validator
	validator.Init()

	return &Application{
		Env:      environment,
		Postgres: pgConn,
	}, nil
}
`
