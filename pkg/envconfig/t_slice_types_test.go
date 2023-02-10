package envconfig

import (
	"os"
	"reflect"
	"testing"
)

func TestSliceTypesDefaults(t *testing.T) {
	os.Clearenv()
	config := new(sliceTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if !reflect.DeepEqual(config.SliceBoolVal, []bool{true, false}) {
		t.Errorf("config.SliceBoolVal should be '[true false]', got: %v", config.SliceBoolVal)
	}
	if !reflect.DeepEqual(config.SliceStringVal, []string{"one", "two"}) {
		t.Errorf("config.SliceStringValue should be [one two], got: %v", config.SliceStringVal)
	}

	if !reflect.DeepEqual(config.SliceUintVal, []uint{1, 2}) {
		t.Errorf("config.SliceUintValue should be '[1 2]', got: %v", config.SliceUintVal)
	}
	if !reflect.DeepEqual(config.SliceUint8Val, []uint8{1, 2}) {
		t.Errorf("config.SliceUint8Value should be '[1 2]', got: %v", config.SliceUint8Val)
	}
	if !reflect.DeepEqual(config.SliceUint16Val, []uint16{1, 2}) {
		t.Errorf("config.SliceUint16Value should be '[1 2]', got: %v", config.SliceUint16Val)
	}
	if !reflect.DeepEqual(config.SliceUint32Val, []uint32{1, 2}) {
		t.Errorf("config.SliceUint32Value should be '[1 2]', got: %v", config.SliceUint32Val)
	}
	if !reflect.DeepEqual(config.SliceUint64Val, []uint64{1, 2}) {
		t.Errorf("config.SliceUint64Value should be '[1 2]', got: %v", config.SliceUint64Val)
	}

	if !reflect.DeepEqual(config.SliceIntVal, []int{1, 2}) {
		t.Errorf("config.SliceIntValue should be '[1 2]', got: %v", config.SliceIntVal)
	}
	if !reflect.DeepEqual(config.SliceInt8Val, []int8{1, 2}) {
		t.Errorf("config.SliceInt8Value should be '[1 2]', got: %v", config.SliceInt8Val)
	}
	if !reflect.DeepEqual(config.SliceInt16Val, []int16{1, 2}) {
		t.Errorf("config.SliceInt16Value should be '[1 2]', got: %v", config.SliceInt16Val)
	}
	if !reflect.DeepEqual(config.SliceInt32Val, []int32{1, 2}) {
		t.Errorf("config.SliceInt32Value should be '[1 2]', got: %v", config.SliceInt32Val)
	}
	if !reflect.DeepEqual(config.SliceInt64Val, []int64{1, 2}) {
		t.Errorf("config.SliceUint64Value should be '[1 2]', got: %v", config.SliceInt64Val)
	}

	if !reflect.DeepEqual(config.SliceFloat32Val, []float32{1.0, 2.0}) {
		t.Errorf("config.SliceFloat32Value should be '[1.0 2.0]', got: '%.1f'", config.SliceFloat32Val)
	}
	if !reflect.DeepEqual(config.SliceFloat64Val, []float64{1.0, 2.0}) {
		t.Errorf("config.SliceFloat64Value should be '[1.0 2.0]', got: '%.1f'", config.SliceFloat64Val)
	}
}

func TestSliceTypes(t *testing.T) {
	os.Clearenv()
	envs := envValues{
		"ENV_SLICE_BOOL":    envValue{checkValue: [2]bool{true, true}},
		"ENV_SLICE_STRING":  envValue{checkValue: [2]string{"three", "four"}},
		"ENV_SLICE_UINT":    envValue{checkValue: [2]uint{2, 3}},
		"ENV_SLICE_UINT8":   envValue{checkValue: [2]uint8{2, 3}},
		"ENV_SLICE_UINT16":  envValue{checkValue: [2]uint16{2, 3}},
		"ENV_SLICE_UINT32":  envValue{checkValue: [2]uint32{2, 3}},
		"ENV_SLICE_UINT64":  envValue{checkValue: [2]uint64{2, 3}},
		"ENV_SLICE_INT":     envValue{checkValue: [2]int{2, 3}},
		"ENV_SLICE_INT8":    envValue{checkValue: [2]int8{2, 3}},
		"ENV_SLICE_INT16":   envValue{checkValue: [2]int16{2, 3}},
		"ENV_SLICE_INT32":   envValue{checkValue: [2]int32{2, 3}},
		"ENV_SLICE_INT64":   envValue{checkValue: [2]int64{2, 3}},
		"ENV_SLICE_FLOAT32": envValue{checkValue: [2]float32{2.0, 3.0}},
		"ENV_SLICE_FLOAT64": envValue{checkValue: [2]float64{2.0, 3.0}},
	}
	if err := envs.set(); err != nil {
		t.Fatalf("envs.set() fail, raw: %v", err)
	}
	config := new(sliceTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if !reflect.DeepEqual(config.SliceBoolVal, []bool{true, true}) {
		t.Errorf("config.SliceBoolVal should be '[true true]', got: %v", config.SliceBoolVal)
	}
	if !reflect.DeepEqual(config.SliceStringVal, []string{"three", "four"}) {
		t.Errorf("config.SliceStringValue should be [three four], got: %v", config.SliceStringVal)
	}

	if !reflect.DeepEqual(config.SliceUintVal, []uint{2, 3}) {
		t.Errorf("config.SliceUintValue should be '[2 3]', got: %v", config.SliceUintVal)
	}
	if !reflect.DeepEqual(config.SliceUint8Val, []uint8{2, 3}) {
		t.Errorf("config.SliceUint8Value should be '[2 3]', got: %v", config.SliceUint8Val)
	}
	if !reflect.DeepEqual(config.SliceUint16Val, []uint16{2, 3}) {
		t.Errorf("config.SliceUint16Value should be '[2 3]', got: %v", config.SliceUint16Val)
	}
	if !reflect.DeepEqual(config.SliceUint32Val, []uint32{2, 3}) {
		t.Errorf("config.SliceUint32Value should be '[2 3]', got: %v", config.SliceUint32Val)
	}
	if !reflect.DeepEqual(config.SliceUint64Val, []uint64{2, 3}) {
		t.Errorf("config.SliceUint64Value should be '[2 3]', got: %v", config.SliceUint64Val)
	}

	if !reflect.DeepEqual(config.SliceIntVal, []int{2, 3}) {
		t.Errorf("config.SliceIntValue should be '[2 3]', got: %v", config.SliceIntVal)
	}
	if !reflect.DeepEqual(config.SliceInt8Val, []int8{2, 3}) {
		t.Errorf("config.SliceInt8Value should be '[2 3]', got: %v", config.SliceInt8Val)
	}
	if !reflect.DeepEqual(config.SliceInt16Val, []int16{2, 3}) {
		t.Errorf("config.SliceInt16Value should be '[2 3]', got: %v", config.SliceInt16Val)
	}
	if !reflect.DeepEqual(config.SliceInt32Val, []int32{2, 3}) {
		t.Errorf("config.SliceInt32Value should be '[2 3]', got: %v", config.SliceInt32Val)
	}
	if !reflect.DeepEqual(config.SliceInt64Val, []int64{2, 3}) {
		t.Errorf("config.SliceUint64Value should be '[2 3]', got: %v", config.SliceInt64Val)
	}

	if !reflect.DeepEqual(config.SliceFloat32Val, []float32{2.0, 3.0}) {
		t.Errorf("config.SliceFloat32Value should be '[2.0 3.0]', got: '%.1f'", config.SliceFloat32Val)
	}
	if !reflect.DeepEqual(config.SliceFloat64Val, []float64{2.0, 3.0}) {
		t.Errorf("config.SliceFloat64Value should be '[2.0 3.0]', got: '%.1f'", config.SliceFloat64Val)
	}
}

// sliceTypes is an test struct that holds parameters of different slice types
type sliceTypes struct {
	SliceBoolVal    []bool    `env:"ENV_SLICE_BOOL" default:"true,false"`
	SliceStringVal  []string  `env:"ENV_SLICE_STRING" default:"one,two"`
	SliceUintVal    []uint    `env:"ENV_SLICE_UINT" default:"1,2"`
	SliceUint8Val   []uint8   `env:"ENV_SLICE_UINT8" default:"1,2"`
	SliceUint16Val  []uint16  `env:"ENV_SLICE_UINT16" default:"1,2"`
	SliceUint32Val  []uint32  `env:"ENV_SLICE_UINT32" default:"1,2"`
	SliceUint64Val  []uint64  `env:"ENV_SLICE_UINT64" default:"1,2"`
	SliceIntVal     []int     `env:"ENV_SLICE_INT" default:"1,2"`
	SliceInt8Val    []int8    `env:"ENV_SLICE_INT8" default:"1,2"`
	SliceInt16Val   []int16   `env:"ENV_SLICE_INT16" default:"1,2"`
	SliceInt32Val   []int32   `env:"ENV_SLICE_INT32" default:"1,2"`
	SliceInt64Val   []int64   `env:"ENV_SLICE_INT64" default:"1,2"`
	SliceFloat32Val []float32 `env:"ENV_SLICE_FLOAT32" default:"1.0,2.0"`
	SliceFloat64Val []float64 `env:"ENV_SLICE_FLOAT64" default:"1.0,2.0"`
}
