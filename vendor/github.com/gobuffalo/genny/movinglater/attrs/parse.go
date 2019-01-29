package attrs

import (
	"strings"

	"github.com/gobuffalo/flect/name"
	"github.com/pkg/errors"
)

//ErrRepeatedAttr is returned when parsing an array with repeated names
var ErrRepeatedAttr = errors.New("duplicate attr name")

//Parse takes a string like name:commonType:goType and turns it into an Attr
func Parse(arg string) (Attr, error) {
	arg = strings.TrimSpace(arg)
	attr := Attr{
		Original:   arg,
		commonType: "string",
	}
	if len(arg) == 0 {
		return attr, errors.New("argument can not be blank")
	}

	parts := strings.Split(arg, ":")
	attr.Name = name.New(parts[0])
	if len(parts) > 1 {
		attr.commonType = parts[1]
	}

	if len(parts) > 2 {
		attr.goType = parts[2]
	}

	return attr, nil
}

//ParseArgs parses passed string args into Attrs
func ParseArgs(args ...string) (Attrs, error) {
	var attrs Attrs
	parsed := map[string]string{}

	for _, arg := range args {
		a, err := Parse(arg)
		if err != nil {
			return attrs, errors.WithStack(err)
		}

		key := a.Name.Underscore().String()
		if parsed[key] != "" {
			return attrs, ErrRepeatedAttr
		}

		parsed[key] = arg
		attrs = append(attrs, a)
	}

	return attrs, nil
}
