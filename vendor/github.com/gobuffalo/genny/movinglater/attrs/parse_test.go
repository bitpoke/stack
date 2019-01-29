package attrs

import (
	"testing"

	"github.com/gobuffalo/flect/name"
	"github.com/stretchr/testify/require"
)

func Test_Parse(t *testing.T) {
	attrs := Attrs{
		{Original: "foo", goType: "string", commonType: "string", Name: name.New("foo")},
		{Original: "foo:int", goType: "int", commonType: "integer", Name: name.New("foo")},
		{Original: "foo:timestamp", goType: "time.Time", commonType: "timestamp", Name: name.New("foo")},
		{Original: "foo:text:exec.Command", goType: "exec.Command", commonType: "text", Name: name.New("foo")},
	}

	for _, a := range attrs {
		t.Run(a.Original, func(st *testing.T) {
			r := require.New(st)
			attr, err := Parse(a.Original)
			r.NoError(err)
			r.Equal(a.Original, attr.Original)
			r.Equal(a.goType, attr.GoType())
			r.Equal(a.commonType, attr.CommonType())
			r.Equal(a.Name, attr.Name)
		})
	}
}

func Test_ParseArgs(t *testing.T) {
	tcases := []struct {
		name string
		args []string
		err  error
	}{
		{
			name: "All Good!",
			args: []string{
				"name",
				"profile:text",
			},
			err: nil,
		},

		{
			name: "Repeated arg",
			args: []string{
				"name",
				"profile:text",
				"profile",
			},
			err: ErrRepeatedAttr,
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.name, func(st *testing.T) {
			r := require.New(st)
			_, err := ParseArgs(tcase.args...)

			if tcase.err != nil {
				r.NotNil(err)
				r.Equal(tcase.err.Error(), err.Error())
			}
		})
	}
}
