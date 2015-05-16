package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

const vmrun = "/Applications/VMware Fusion.app/Contents/Library/vmrun"

var (
	home = os.Getenv("HOME")
	vmdir = fmt.Sprintf("%s/%s", home, "Documents/Virtual Machines.localized")
)

func vmname(vmpath string) string {
	relpath := strings.TrimPrefix(vmpath, vmdir)
	vmwarevm := path.Base(path.Dir(relpath))
	return strings.TrimSuffix(vmwarevm, ".vmwarevm")
}

func list() (map[string]string) {
	glob := fmt.Sprintf("%s/**/*.vmx", vmdir)
	paths, err := filepath.Glob(glob)
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}

	vmx := make(map[string]string)
	for _, path := range paths {
		name := vmname(path)
		vmx[name] = path
	}

	return vmx
}

func main() {
	args := os.Args[1:]
	out, err := exec.Command(vmrun, args...).Output()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}
	fmt.Printf("%s", out)

	if len(args) > 0 && args[0] == "list" {
		fmt.Printf("\nVMs in %s:\n", vmdir)
		for name, path := range list() {
			fmt.Printf("  %+20s (%s)\n", name, path)
		}
	}
}
