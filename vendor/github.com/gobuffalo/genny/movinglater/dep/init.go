package dep

import (
	"os/exec"

	"github.com/gobuffalo/genny"
)

func Init(path string, verbose bool) (*genny.Generator, error) {
	g := genny.New()
	var args []string
	if verbose {
		args = append(args, "-v")
	}

	g.RunFn(InstallDep(args...))
	cmd := exec.Command("dep", "init")
	if verbose {
		cmd.Args = append(cmd.Args, args...)
	}
	g.Command(cmd)
	return g, nil
}
