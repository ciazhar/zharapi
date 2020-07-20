package controller

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func InitGRPCController(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init grpc")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(GRPCTemplate))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name) + "/controller/grpc/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name)+"/controller/grpc/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/controller/grpc/" + strings.ToLower(d.Name) + "_grpc_controller.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var GRPCTemplate = `
package grpc

import (
	"context"
	"{{.Package }}/common/rest"	
	"{{.Package }}/grpc/generated/golang"	
	"{{.Package }}/src/{{.Name | toLower }}/model"
	"{{.Package }}/src/{{.Name | toLower }}/usecase"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

type {{.Name }}Controller struct {
	uc usecase.{{.Name }}UseCase
}

func (u {{.Name }}Controller) Store(ctx context.Context, {{.Name | toLower }} *golang.{{.Name }}) (*golang.{{.Name }}, error) {
	new{{.Name }} := model.{{.Name }}{
		Id:        {{.Name | toLower }}.Id,
	}	
	if err := u.uc.Store(&new{{.Name }}); err != nil {
		return {{.Name | toLower }}, err
	}
	return {{.Name | toLower }}, nil
}

func (u {{.Name }}Controller) Fetch(request *golang.ListAll{{.Name }}Request, server golang.{{.Name }}Service_FetchServer) error {
	param := rest.NewParam()
	param.Offset = 0
	param.Limit = 10

	{{.Name | toLower }}, err := u.uc.Fetch(param)
	if err != nil {
		return err
	}

	new{{.Name }} := make([]*golang.{{.Name }},0)
	for i := range {{.Name | toLower}} {
		createdAt, err := ptypes.TimestampProto({{.Name | toLower}}[i].CreatedAt)
		if err!=nil {
			return err
		}
		updatedAt, err := ptypes.TimestampProto({{.Name | toLower}}[i].UpdatedAt)
		if err!=nil {
			return err
		}
		deletedAt, err := ptypes.TimestampProto({{.Name | toLower}}[i].DeletedAt)
		if err!=nil {
			return err
		}
		
		new{{.Name }}= append(new{{.Name }}, &golang.{{.Name }}{
			Id:        {{.Name | toLower}}[i].Id,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			DeletedAt: deletedAt,
		})
	}

	return server.Send(&golang.ListAll{{.Name }}Response{
		{{.Name }}: new{{.Name }},
	})
}

func (u {{.Name }}Controller) Update(ctx context.Context, {{.Name | toLower }} *golang.{{.Name }}) (*golang.{{.Name }}, error) {
	new{{.Name }} := model.{{.Name }}{
		Id:        {{.Name | toLower }}.Id,
	}
	if err := u.uc.Update(&new{{.Name }}); err != nil {
		return {{.Name | toLower }}, err
	}
	return {{.Name | toLower }}, nil
}

func New{{.Name }}GRPCController(server *grpc.Server, useCase usecase.{{.Name }}UseCase) {
	golang.Register{{.Name }}ServiceServer(server, &{{.Name }}Controller{
		uc: useCase,
	})
}
`
