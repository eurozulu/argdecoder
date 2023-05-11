package argdecoder

import (
	"fmt"
	"reflect"
	"strings"
)

type mapDecoder struct {
	args []string
}

func (md mapDecoder) Apply(v interface{}) ([]string, error) {
	mapValue := reflect.ValueOf(v)
	if mapValue.Kind() != reflect.Map {
		return nil, fmt.Errorf("can not decode into non map value")
	}
	if mapValue.Type().Key().Kind() != reflect.String {
		return nil, fmt.Errorf("can not decode into maps without string keymanager")
	}

	if mapValue.Type().Elem().Kind() != reflect.Interface {
		return nil, fmt.Errorf("can not decode into maps without interface values")
	}

	params, flags := parseArgs(md.args)
	m := map[string]interface{}{}
	if len(params) > 0 {
		m[""] = params
	}
	for f, sv := range flags {
		m[f] = sv
	}
	mapValue.Set(reflect.ValueOf(m))
	return nil, nil
}

func parseArgs(args []string) (params []string, flags map[string]*string) {
	flags = map[string]*string{}
	for index := 0; index < len(args); index++ {
		if !strings.HasPrefix(args[index], "-") {
			params = append(params, args[index])
			continue
		}
		flag := strings.ToLower(strings.TrimLeft(args[index], "-"))
		var value *string
		if index+1 < len(args) && !strings.HasPrefix(args[index+1], "-") {
			index++
			value = &args[index]
		}
		flags[flag] = value
	}
	return params, flags
}
