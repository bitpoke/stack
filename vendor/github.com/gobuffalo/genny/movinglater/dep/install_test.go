package dep

import (
	"errors"
	"testing"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_InstallDep_NotInstalled(t *testing.T) {
	envy.Set("GO_BIN", "go")
	r := require.New(t)

	run := gentest.NewRunner()
	run.LookPathFn = func(s string) (string, error) {
		return s, errors.New("couldn't find dep")
	}
	run.WithRun(InstallDep())
	r.NoError(run.Run())

	res := run.Results()

	cmds := []string{"go get github.com/golang/dep/cmd/dep"}
	r.Len(res.Commands, len(cmds))
}

func Test_InstallDep_Installed(t *testing.T) {
	envy.Set("GO_BIN", "go")
	r := require.New(t)

	run := gentest.NewRunner()
	run.LookPathFn = func(s string) (string, error) {
		return s, nil
	}
	run.WithRun(InstallDep())
	r.NoError(run.Run())

	res := run.Results()
	r.Len(res.Commands, 0)
}
