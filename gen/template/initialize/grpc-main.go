package initialize

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitGRPCMain(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init grpc")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(GRPCMain))

	if _, err := os.Stat("cmd"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "cmd")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("cmd/grpc-main.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var GRPCMain = `
package main

import (
	"{{.Package}}/app"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
)

func main() {
	//setup app
	application, err := app.SetupApp()
	if err != nil {
		panic(err)
	}

	//setup grpc
	if err := InitGRPC(application); err != nil {
		panic(err)
	}
}

func InitGRPC(application *app.Application) error {

	address := application.Env.Get("grpc.address")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	//init grpc server
	s := grpc.NewServer()

	//init client
	

	//serve grpc server
	log.Info().Caller().Msg("Running GRPC in port : " + address)
	return s.Serve(lis)
}
`
