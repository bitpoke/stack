package plugin

import (
	"os/user"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Options_Validate(t *testing.T) {
	r := require.New(t)

	opts := &Options{}
	err := opts.Validate()
	r.Error(err)

	opts.PluginPkg = "github.com/foo/bar"

	err = opts.Validate()
	r.NoError(err)
	r.Equal("github.com/foo/buffalo-bar", opts.PluginPkg)

	year := time.Now().Year()
	r.Equal(opts.Year, year)

	u, err := user.Current()
	r.NoError(err)
	r.Equal(u.Name, opts.Author)
}
