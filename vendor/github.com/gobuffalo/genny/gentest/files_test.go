package gentest

import (
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_CompareFiles(t *testing.T) {
	r := require.New(t)

	exp := []string{"b.html", "a.html"}
	act := []genny.File{
		genny.NewFileS("a.html", "A"),
		genny.NewFileS("b.html", "B"),
	}
	r.NoError(CompareFiles(exp, act))
}
