package initialize

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"html/template"
	"os"
	"path/filepath"
)

func InitApp(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init app")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(AppTemplate))

	if _, err := os.Stat("app"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "app")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("app/app.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
