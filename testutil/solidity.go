package testutil

import (
	"os/exec"
)

func IsSolcInstalled() bool {
	output, err := exec.Command("solc", "--version").Output()
	if err != nil {
		return false
	}

	return len(output) > 0
}
