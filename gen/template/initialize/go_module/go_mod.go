package go_module

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"os/exec"
)

func InitGoModule(d data.Data, funcMap map[string]interface{}) {
	fmt.Println("inti go mod")

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		output, err := exec.Command("go", "mod", "init", d.Package).CombinedOutput()
		if err != nil {
			os.Stderr.WriteString(err.Error())
			fmt.Println()
		}
		fmt.Println(string(output))
	}
}
