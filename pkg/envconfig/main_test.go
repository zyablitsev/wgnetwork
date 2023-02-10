package envconfig

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

type envValue struct{ checkValue, configValue interface{} }

type envValues map[string]envValue

func (env envValues) set() error {
	for k, v := range env {
		envString := parseValue(reflect.ValueOf(v.checkValue))
		if err := os.Setenv(k, envString); err != nil {
			return err
		}
	}
	return nil
}

func parseValue(v reflect.Value) string {
	value := ""
	switch v.Kind() {
	case reflect.Bool:
		value = strconv.FormatBool(v.Bool())
	case reflect.String:
		value = v.String()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = strconv.FormatUint(v.Uint(), 10)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value = strconv.FormatInt(v.Int(), 10)
	case reflect.Float32:
		value = strconv.FormatFloat(v.Float(), 'E', -1, 32)
	case reflect.Float64:
		value = strconv.FormatFloat(v.Float(), 'E', -1, 64)
	case reflect.Array, reflect.Slice:
		slice := make([]string, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			s := parseValue(v.Index(i))
			slice = append(slice, s)
		}
		value = strings.Join(slice, ",")
	case reflect.Map:
		slice := make([]string, 0, v.Len())
		mapKeys := v.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			s := make([]string, 0, 2)
			s = append(s, parseValue(mapKeys[i]))
			s = append(s, parseValue(v.MapIndex(mapKeys[i])))
			slice = append(slice, strings.Join(s, ":"))
		}
		value = strings.Join(slice, ",")
	}
	return value
}
