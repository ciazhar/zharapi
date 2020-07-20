package main

import (
	"github.com/ciazhar/zharapi/gen/template/data"
	"github.com/ciazhar/zharapi/gen/usecase"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
	"text/template"
)

func main() {

	var appName string
	var packageName string
	var moduleName string
	var protocol string
	var db string
	var mq string

	funcMap := template.FuncMap{
		"toLower": strings.ToLower,
	}

	var d data.Data

	app := &cli.App{
		Name:  "zharapi",
		Usage: "golang project generator",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "app",
				Usage:       "app name",
				Destination: &appName,
			},
			&cli.StringFlag{
				Name:        "package",
				Usage:       "package name",
				Destination: &packageName,
			},
			&cli.StringFlag{
				Name:        "name",
				Usage:       "name for the module",
				Destination: &moduleName,
			},
			&cli.StringFlag{
				Name:        "protocol",
				Usage:       "protocol for the module",
				Destination: &protocol,
			},
			&cli.StringFlag{
				Name:        "db",
				Usage:       "db for the module",
				Destination: &db,
			},
			&cli.StringFlag{
				Name:        "mq",
				Usage:       "mq for the module",
				Destination: &mq,
			},
		},
		Action: func(c *cli.Context) error {

			d.Package = packageName
			d.Name = moduleName

			switch c.String("app") {
			case "init":
				usecase.Init(d, funcMap)
			case "module":
			}

			return nil
		},
	}

	if err := (app).Run(os.Args); err != nil {
		panic(err)
	}
}
