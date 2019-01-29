package attrs

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ParseNamedArgs(t *testing.T) {
	r := require.New(t)

	na, err := ParseNamedArgs("widget", "name", "age:int")
	r.NoError(err)

	r.Equal("widget", na.Name.String())
	r.Len(na.Attrs, 2)
}

func Test_ParseNamedArgs_NoName(t *testing.T) {
	r := require.New(t)

	_, err := ParseNamedArgs()
	r.Error(err)

}

func Test_NamedAttrs_String(t *testing.T) {
	table := []struct {
		in  []string
		out string
	}{
		{[]string{"foo", "bar:baz"}, "foo bar:baz"},
	}

	for _, tt := range table {
		t.Run(strings.Join(tt.in, " "), func(st *testing.T) {
			r := require.New(st)
			n, err := ParseNamedArgs(tt.in...)
			r.NoError(err)
			r.Equal(tt.out, n.String())
		})
	}
}
