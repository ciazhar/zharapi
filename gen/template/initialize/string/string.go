package string

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitString(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init string")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(StringTemplate))

	if _, err := os.Stat("common/string"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/string")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/string/string.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var StringTemplate = `
package string

func Contains(s []interface{}, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
`
