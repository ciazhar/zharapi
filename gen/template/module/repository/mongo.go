package repository

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func InitMongoRepository(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init mongo repository")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(MongoTemplate))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name) + "/repository/mongo/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name)+"/repository/mongo/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/repository/mongo/" + strings.ToLower(d.Name) + "_mongo_repo.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var MongoTemplate = `
package mongodb

import (
	"context"
	"{{.Package}}/app"
	"{{.Package}}/common/rest"
	"{{.Package}}/src/{{.Name | toLower}}/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type {{.Name}}MongoRepository interface {
	Fetch(param rest.Param) ([]*model.{{.Name}}, error)
	GetByID(id string) (model.{{.Name}}, error)
	Store(req *model.{{.Name}}) error
	Update(req *model.{{.Name}}) error
}

type repository struct {
	app *app.Application
	collection *mongo.Collection
}

func (r repository) Fetch(param rest.Param) ([]*model.{{.Name}}, error) {
	{{.Name | toLower}}s := make([]*model.{{.Name}}, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := r.collection.Find(ctx, param.Query)
	if err != nil {
		return {{.Name | toLower}}s, err
	}

	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var {{.Name | toLower}} *model.{{.Name}}
		err := cur.Decode(&{{.Name | toLower}})
		if err != nil {
			return {{.Name | toLower}}s, err
		}
		{{.Name | toLower}}s = append({{.Name | toLower}}s, {{.Name | toLower}})
	}
	if err := cur.Err(); err != nil {
		return {{.Name | toLower}}s, err
	}

	return {{.Name | toLower}}s, nil
}

func (r repository) GetByID(id string) (model.{{.Name}}, error) {
	{{.Name | toLower}} := model.{{.Name }}{}
	
	objId, err := primitive.ObjectIDFromHex(id)
	if err!=nil {
		return model.{{.Name }}{}, err
	}
	
	filter := bson.M{
		"_id":objId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = r.collection.FindOne(ctx, filter).Decode(&{{.Name | toLower}})

	return {{.Name | toLower}}, err
}

func (r repository) Store(req *model.{{.Name}}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.collection.InsertOne(ctx, req)
	return err
}

func (r repository) Update(req *model.{{.Name}}) error {
	return nil
}

func New{{.Name}}MongoRepository(app *app.Application) {{.Name}}MongoRepository {
	db := app.Env.Get("mongo.database")

	r := repository{
		app: app,
		collection: app.Mongo.Database(db).Collection("{{.Name | toLower}}"),
	}

	return r
}

`
