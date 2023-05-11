package argdecoder

import (
	"fmt"
	"reflect"
)

type ArgumentDecoder interface {
	// Apply will apply the arguments to the given object.
	// Arguments are parsed into flags, args beginning with '-', with the remaining arg value defining the flag name.
	// this name is matched to a public field in the given value either directly by name or via a `flag` tag.
	// If the arguments has another arg following the flag, and not starting with -, this is used as the value for that flag
	// which is assigned to the Field, assuming it can be coerced into the relevant type, otherwise an error is thrown.
	// Any arguments not matched to fields are returned.  As such, multiple objects may be passed to the same decoder, each "consuming" their flags.
	Apply(v interface{}) ([]string, error)
}

func ApplyArguments(args []string, v interface{}) ([]string, error) {
	vv := reflect.ValueOf(v)
	if vv.IsNil() {
		return nil, fmt.Errorf("can not apply to nil value")
	}

	switch vv.Elem().Kind() {
	case reflect.Pointer:
		return ApplyArguments(args, vv.Elem())
	case reflect.Struct:
		return structParser{args: args}.Apply(v)
	case reflect.Map:
		return mapDecoder{args: args}.Apply(v)
	case reflect.Slice:
		return stringSliceDecoder{args: args}.Apply(v)
	case reflect.String:
		return stringDecoder{args: args}.Apply(v)
	default:
		return nil, fmt.Errorf("%s is not ap supported type", vv.Type().String())
	}
}
