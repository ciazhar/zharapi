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

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ConfigTemplate))

	f, err := os.Create("config.json")
	if err != nil {
		panic(err)
	}

	t.Execute(f, d)
}

var ConfigTemplate = `
{
  "name": "GRPC COMMON CONFIG",
  "profile": "debug",
  "version": "v6.1.1",
  "port":"8080",
  "postgres.host": "localhost",
  "postgres.username": "postgres",
  "postgres.password": "",
  "postgres.database": "orm_test",
  "postgres.port": "5432",
  "grpc.address": "0.0.0.0:50051"
}
`
