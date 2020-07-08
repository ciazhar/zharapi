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

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ResponseTemplate))
	t.Execute(os.Stdout, d)
}

var ResponseTemplate = `
package rest

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, httStatus int, payload interface{}) error {
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(httStatus)
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}
`
