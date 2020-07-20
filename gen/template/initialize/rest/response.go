package rest

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitResponse(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init response")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ResponseTemplate))

	if _, err := os.Stat("common/rest"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/rest")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/rest/response.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
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
