package argdecoder

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const TagKey = "flag"

var errUnknownField = fmt.Errorf("not a known field")

var errBoolFailedToParse = fmt.Errorf("can not parse as bool")

type structParser struct {
	args []string
}

func (sd structParser) Apply(v interface{}) ([]string, error) {
	receiverValue := reflect.ValueOf(v)
	if receiverValue.Kind() != reflect.Ptr ||
		receiverValue.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("value must be a pointer to a struct")
	}

	var unused []string
	for index := 0; index < len(sd.args); index++ {
		arg := sd.args[index]
		if !isFlagArg(arg) {
			unused = append(unused, arg)
			continue
		}
		// arg is a flag, look for next arg as value
		name := strings.TrimLeft(arg, "-")
		var value *string
		if index+1 < len(sd.args) && !isFlagArg(sd.args[index+1]) {
			index++
			value = &sd.args[index]
		}
		fld, err := fieldForName(name, receiverValue.Elem())
		if err != nil {
			// No matching field for this flag
			unused = append(unused, arg)
			if value != nil {
				unused = append(unused, *value)
			}
			continue
		}

		typedValue, err := stringAsValue(value, fld.Type)
		if err != nil {
			if err != errBoolFailedToParse {
				return nil, err
			}
			// if bool parse error assume it's a parameter following a nil bool flag.
			// step back and process next arg as regular arg, not value
			index--
			typedValue = reflect.ValueOf(true)
		}
		receiverValue.Elem().FieldByIndex(fld.Index).Set(typedValue)

	}
	return unused, nil
}

func fieldForName(name string, receiver reflect.Value) (reflect.StructField, error) {
	t := receiver.Type()
	for i := 0; i < t.NumField(); i++ {
		fld := t.Field(i)
		if !strings.EqualFold(name, fld.Name) &&
			!isTagName(name, fld.Tag.Get(TagKey)) {
			continue
		}
		if !fld.IsExported() {
			return fld, fmt.Errorf("field %s is not an exported field", fld.Name)
		}
		return fld, nil
	}
	return reflect.StructField{}, errUnknownField
}

func stringAsValue(svalue *string, vtype reflect.Type) (reflect.Value, error) {
	if svalue == nil {
		switch vtype.Kind() {
		case reflect.Bool:
			// bool defaults to true when no value given
			return reflect.ValueOf(true), nil
		case reflect.Pointer, reflect.Slice:
			return reflect.New(vtype).Elem(), nil
		default:
			return reflect.Value{}, fmt.Errorf("no value given")
		}
	}
	s := strings.TrimSpace(*svalue)
	switch vtype.Kind() {
	case reflect.Pointer:
		return pointerValueForString(svalue, vtype)

	case reflect.Slice:
		return sliceValueForString(svalue, vtype)

	case reflect.String:
		return reflect.ValueOf(s).Convert(vtype), nil

	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, errBoolFailedToParse
		}
		return reflect.ValueOf(b).Convert(vtype), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into an uint value  %v", s, err)
		}
		return reflect.ValueOf(ui).Convert(vtype), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into a int value  %v", s, err)
		}
		return reflect.ValueOf(i).Convert(vtype), nil

	case reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into a floatvalue  %v", s, err)
		}
		return reflect.ValueOf(f).Convert(vtype), nil

	case reflect.Float32:
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Failed to convert %s into a floatvalue  %v", s, err)
		}
		return reflect.ValueOf(f).Convert(vtype), nil

	default:
		return reflect.Value{}, fmt.Errorf("%s type is not supported", vtype.Kind().String())
	}
}

func isFlagArg(arg string) bool {
	return strings.HasPrefix(arg, "-")
}

func isTagName(name, tagValue string) bool {
	for _, tn := range strings.Split(tagValue, ",") {
		if strings.EqualFold(strings.TrimSpace(tn), name) {
			return true
		}
	}
	return false
}

func pointerValueForString(svalue *string, vtype reflect.Type) (reflect.Value, error) {
	v, err := stringAsValue(svalue, vtype.Elem())
	if err != nil {
		return v, err
	}
	vp := reflect.New(vtype.Elem())
	vp.Elem().Set(v)
	return vp, nil
}

func sliceValueForString(sp *string, vtype reflect.Type) (reflect.Value, error) {
	if sp == nil {
		return reflect.New(vtype), nil
	}
	ss := strings.Split(*sp, ",")
	el := vtype.Elem()
	sv := reflect.MakeSlice(vtype, len(ss), len(ss))
	for i, sz := range ss {
		v, err := stringAsValue(&sz, el)
		if err != nil {
			return v, err
		}
		sv.Index(i).Set(v)
	}
	return sv, nil
}
