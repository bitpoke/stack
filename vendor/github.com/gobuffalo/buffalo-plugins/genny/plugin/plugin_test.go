package plugin

import (
	"strings"
	"testing"

	"github.com/gobuffalo/genny/gentest"
	"github.com/gobuffalo/genny/movinglater/gotools/gomods"
	"github.com/stretchr/testify/require"
)

func Test_Generator(t *testing.T) {
	r := require.New(t)
	opts := &Options{
		PluginPkg: "github.com/foo/buffalo-bar",
		Year:      1999,
		Author:    "Homer Simpson",
		ShortName: "bar",
	}

	run := gentest.NewRunner()

	gg, err := New(opts)
	r.NoError(err)
	run.WithGroup(gg)

	r.NoError(run.Run())
	res := run.Results()

	var cmds []string
	if !gomods.On() {
		cmds = []string{"git init", "go get github.com/alecthomas/gometalinter", "gometalinter --install"}
	} else {
		cmds = []string{"git init", "go mod init github.com/foo/buffalo-bar", "go get github.com/alecthomas/gometalinter", "gometalinter --install", "go mod tidy"}
	}

	r.Len(res.Commands, len(cmds))
	for i, x := range cmds {
		r.Equal(x, strings.TrimSpace(strings.Join(res.Commands[i].Args, " ")))
	}

	files := []string{
		".gitignore",
		".gometalinter.json",
		".goreleaser.yml.plush",
		".travis.yml",
		"LICENSE",
		"Makefile",
		"README.md",
		"bar/listen.go",
		"bar/version.go",
		"cmd/available.go",
		"cmd/bar.go",
		"cmd/root.go",
		"cmd/version.go",
		"main.go",
	}
	r.NoError(gentest.CompareFiles(files, res.Files))

	f, err := res.Find("README.md")
	r.NoError(err)
	r.Contains(f.String(), opts.PluginPkg)

	f, err = res.Find("cmd/version.go")
	r.NoError(err)
	r.Contains(f.String(), opts.PluginPkg+"/"+opts.ShortName)
	r.Contains(f.String(), opts.ShortName+".Version")

	f, err = res.Find("main.go")
	r.NoError(err)
	r.Contains(f.String(), "github.com/foo/buffalo-bar/cmd")

}
