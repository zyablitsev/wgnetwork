package envconfig

import (
	"os"
	"testing"
)

func TestArrTypesDefaults(t *testing.T) {
	os.Clearenv()
	config := new(arrTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if config.ArrBoolVal != [2]bool{true, false} {
		t.Errorf("config.ArrBoolVal should be '[true false]', got: %v", config.ArrBoolVal)
	}
	if config.ArrStringVal != [2]string{"one", "two"} {
		t.Errorf("config.ArrStringValue should be [one two], got: %v", config.ArrStringVal)
	}

	if config.ArrUintVal != [2]uint{1, 2} {
		t.Errorf("config.ArrUintValue should be '[1 2]', got: %v", config.ArrUintVal)
	}
	if config.ArrUint8Val != [2]uint8{1, 2} {
		t.Errorf("config.ArrUint8Value should be '[1 2]', got: %v", config.ArrUint8Val)
	}
	if config.ArrUint16Val != [2]uint16{1, 2} {
		t.Errorf("config.ArrUint16Value should be '[1 2]', got: %v", config.ArrUint16Val)
	}
	if config.ArrUint32Val != [2]uint32{1, 2} {
		t.Errorf("config.ArrUint32Value should be '[1 2]', got: %v", config.ArrUint32Val)
	}
	if config.ArrUint64Val != [2]uint64{1, 2} {
		t.Errorf("config.ArrUint64Value should be '[1 2]', got: %v", config.ArrUint64Val)
	}

	if config.ArrIntVal != [2]int{1, 2} {
		t.Errorf("config.ArrIntValue should be '[1 2]', got: %v", config.ArrIntVal)
	}
	if config.ArrInt8Val != [2]int8{1, 2} {
		t.Errorf("config.ArrInt8Value should be '[1 2]', got: %v", config.ArrInt8Val)
	}
	if config.ArrInt16Val != [2]int16{1, 2} {
		t.Errorf("config.ArrInt16Value should be '[1 2]', got: %v", config.ArrInt16Val)
	}
	if config.ArrInt32Val != [2]int32{1, 2} {
		t.Errorf("config.ArrInt32Value should be '[1 2]', got: %v", config.ArrInt32Val)
	}
	if config.ArrInt64Val != [2]int64{1, 2} {
		t.Errorf("config.ArrUint64Value should be '[1 2]', got: %v", config.ArrInt64Val)
	}

	if config.ArrFloat32Val != [2]float32{1.0, 2.0} {
		t.Errorf("config.ArrFloat32Value should be '[1.0 2.0]', got: '%.1f'", config.ArrFloat32Val)
	}
	if config.ArrFloat64Val != [2]float64{1.0, 2.0} {
		t.Errorf("config.ArrFloat64Value should be '[1.0 2.0]', got: '%.1f'", config.ArrFloat64Val)
	}
}

func TestArrTypes(t *testing.T) {
	os.Clearenv()
	envs := envValues{
		"ENV_ARR_BOOL":    envValue{checkValue: [2]bool{true, true}},
		"ENV_ARR_STRING":  envValue{checkValue: [2]string{"three", "four"}},
		"ENV_ARR_UINT":    envValue{checkValue: [2]uint{2, 3}},
		"ENV_ARR_UINT8":   envValue{checkValue: [2]uint8{2, 3}},
		"ENV_ARR_UINT16":  envValue{checkValue: [2]uint16{2, 3}},
		"ENV_ARR_UINT32":  envValue{checkValue: [2]uint32{2, 3}},
		"ENV_ARR_UINT64":  envValue{checkValue: [2]uint64{2, 3}},
		"ENV_ARR_INT":     envValue{checkValue: [2]int{2, 3}},
		"ENV_ARR_INT8":    envValue{checkValue: [2]int8{2, 3}},
		"ENV_ARR_INT16":   envValue{checkValue: [2]int16{2, 3}},
		"ENV_ARR_INT32":   envValue{checkValue: [2]int32{2, 3}},
		"ENV_ARR_INT64":   envValue{checkValue: [2]int64{2, 3}},
		"ENV_ARR_FLOAT32": envValue{checkValue: [2]float32{2.0, 3.0}},
		"ENV_ARR_FLOAT64": envValue{checkValue: [2]float64{2.0, 3.0}},
	}
	if err := envs.set(); err != nil {
		t.Fatalf("envs.set() fail, raw: %v", err)
	}
	config := new(arrTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if config.ArrBoolVal != [2]bool{true, true} {
		t.Errorf("config.ArrBoolVal should be '[true true]', got: %v", config.ArrBoolVal)
	}
	if config.ArrStringVal != [2]string{"three", "four"} {
		t.Errorf("config.ArrStringValue should be [three four], got: %v", config.ArrStringVal)
	}

	if config.ArrUintVal != [2]uint{2, 3} {
		t.Errorf("config.ArrUintValue should be '[2 3]', got: %v", config.ArrUintVal)
	}
	if config.ArrUint8Val != [2]uint8{2, 3} {
		t.Errorf("config.ArrUint8Value should be '[2 3]', got: %v", config.ArrUint8Val)
	}
	if config.ArrUint16Val != [2]uint16{2, 3} {
		t.Errorf("config.ArrUint16Value should be '[2 3]', got: %v", config.ArrUint16Val)
	}
	if config.ArrUint32Val != [2]uint32{2, 3} {
		t.Errorf("config.ArrUint32Value should be '[2 3]', got: %v", config.ArrUint32Val)
	}
	if config.ArrUint64Val != [2]uint64{2, 3} {
		t.Errorf("config.ArrUint64Value should be '[2 3]', got: %v", config.ArrUint64Val)
	}

	if config.ArrIntVal != [2]int{2, 3} {
		t.Errorf("config.ArrIntValue should be '[2 3]', got: %v", config.ArrIntVal)
	}
	if config.ArrInt8Val != [2]int8{2, 3} {
		t.Errorf("config.ArrInt8Value should be '[2 3]', got: %v", config.ArrInt8Val)
	}
	if config.ArrInt16Val != [2]int16{2, 3} {
		t.Errorf("config.ArrInt16Value should be '[2 3]', got: %v", config.ArrInt16Val)
	}
	if config.ArrInt32Val != [2]int32{2, 3} {
		t.Errorf("config.ArrInt32Value should be '[2 3]', got: %v", config.ArrInt32Val)
	}
	if config.ArrInt64Val != [2]int64{2, 3} {
		t.Errorf("config.ArrUint64Value should be '[2 3]', got: %v", config.ArrInt64Val)
	}

	if config.ArrFloat32Val != [2]float32{2.0, 3.0} {
		t.Errorf("config.ArrFloat32Value should be '[2.0 3.0]', got: '%.1f'", config.ArrFloat32Val)
	}
	if config.ArrFloat64Val != [2]float64{2.0, 3.0} {
		t.Errorf("config.ArrFloat64Value should be '[2.0 3.0]', got: '%.1f'", config.ArrFloat64Val)
	}
}

// arrTypes is an test struct that holds parameters of different arr types
type arrTypes struct {
	ArrBoolVal    [2]bool    `env:"ENV_ARR_BOOL" default:"true,false"`
	ArrStringVal  [2]string  `env:"ENV_ARR_STRING" default:"one,two"`
	ArrUintVal    [2]uint    `env:"ENV_ARR_UINT" default:"1,2"`
	ArrUint8Val   [2]uint8   `env:"ENV_ARR_UINT8" default:"1,2"`
	ArrUint16Val  [2]uint16  `env:"ENV_ARR_UINT16" default:"1,2"`
	ArrUint32Val  [2]uint32  `env:"ENV_ARR_UINT32" default:"1,2"`
	ArrUint64Val  [2]uint64  `env:"ENV_ARR_UINT64" default:"1,2"`
	ArrIntVal     [2]int     `env:"ENV_ARR_INT" default:"1,2"`
	ArrInt8Val    [2]int8    `env:"ENV_ARR_INT8" default:"1,2"`
	ArrInt16Val   [2]int16   `env:"ENV_ARR_INT16" default:"1,2"`
	ArrInt32Val   [2]int32   `env:"ENV_ARR_INT32" default:"1,2"`
	ArrInt64Val   [2]int64   `env:"ENV_ARR_INT64" default:"1,2"`
	ArrFloat32Val [2]float32 `env:"ENV_ARR_FLOAT32" default:"1.0,2.0"`
	ArrFloat64Val [2]float64 `env:"ENV_ARR_FLOAT64" default:"1.0,2.0"`
}
