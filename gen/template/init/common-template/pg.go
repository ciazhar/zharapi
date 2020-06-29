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

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(PgTemplate))
	t.Execute(os.Stdout, d)
}

var PgTemplate = `
package db

import (
	"context"
	"{{.Package}}/common/env"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/rs/zerolog/log"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	if err != nil {
		return err
	}
	log.Debug().Msg(query)
	return nil
}

func InitPG(environment *env.Environtment) *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     environment.Get("postgres.username"),
		Password: environment.Get("postgres.password"),
		Database: environment.Get("postgres.database"),
		Addr:     environment.Get("postgres.host") + ":" + environment.Get("postgres.port"),
	})
	if gin.Mode() == gin.DebugMode {
		db.AddQueryHook(dbLogger{})
	}
	return db
}
`
