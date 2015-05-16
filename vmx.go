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

func extractVMName(vmpath string) string {
	relpath := strings.TrimPrefix(vmpath, vmdir)
	vmwarevm := path.Base(path.Dir(relpath))
	return strings.TrimSuffix(vmwarevm, ".vmwarevm")
}

func listVMs() (map[string]string) {
	glob := fmt.Sprintf("%s/**/*.vmx", vmdir)
	paths, err := filepath.Glob(glob)
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}

	vmx := make(map[string]string)
	for _, path := range paths {
		name := extractVMName(path)
		vmx[name] = path
	}

	return vmx
}

// Returns vmx file path from name.
// If name is vmx file path itself, returns that unchanged.
// If vmx path is not found, returns empty string.
func findVmxPath(name string) string {
	path := name
	if _, err := os.Stat(name); os.IsNotExist(err) {
		// If calling listVMs is heavy, consider to cache this value.
		path = listVMs()[name]
	}
	return path
}

// Convert VM name to vmx path in original args.
func convertArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	originalArgs := args
	args = make([]string, len(args))
	copy(args, originalArgs)

	var command, vmName, guiMode string
	optionArgs := false
	for i, v := range args {
		switch {
		case optionArgs:
			optionArgs = false
		case strings.HasPrefix(v, "-"):
			optionArgs = true
		case command == "":
			command = v
		case vmName == "":
			vmName = v
			args[i] = findVmxPath(v)
		case v == "gui" || v == "nogui":
			guiMode = v
		}
	}
	if command == "start" && guiMode == "" {
		args = append(args, "nogui")
	}
	return args
}


func main() {
	args := convertArgs(os.Args[1:])
	cmd := exec.Command(vmrun, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}
	fmt.Printf("%s", out)

	if len(args) > 0 && args[0] == "list" {
		fmt.Printf("\nVMs in %s:\n", vmdir)
		vms := listVMs()
		maxNameLen := 0
		for name, _ := range vms {
			l := len(name)
			if maxNameLen < l {
				maxNameLen = l
			}
		}

		format := fmt.Sprintf("  %%-%ds (%%s)\n", maxNameLen + 2)
		for name, path := range vms {
			fmt.Printf(format, name, path)
		}
	}

	if !cmd.ProcessState.Success() {
		os.Exit(1)
	}
}
