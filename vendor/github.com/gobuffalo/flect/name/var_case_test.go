package name

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_VarCaseSingle(t *testing.T) {
	table := []tt{
		{"foo_bar", "fooBar"},
		{"admin/widget", "adminWidget"},
		{"widget", "widget"},
		{"widgets", "widget"},
		{"User", "user"},
		{"FooBar", "fooBar"},
		{"status", "status"},
		{"statuses", "status"},
		{"Status", "status"},
		{"Statuses", "status"},
	}

	for _, tt := range table {
		t.Run(tt.act, func(st *testing.T) {
			r := require.New(st)
			r.Equal(tt.exp, VarCaseSingle(tt.act))
			r.Equal(tt.exp, VarCaseSingle(tt.exp))
		})
	}
}

func Test_VarCasePlural(t *testing.T) {
	table := []tt{
		{"foo_bar", "fooBars"},
		{"admin/widget", "adminWidgets"},
		{"widget", "widgets"},
		{"widgets", "widgets"},
		{"User", "users"},
		{"FooBar", "fooBars"},
		{"status", "statuses"},
		{"statuses", "statuses"},
		{"Status", "statuses"},
		{"Statuses", "statuses"},
	}

	for _, tt := range table {
		t.Run(tt.act, func(st *testing.T) {
			r := require.New(st)
			r.Equal(tt.exp, VarCasePlural(tt.act))
			r.Equal(tt.exp, VarCasePlural(tt.exp))
		})
	}
}
