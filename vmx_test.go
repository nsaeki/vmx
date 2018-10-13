package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func Setup(t *testing.T) (string, func()) {
	vmdir, err := ioutil.TempDir("", "vmxtest")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("Create %s", vmdir))
	setVMDir(vmdir)

	return vmdir, func() {
		t.Log(fmt.Sprintf("Remove %s", vmdir))
		if err := os.RemoveAll(vmdir); err != nil {
			t.Fatal(err)
		}
	}
}

func CreateVM(vmdir, vmname, vmx string) {
	os.MkdirAll(path.Join(vmdir, vmname), os.ModePerm)
	if vmx != "" {
		os.Create(path.Join(vmdir, vmname, vmx))
	}
}

func TestListVMs(t *testing.T) {
	vmdir, f := Setup(t)
	defer f()
	CreateVM(vmdir, "Debian.vmwarevm", "Debian.vmx")
	CreateVM(vmdir, "Debian (stretch).vmwarevm", "Debian.vmx")
	CreateVM(vmdir, "Empty", "")

	vms := listVMs()
	expected := 2
	if len(vms) != expected {
		t.Errorf("Number of VMs is %d: expected %d", len(vms), expected)
	}

	if _, ok := vms["Debian"]; !ok {
		t.Errorf("Some VMs has not been found: VM name = %s", "Debian")
	}
}

func TestConvertArgs(t *testing.T) {
	cases := []struct {
		in, want []string
	}{
		{
			[]string{""},
			[]string{""},
		},
		{
			[]string{"list"},
			[]string{"list"},
		},
		{
			[]string{"start", "Debian"},
			[]string{"start", "%", "nogui"},
		},
		{
			[]string{"start", "Debian", "gui"},
			[]string{"start", "%", "gui"},
		},
		{
			[]string{"-T", "ws", "start", "Debian"},
			[]string{"-T", "ws", "start", "%", "nogui"},
		},
		{
			[]string{"-T", "ws", "stop", "-u", "user", "Debian", "-p", "pass"},
			[]string{"-T", "ws", "stop", "-u", "user", "%", "-p", "pass"},
		},
	}

	for _, c := range cases {
		got := convertArgs(c.in)
		ok := true
		for i := 0; i < len(got); i++ {
			if c.want[i] != "%" && c.want[i] != got[i] {
				ok = false
				break
			}
		}
		if !ok {
			t.Errorf("input: %v\ngot: %v\n expected: %v", c.in, got, c.want)
		}
	}
}
