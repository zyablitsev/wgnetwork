package envconfig

import (
	"net"
	"reflect"
	"strconv"
	"strings"
)

func handleArray(fieldType reflect.Type, fieldValue reflect.Value, setValue string) {
	elems := strings.Split(setValue, ",")
	slice := fieldValue.Slice(0, fieldValue.Len())
	switch slice.Interface().(type) {
	case []bool:
		for idx, i := range elems {
			if v, err := strconv.ParseBool(i); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(v))
			}
		}
	case []string:
		for idx, i := range elems {
			fieldValue.Index(idx).Set(reflect.ValueOf(i))
		}
	case []uint:
		for idx, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 0); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(uint(v)))
			}
		}
	case []uint8:
		for idx, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 8); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(uint8(v)))
			}
		}
	case []uint16:
		for idx, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 16); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(uint16(v)))
			}
		}
	case []uint32:
		for idx, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 32); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(uint32(v)))
			}
		}
	case []uint64:
		for idx, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 64); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(v))
			}
		}
	case []int:
		for idx, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 0); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(int(v)))
			}
		}
	case []int8:
		for idx, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 8); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(int8(v)))
			}
		}
	case []int16:
		for idx, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 16); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(int16(v)))
			}
		}
	case []int32:
		for idx, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 32); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(int32(v)))
			}
		}
	case []int64:
		for idx, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 64); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(v))
			}
		}
	case []float32:
		for idx, i := range elems {
			if v, err := strconv.ParseFloat(i, 32); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(float32(v)))
			}
		}
	case []float64:
		for idx, i := range elems {
			if v, err := strconv.ParseFloat(i, 64); err == nil {
				fieldValue.Index(idx).Set(reflect.ValueOf(v))
			}
		}
	}
}

func handleSlice(fieldType reflect.Type, fieldValue reflect.Value, setValue string) {
	splitted := strings.Split(setValue, ",")
	elems := make([]string, 0, len(splitted))
	for _, v := range splitted {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		elems = append(elems, v)
	}

	if len(elems) == 0 {
		return
	}

	slicePtr := reflect.Value{}
	switch fieldValue.Interface().(type) {
	case []bool:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((bool)(false))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseBool(i); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(v)))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []string:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((string)(""))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(i)))
		}
		fieldValue.Set(slicePtr.Elem())
	case []uint:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((uint)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 0); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(uint(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []uint8:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((uint8)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 8); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(uint8(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []uint16:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((uint16)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 16); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(uint16(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []uint32:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((uint32)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 32); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(uint32(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []uint64:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((uint64)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseUint(i, 10, 64); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(uint64(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []int:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((int)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 0); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(int(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []int8:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((int8)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 8); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(int8(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []int16:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((int16)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 16); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(int16(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []int32:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((int32)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 32); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(int32(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []int64:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((int64)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseInt(i, 10, 64); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(int64(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []float32:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((float32)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseFloat(i, 32); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(float32(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case []float64:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((float64)(0))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, i := range elems {
			if v, err := strconv.ParseFloat(i, 64); err == nil {
				slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(float64(v))))
			}
		}
		fieldValue.Set(slicePtr.Elem())
	case net.IP:
		fieldValue.Set(reflect.ValueOf(net.ParseIP(setValue).To4()))
	case []net.IP:
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf((net.IP)(nil))), 0, len(elems))
		slicePtr = reflect.New(slice.Type())
		slicePtr.Elem().Set(slice)
		for _, v := range elems {
			slicePtr.Elem().Set(reflect.Append(slicePtr.Elem(), reflect.ValueOf(net.ParseIP(v).To4())))
		}
		fieldValue.Set(slicePtr.Elem())
	}
}

func handleMap(fieldType reflect.Type, fieldValue reflect.Value, setValue string) {
	elems := strings.Split(setValue, ",")
	mPtr := reflect.Value{}

	switch fieldValue.Interface().(type) {
	case map[bool]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[bool]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[bool]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseBool(kv[0]); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[string]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseBool(kv[1]); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(val))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(kv[1]))
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(uint(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(uint8(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(uint16(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(uint32(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(val))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(int(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(int8(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(int16(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(int32(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(val))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(float32(val)))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[string]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[string]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
				m.SetMapIndex(reflect.ValueOf(kv[0]), reflect.ValueOf(val))
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[uint]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[uint8]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint8]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint8]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[uint16]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint16]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint16]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[uint32]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint32]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint32]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[uint64]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[uint64]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[uint64]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseUint(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(uint64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[int]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 0); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[int8]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int8]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int8]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 8); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int8(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[int16]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int16]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int16]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 16); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int16(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[int32]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int32]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int32]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 32); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[int64]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[int64]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[int64]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseInt(kv[0], 10, 64); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(int64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[float32]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float32]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float32]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 32); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(float32(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())

	case map[float64]bool:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]bool)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseBool(kv[1]); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]string:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]string)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(kv[1]))
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]uint:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]uint)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(uint(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]uint8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]uint8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(uint8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]uint16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]uint16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(uint16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]uint32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]uint32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(uint32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]uint64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]uint64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseUint(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]int:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]int)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 0); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(int(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]int8:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]int8)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 8); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(int8(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]int16:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]int16)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 16); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(int16(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]int32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]int32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(int32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]int64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]int64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]float32:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]float32)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 32); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(float32(val)))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	case map[float64]float64:
		m := reflect.MakeMapWithSize(reflect.TypeOf((map[float64]float64)(nil)), len(elems))
		mPtr = reflect.New(m.Type())
		mPtr.Elem().Set(m)
		for _, i := range elems {
			kv := strings.Split(i, ":")
			if len(kv) != 2 {
				continue
			}
			if key, err := strconv.ParseFloat(kv[0], 64); err == nil {
				if val, err := strconv.ParseFloat(kv[1], 64); err == nil {
					m.SetMapIndex(reflect.ValueOf(float64(key)), reflect.ValueOf(val))
				}
			}
		}
		fieldValue.Set(mPtr.Elem())
	}
}
