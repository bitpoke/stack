package gentest

import (
	"fmt"
	"os/exec"
	"strings"
)

// CompareCommands ...
func CompareCommands(exp []string, act []*exec.Cmd) error {
	if len(exp) != len(act) {
		return fmt.Errorf("len(exp) != len(act) [%d != %d]", len(exp), len(act))
	}
	for i, c := range act {
		e := exp[i]
		a := strings.Join(c.Args, " ")
		if a != e {
			return fmt.Errorf("expect %q got %q", a, e)
		}
	}
	return nil
}
