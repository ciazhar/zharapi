package initialize

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitMain(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init main")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(MainTemplate))

	if _, err := os.Stat("cmd"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "cmd")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("cmd/main.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var MainTemplate = `
package main

import (
	"{{.Package}}/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {

	//setup app
	application, err := app.SetupApp()
	if err != nil {
		panic(err)
	}

	//setup http
	if err := InitHTTP(application); err != nil {
		panic(err)
	}
}

func InitHTTP(application *app.Application) error {
	//config router api
	router := gin.New()

	//middleware
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(logger.SetLogger())

	//run
	log.Info().Caller().Msg("Running in port : " + application.Env.Get("port"))
	return router.Run(":" + application.Env.Get("port"))
}
`
