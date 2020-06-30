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
	flag.StringVar(&d.Name, "name", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
	flag.Parse()

	t := template.Must(template.New("t").Funcs(funcMap).Parse(ProtoTemplate))
	t.Execute(os.Stdout, d)
}

var ProtoTemplate = `syntax = "proto3";

package proto;

import "google/api/annotations.proto";
import "google/protobuf/timestamp.proto";

option go_package = "generated/golang";

message {{.Name}} {
  // Output only.
  string id = 1;
  google.protobuf.Timestamp created_at = 2;
  google.protobuf.Timestamp updated_at = 3;
  google.protobuf.Timestamp deleted_at = 4;
}

message ListAll{{.Name}}Request {}
message ListAll{{.Name}}Response {
  repeated {{.Name}} {{.Name | toLower }} = 1;
}

service {{.Name}}Service {
  rpc Store ({{.Name}}) returns ({{.Name}}) {
    option (google.api.http) = {
      post: "/{{.Name | toLower }}"
      body: "*"
    };
  }

  rpc Fetch (ListAll{{.Name}}Request) returns (stream ListAll{{.Name}}Response) {
    option (google.api.http) = {
      get: "/{{.Name | toLower }}"
    };
  }

  rpc Update ({{.Name}}) returns ({{.Name}}) {
    // Update maps to HTTP Patch.
    option (google.api.http) = {
      put: "/{{.Name | toLower }}"
      body: "*"
      additional_bindings {
        patch: "/{{.Name | toLower }}"
        body: "*"
      }
    };
  };
}
`
