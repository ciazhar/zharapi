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

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(GRPCMain))
	t.Execute(os.Stdout, d)
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
