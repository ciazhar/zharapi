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
	flag.Parse()

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(GatewayTemplate))
	t.Execute(os.Stdout, d)
}

var GatewayTemplate = `
package main

import (
	"flag"
	"{{.Package}}/common/middleware"
	"{{.Package}}/grpc/generated/golang"
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
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}

	//register module here
	if err := golang.RegisterListServiceHandlerFromEndpoint(ctx, mux, *serverEndpoint, dialOpts); err != nil {
		return mux, err
	}

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
