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
	flag.Parse()

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(MainTemplate))
	t.Execute(os.Stdout, d)
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
