package model

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func InitProto(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init proto")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ProtoTemplate))

	if _, err := os.Stat("grpc/proto/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "grpc/proto/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	if _, err := os.Stat("grpc/generated/golang/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "grpc/generated/golang/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	if _, err := os.Stat("grpc/generated/swagger/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "grpc/generated/swagger/")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("grpc/proto/" + strings.ToLower(d.Name) + ".proto")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	} else {
		output, err := exec.Command("protoc",
			"-I/usr/local/include",
			"-I.",
			"-I"+os.Getenv("GOPATH")+"/src",
			"-I"+os.Getenv("GOPATH")+"/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis",
			"-I"+os.Getenv("GOPATH")+"/src/github.com/grpc-ecosystem/grpc-gateway",
			"--go_out=plugins=grpc:./grpc",
			"--grpc-gateway_out=logtostderr=true:./grpc",
			"--swagger_out=allow_merge=true,merge_file_name=global:./grpc/generated/swagger",
			"grpc/proto/"+strings.ToLower(d.Name)+".proto").CombinedOutput()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			fmt.Println()
		}
		fmt.Println(string(output))
	}
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
