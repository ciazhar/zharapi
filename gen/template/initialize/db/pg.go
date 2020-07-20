package db

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitPG(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init pg")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(PgTemplate))

	if _, err := os.Stat("common/db/pg"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/db/pg")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/db/pg/pg.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
