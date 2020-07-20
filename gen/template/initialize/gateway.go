package initialize

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"html/template"
	"os"
	"path/filepath"
)

func InitGateway(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init gateway")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(GatewayTemplate))

	if _, err := os.Stat("cmd"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "cmd")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("cmd/gateway.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var GatewayTemplate = `
package main

import (
	"flag"
	"{{.Package}}/common/middleware"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
)

var (
	//register server url here
	serverEndpoint = flag.String("/server", "localhost:50051", "endpoint of YourService")
)

func newGateway(ctx context.Context, opts ...runtime.ServeMuxOption) (http.Handler, error) {
	mux := runtime.NewServeMux(opts...)
	_ = []grpc.DialOption{grpc.WithInsecure()}

	//register module here

	return mux, nil
}

func Run(address string, opts ...runtime.ServeMuxOption) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()

	gw, err := newGateway(ctx, opts...)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)
	mux.Handle("/swagger/", http.StripPrefix("/swagger", http.FileServer(http.Dir("grpc/generated/swagger/grpc/proto"))))

	return http.ListenAndServe(address, middleware.AllowCORS(middleware.CheckAuth(mux)))

}

func main() {
	flag.Parse()
	defer glog.Flush()

	log.Info().Caller().Msg("Running in port : 8080")
	if err := Run(":8080"); err != nil {
		glog.Fatal(err)
	}
}

`
