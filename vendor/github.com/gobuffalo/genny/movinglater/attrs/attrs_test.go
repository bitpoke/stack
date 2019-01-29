package attrs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Attr_GoType(t *testing.T) {

	tt := []struct {
		ct string
		gt string
	}{
		{"timestamp", "time.Time"},
		{"datetime", "time.Time"},
		{"date", "time.Time"},
		{"time", "time.Time"},
		{"text", "string"},
		{"Text", "string"},
		{"nulls.text", "nulls.String"},
		{"nulls.Text", "nulls.String"},
		{"uuid", "uuid.UUID"},
		{"json", "slices.Map"},
		{"jsonb", "slices.Map"},
		{"[]string", "slices.String"},
		{"[]int", "slices.Int"},
		{"slices.float", "slices.Float"},
		{"[]float", "slices.Float"},
		{"[]float32", "slices.Float"},
		{"[]float64", "slices.Float"},
		{"decimal", "float64"},
		{"float", "float64"},
		{"[]byte", "[]byte"},
		{"blob", "[]byte"},
	}

	for _, test := range tt {
		t.Run(test.ct+"/"+test.gt, func(st *testing.T) {
			r := require.New(st)
			a := Attr{commonType: test.ct}
			r.Equal(test.gt, a.GoType())
		})
	}
}

func Test_Attr_CommonType(t *testing.T) {

	tt := []struct {
		pt string
		ct string
	}{
		{"timestamp", "timestamp"},
		{"datetime", "timestamp"},
		{"date", "date"},
		{"time", "timestamp"},
		{"text", "text"},
		{"Text", "text"},
		{"nulls.text", "text"},
		{"nulls.Text", "text"},
		{"uuid", "uuid"},
		{"slices.Map", "json"},
		{"slices.String", "[]string"},
		{"slices.Int", "[]int"},
		{"slices.float", "[]float"},
		{"slices.Float", "[]float"},
		{"[]float64", "[]float"},
		{"float64", "decimal"},
		{"float", "decimal"},
		{"[]byte", "[]byte"},
	}

	for _, test := range tt {
		t.Run(test.pt+"/"+test.ct, func(st *testing.T) {
			r := require.New(st)
			a := Attr{commonType: test.pt}
			r.Equal(test.ct, a.CommonType())
		})
	}

}
