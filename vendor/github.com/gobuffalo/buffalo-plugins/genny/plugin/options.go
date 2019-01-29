package plugin

import (
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/markbates/going/defaults"
	"github.com/pkg/errors"
)

type Options struct {
	PluginPkg string
	Year      int
	Author    string
	ShortName string
	License   string
	Root      string
}

func (opts *Options) Validate() error {
	if opts.PluginPkg == "" {
		return errors.New("plugin has to have a package name")
	}
	if len(opts.License) == 0 {
		opts.License = "mit"
	}

	name := path.Base(opts.PluginPkg)
	opts.ShortName = strings.TrimPrefix(name, "buffalo-")
	if !strings.HasPrefix(name, "buffalo-") {
		name = "buffalo-" + name
	}

	dir := path.Dir(opts.PluginPkg)
	opts.PluginPkg = path.Join(dir, name)
	opts.Root = filepath.Join(opts.Root, opts.PluginPkg)

	if opts.Year == 0 {
		opts.Year = time.Now().Year()
	}

	if len(opts.Author) == 0 {
		u, err := user.Current()
		if err != nil {
			return errors.WithStack(err)
		}
		opts.Author = defaults.String(u.Name, u.Username)
	}

	return nil
}
