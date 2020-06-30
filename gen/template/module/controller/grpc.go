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
	flag.StringVar(&d.Name, "name", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
	flag.Parse()

	t := template.Must(template.New("t").Funcs(funcMap).Parse(GRPCTemplate))
	t.Execute(os.Stdout, d)
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

func (u {{.Name }}Controller) Fetch(request *golang.{{.Name }}All{{.Name }}Request, server golang.{{.Name }}Service_FetchServer) error {
	{{.Name | toLower }}, err := u.uc.Fetch(rest.NewParam())
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

	return server.Send(&golang.{{.Name }}All{{.Name }}Response{
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
