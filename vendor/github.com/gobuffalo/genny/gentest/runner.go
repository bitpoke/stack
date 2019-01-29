package gentest

import (
	"context"

	"github.com/gobuffalo/genny"
	"github.com/pkg/errors"
)

// NewRunner is a dry runner with gentest.NewLogger()
func NewRunner() *genny.Runner {
	r := genny.DryRunner(context.Background())
	r.Logger = NewLogger()
	return r
}

// Run the generator and return results or an error
func Run(g *genny.Generator) (genny.Results, error) {
	return RunNew(g, nil)
}

// RunNew the generator and return results or an error
func RunNew(g *genny.Generator, err error) (genny.Results, error) {
	if err != nil {
		return genny.Results{}, errors.WithStack(err)
	}

	r := NewRunner()
	r.With(g)

	return sprint(r)
}

// RunGroup runs the group and returns results or an error
func RunGroup(gg *genny.Group) (genny.Results, error) {
	r := NewRunner()
	r.WithGroup(gg)
	return sprint(r)
}

func sprint(r *genny.Runner) (genny.Results, error) {
	if err := r.Run(); err != nil {
		return r.Results(), errors.WithStack(err)
	}
	return r.Results(), nil
}
