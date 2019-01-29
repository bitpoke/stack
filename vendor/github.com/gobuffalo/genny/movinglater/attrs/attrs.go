package attrs

import (
	"strings"

	"github.com/gobuffalo/flect/name"
)

//Attr is buffalo's implementation for model attributes
type Attr struct {
	Original   string
	Name       name.Ident
	commonType string
	goType     string
}

func (a Attr) String() string {
	return a.Original
}

//GoType returns the Go type for an Attr based on its commonType
func (a Attr) GoType() string {
	if a.goType != "" {
		return a.goType
	}

	switch strings.ToLower(a.commonType) {
	case "text":
		return "string"
	case "timestamp", "datetime", "date", "time":
		return "time.Time"
	case "nulls.text":
		return "nulls.String"
	case "uuid":
		return "uuid.UUID"
	case "json", "jsonb":
		return "slices.Map"
	case "[]string":
		return "slices.String"
	case "[]int":
		return "slices.Int"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "slices.Float"
	case "decimal", "float":
		return "float64"
	case "[]byte", "blob":
		return "[]byte"
	}

	return a.commonType
}

//CommonType returns the common type of an attribute,
//this common type is used later for things like determining
//the database column type depending on the database.
func (a Attr) CommonType() string {
	return commonType(a.commonType)
}

func commonType(s string) string {
	switch strings.ToLower(s) {
	case "int":
		return "integer"
	case "time", "datetime":
		return "timestamp"
	case "uuid.uuid", "uuid":
		return "uuid"
	case "nulls.float32", "nulls.float64":
		return "float"
	case "slices.string", "slices.uuid", "[]string":
		return "[]string"
	case "slices.float", "[]float", "[]float32", "[]float64":
		return "[]float"
	case "slices.int":
		return "[]int"
	case "slices.map":
		return "json"
	case "float32", "float64", "float":
		return "decimal"
	case "blob", "[]byte":
		return "[]byte"
	default:
		if strings.HasPrefix(s, "nulls.") {
			return commonType(strings.Replace(s, "nulls.", "", -1))
		}
		return strings.ToLower(s)
	}
}

//Attrs is a slice of Attr
type Attrs []Attr

func (ats Attrs) Slice() []string {
	var x []string
	for _, a := range ats {
		x = append(x, a.Original)
	}
	return x
}
