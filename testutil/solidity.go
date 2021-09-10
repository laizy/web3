package testutil

import (
	"fmt"
	"os/exec"
)

func IsSolcInstalled() bool {
	output, err := exec.Command("solc", "--version").Output()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return len(output) > 0
}
