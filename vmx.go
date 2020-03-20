package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var (
	vmrun  = "/Applications/VMware Fusion.app/Contents/Library/vmrun"
	home   = os.Getenv("HOME")
	vmdir  = fmt.Sprintf("%s/%s", home, "Virtual Machines.localized")
	vmdirs = []string{
		fmt.Sprintf("%s/%s", home, "Virtual Machines.localized"),
		fmt.Sprintf("%s/%s", home, "Documents/Virtual Machines.localized"),
	}
)

// Changes VM Directory. Currently this function only for testing.
func setVMDir(dir string) {
	vmdir = dir
}

func setVMDirs(dirs []string) {
	vmdirs = dirs
}

func findVMDir() string {
	for _, d := range vmdirs {
		fi, err := os.Stat(d)
		if err == nil && fi.IsDir() {
			return d
		}
	}
	return ""
}

func extractVMName(vmxpath string) string {
	relpath := strings.TrimPrefix(vmxpath, vmdir)
	vmwarevm := path.Base(path.Dir(relpath))
	return strings.TrimSuffix(vmwarevm, ".vmwarevm")
}

func listVMs() map[string]string {
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
// If vmx path is not found or two or more vmx found, returns empty string and non-nil error
func findVmxPath(name string) (string, error) {
	path := name
	if _, err := os.Stat(name); os.IsNotExist(err) {
		glob := fmt.Sprintf("%s/%s*/*.vmx", vmdir, name)
		paths, err := filepath.Glob(glob)
		if err != nil {
			return "", err
		} else if paths == nil {
			return "", fmt.Errorf("No VM image found like: %s", name)
		} else if len(paths) > 1 {
			return "", fmt.Errorf("Two or more VM images found like: %s", name)
		} else {
			path = paths[0]
		}
	}
	return path, nil
}

// Converts VM name to vmx path.
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
			vmx, err := findVmxPath(v)
			if vmx == "" || err != nil {
				if err != nil {
					fmt.Print(err)
				}
				list()
				os.Exit(2)
			}
			args[i] = vmx
		case v == "gui" || v == "nogui":
			guiMode = v
		}
	}
	if command == "start" && guiMode == "" {
		args = append(args, "nogui")
	}
	return args
}

func list() {
	fmt.Printf("\nVMs in %s:\n", vmdir)
	vms := listVMs()
	maxNameLen := 0
	for name, _ := range vms {
		l := len(name)
		if maxNameLen < l {
			maxNameLen = l
		}
	}

	format := fmt.Sprintf("  %%-%ds (%%s)\n", maxNameLen+2)
	for name, path := range vms {
		fmt.Printf(format, name, path)
	}
}

func main() {
	findVMDir()
	args := convertArgs(os.Args[1:])

	if len(args) > 0 && args[0] == "list" {
		list()
		return
	}

	cmd := exec.Command(vmrun, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s\n\n", err.Error())
	}
	fmt.Printf("%s", out)

	if !cmd.ProcessState.Success() {
		os.Exit(1)
	}
}
