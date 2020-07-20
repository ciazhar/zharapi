package error

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func InitError(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init error")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(ErrorTemplate))

	if _, err := os.Stat("common/error"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/error")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/error/error.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	} else {
		output, err := exec.Command("gomodifytags", "-file", "common/error/error.go", "-struct", "Error", "-add-tags", "json", "-w").CombinedOutput()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			fmt.Println()
		}
		fmt.Println(string(output))
	}
}

var ErrorTemplate = `
package error

type Status string

const (
	WrongInput Status = "01"
	NotExist   Status = "02"
)

type Error struct {
	Error  string 
	Status Status 
}

func New(err error) Error {
	return Error{Error: err.Error()}
}

func NewF(err string) Error {
	return Error{
		Error:  err,
	}
}

func NewS(err error, status Status) Error {
	return Error{Error: err.Error(), Status: status}
}
`
