package go_module

import (
	"fmt"
	"os"
	"os/exec"
)

func GoModTidy() {
	output, err := exec.Command("go", "mod", "tidy").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		fmt.Println()
	}
	fmt.Println(string(output))
}
