package main

import (
	"fmt"
	"os"
	"os/exec"
)

const VMRUN = "/Applications/VMware Fusion.app/Contents/Library/vmrun"

func main() {
	args := os.Args[1:]
	out, err := exec.Command(VMRUN, args...).Output()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}
	fmt.Printf("%s", out)
}
