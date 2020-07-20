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

func InitModel(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init model")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(Template))

	if _, err := os.Stat("src/" + strings.ToLower(d.Name) + "/model/"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "src/"+strings.ToLower(d.Name)+"/model/")
		err := os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	f, err := os.Create("src/" + strings.ToLower(d.Name) + "/model/" + strings.ToLower(d.Name) + ".go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	} else {
		output, err := exec.Command("gomodifytags", "-file", "src/"+d.Name+"/model/"+strings.ToLower(d.Name)+".go", "-struct", d.Name, "-add-tags", "json", "--skip-unexported", "-w").CombinedOutput()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			fmt.Println()
		}
		fmt.Println(string(output))

		output, err = exec.Command("gomodifytags", "-file", "src/"+d.Name+"/model/"+strings.ToLower(d.Name)+".go", "-line", "6", "-add-tags", "pg:"+strings.ToLower(d.Name), "-w").CombinedOutput()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			fmt.Println()
		}
		fmt.Println(string(output))

	}

}

var Template = `
package model

import "time"

type {{.Name}} struct {
	tableName struct{}  
	Id        string
	CreatedAt time.Time 
	UpdatedAt time.Time 
	DeletedAt time.Time 
}
`
