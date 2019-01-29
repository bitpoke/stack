package dep

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/gentest"
	"github.com/stretchr/testify/require"
)

func Test_Init_WithDep(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	run.LookPathFn = func(s string) (string, error) {
		if s == "dep" {
			return "dep", nil
		}
		return exec.LookPath(s)
	}

	err := run.WithNew(Init("", false))
	r.NoError(err)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 1)

	c := res.Commands[0]
	r.Equal("dep init", strings.Join(c.Args, " "))
}

func Test_Init_WithoutDep(t *testing.T) {
	r := require.New(t)

	run := gentest.NewRunner()
	run.LookPathFn = func(s string) (string, error) {
		if s == "dep" {
			return "", os.ErrNotExist
		}
		return exec.LookPath(s)
	}

	err := run.WithNew(Init("", false))
	r.NoError(err)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 2)

	c := res.Commands[0]
	r.Equal(genny.GoBin()+" get github.com/golang/dep/cmd/dep", strings.Join(c.Args, " "))

	c = res.Commands[1]
	r.Equal("dep init", strings.Join(c.Args, " "))
}
