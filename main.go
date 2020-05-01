package confstruct

import (
	"fmt"
	"reflect"
	"strings"
)

const tagName = "conf"

// Populate takes an annotated struct, and populates
// the fields with the corresponding values from the
// environment variables
func Populate(cs interface{}) error {
	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(cs)

	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	} else {
		return fmt.Errorf("argument must be pointer to struct")
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("can't populate non-struct type")
	}

	for i := 0; i < v.NumField(); i++ {
		// Get the field tag value
		tag := v.Type().Field(i).Tag.Get(tagName)

		// Ignore empty/unset
		if tag == "" || tag == "-" {
			continue
		}

		// Get the type of the field
		var t string
		var ptr bool
		if v.Type().Field(i).Type.Kind() == reflect.Ptr {
			// it's a pointer, use underlying type
			t = v.Type().Field(i).Type.Elem().Name()
			ptr = true
		} else {
			// no pointer, straight type
			t = v.Type().Field(i).Type.Name()
			ptr = false
		}

		f, err := getFetcher(tag, t, ptr)
		if err != nil {
			return err
		}
		if v.Field(i).CanAddr() && v.Field(i).CanSet() {
			rval, err := f.Fetch()
			if err != nil {
				return err
			}
			// Set the value
			v.Field(i).Set(rval)
		}
	}

	return nil
}

// Fetcher gets values from environment variables
type Fetcher interface {
	// Fetch a value
	Fetch() (reflect.Value, error)
}

// Returns Fetcher struct corresponding to field's datatype
func getFetcher(tag string, t string, ptr bool) (Fetcher, error) {
	var defVal, format string
	hasDef := false
	args := strings.Split(tag, ",")
	varName := strings.TrimSpace(args[0])
	if len(args) > 1 {
		// yes there's options
		for i := 1; i < len(args); i++ {
			argParts := strings.Split(args[i], "=")
			switch strings.TrimSpace(argParts[0]) {
			case "default":
				defVal = argParts[1]
				hasDef = true
			case "format":
				format = argParts[1]
			}
		}
	}

	switch t {
	case "string":
		return StringFetcher{
			VarName: varName,
			IsPtr:   ptr,
			HasDef:  hasDef,
			Default: defVal,
		}, nil
	case "int":
		return IntFetcher{
			VarName: varName,
			IsPtr:   ptr,
			HasDef:  hasDef,
			Default: defVal,
		}, nil
	case "bool":
		return BoolFetcher{
			VarName: varName,
			IsPtr:   ptr,
			HasDef:  hasDef,
			Default: defVal,
		}, nil
	case "float64":
		return Float64Fetcher{
			VarName: varName,
			IsPtr:   ptr,
			HasDef:  hasDef,
			Default: defVal,
		}, nil
	case "Time":
		return TimeFetcher{
			VarName: varName,
			IsPtr:   ptr,
			HasDef:  hasDef,
			Default: defVal,
			Format:  format,
		}, nil
	case "Duration":
		return DurationFetcher{
			VarName: varName,
			IsPtr:   ptr,
			HasDef:  hasDef,
			Default: defVal,
		}, nil
	case "URL":
		return URLFetcher{
			VarName: varName,
			IsPtr:   ptr,
			HasDef:  hasDef,
			Default: defVal,
		}, nil
	}
	// Dunno...
	return nil, fmt.Errorf("confstruct does not (yet) support variables of type: %s", t)
}
