package rest

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitRequest(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init request")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(RequestTemplate))

	if _, err := os.Stat("common/rest"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/rest")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/rest/request.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var RequestTemplate = `
package rest

type Param struct {
	Query  map[string]interface{}
	Var    map[string]interface{}
	Offset int
	Limit  int
}

func NewParam() Param {
	return Param{
		Query:  map[string]interface{}{},
		Var:    map[string]interface{}{},
		Offset: 1,
		Limit:  10,
	}
}

func (it *Param) GetInclude() string {
	if it.Var["include"] != nil {
		return it.Var["include"].(string)
	}
	return ""
}

`
