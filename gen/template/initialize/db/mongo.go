package db

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitMongo(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init mongo")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(MongoTemplate))

	if _, err := os.Stat("common/db/mongo"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/db/mongo")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/db/mongo/mongo.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var MongoTemplate = `
package db

import (
	"context"
	"fmt"
	"{{.Package}}/common/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func InitMongoDB(environment *env.Environtment) *mongo.Client {

	uri := environment.Get("mongo.url")

	//connect client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx,options.Client().ApplyURI(uri))
	if err!=nil {
		panic(err)
	}

	//ping
	err = client.Ping(ctx, readpref.Primary())
	if err!=nil {
		panic(err)
	}

	fmt.Println("Successfully connected and pinged.")

	return client
}`
