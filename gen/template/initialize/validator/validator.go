package validator

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitValidator(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init validator")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ValidatorTemplate))

	if _, err := os.Stat("common/validator"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/validator")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/validator/validator.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var ValidatorTemplate = `
package validator

import "github.com/asaskevich/govalidator"

var MustCheck = false

func Init() {
	MustCheck = true
}

func Struct(payload interface{}) error {
	if MustCheck {
		//validate valid tag
		if _, err := govalidator.ValidateStruct(payload); err != nil {
			return err
		}
	}
	return nil
}
`
