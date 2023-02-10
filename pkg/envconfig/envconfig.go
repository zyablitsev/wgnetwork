package envconfig

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"time"
)

// ReadEnv get a config struct ptr
// and populate it fields with values from environment varibales
func ReadEnv(config interface{}) error {
	err := handleStruct(reflect.ValueOf(config), handleField)
	return err
}

func handlePtr(v reflect.Value) (reflect.Value, error) {
	if v.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("handlePtr: value isn't a pointer")
	}
	return v.Elem(), nil
}

func handleStruct(v reflect.Value, f func(reflect.StructField, reflect.Value)) error {
	if v.Kind() == reflect.Ptr {
		v, _ = handlePtr(v)
	}
	if v.Kind() != reflect.Struct {
		return errors.New("handleStruct: value isn't a struct")
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			if err := handleStruct(v.Field(i), f); err != nil {
				return err
			}
			continue
		}
		if !v.Field(i).IsValid() ||
			!v.Field(i).CanSet() {
			continue
		}
		f(v.Type().Field(i), v.Field(i))
	}
	return nil
}

func handleValue(
	fieldType reflect.Type,
	fieldTag reflect.StructTag,
	fieldValue reflect.Value,
	setValue string,
) {
	switch fieldType.Kind() {
	case reflect.Bool:
		if v, err := strconv.ParseBool(setValue); err == nil {
			fieldValue.SetBool(v)
		}
	case reflect.String:
		fieldValue.SetString(setValue)
	case reflect.Uint:
		if v, err := strconv.ParseUint(setValue, 10, 0); err == nil {
			fieldValue.SetUint(v)
		}
	case reflect.Uint8:
		if v, err := strconv.ParseUint(setValue, 10, 8); err == nil {
			fieldValue.SetUint(v)
		}
	case reflect.Uint16:
		if v, err := strconv.ParseUint(setValue, 10, 16); err == nil {
			fieldValue.SetUint(v)
		}
	case reflect.Uint32:
		if v, err := strconv.ParseUint(setValue, 10, 32); err == nil {
			fieldValue.SetUint(v)
		}
	case reflect.Uint64:
		if v, err := strconv.ParseUint(setValue, 10, 64); err == nil {
			fieldValue.SetUint(v)
		}
	case reflect.Int:
		if v, err := strconv.ParseInt(setValue, 10, 0); err == nil {
			fieldValue.SetInt(v)
		}
	case reflect.Int8:
		if v, err := strconv.ParseInt(setValue, 10, 8); err == nil {
			fieldValue.SetInt(v)
		}
	case reflect.Int16:
		if v, err := strconv.ParseInt(setValue, 10, 16); err == nil {
			fieldValue.SetInt(v)
		}
	case reflect.Int32:
		if v, err := strconv.ParseInt(setValue, 10, 32); err == nil {
			fieldValue.SetInt(v)
		}
	case reflect.Int64:
		if _, ok := fieldValue.Interface().(time.Duration); ok {
			d, err := time.ParseDuration(setValue)
			if err == nil {
				fieldValue.Set(reflect.ValueOf(d))
			}
			return
		}

		if v, err := strconv.ParseInt(setValue, 10, 64); err == nil {
			fieldValue.SetInt(v)
		}
	case reflect.Float32:
		if v, err := strconv.ParseFloat(setValue, 32); err == nil {
			fieldValue.SetFloat(v)
		}
	case reflect.Float64:
		if v, err := strconv.ParseFloat(setValue, 64); err == nil {
			fieldValue.SetFloat(v)
		}
	case reflect.Array:
		handleArray(fieldType, fieldValue, setValue)
	case reflect.Slice:
		handleSlice(fieldType, fieldValue, setValue)
	case reflect.Map:
		handleMap(fieldType, fieldValue, setValue)
	}
}

func handleField(structField reflect.StructField, structValue reflect.Value) {
	if structValue.Kind() == reflect.Ptr && structValue.IsNil() {
		return
	}
	envName := structField.Tag.Get("env")
	envVal := os.Getenv(envName)
	defaultVal := structField.Tag.Get("default")
	setValue := ""
	if defaultVal != "" {
		setValue = defaultVal
	}
	if envVal != "" {
		setValue = envVal
	}
	handleValue(structField.Type, structField.Tag, structValue, setValue)
}
