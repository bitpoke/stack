package git

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

var ErrWorkingTreeClean = errors.New("working tree clean")

func Run(args ...string) genny.RunFn {
	return func(r *genny.Runner) error {
		cmd := exec.Command("git", args...)
		err := r.Exec(cmd)
		if err != nil {
			if workingDirIsClean() {
				return ErrWorkingTreeClean
			}
			return errors.WithStack(err)
		}
		return nil
	}

}

func workingDirIsClean() bool {
	bb := &bytes.Buffer{}
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Stdout = bb
	err := cmd.Run()
	if err != nil {
		return false
	}
	return strings.TrimSpace(bb.String()) == ""
}
