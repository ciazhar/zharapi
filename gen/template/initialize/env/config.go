package config

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"text/template"
)

func InitConfigFile(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init config")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ConfigTemplate))

	f, err := os.Create("config.json")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
