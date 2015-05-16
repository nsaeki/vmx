package main

import (
	"testing"
)

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
