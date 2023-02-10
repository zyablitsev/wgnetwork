package envconfig

import (
	"os"
	"reflect"
	"testing"
)

func TestMapTypesDefaults(t *testing.T) {
	os.Clearenv()
	config := new(mapTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if !reflect.DeepEqual(config.MapBoolBoolVal, map[bool]bool{true: true, false: false}) {
		t.Errorf("config.MapBoolBoolVal should be 'map[true:true false:false, got: '%v'", config.MapBoolBoolVal)
	}

	if !reflect.DeepEqual(config.MapBoolStringVal, map[bool]string{true: "s1", false: "s2"}) {
		t.Errorf("config.MapBoolStringVal should be 'map[true:s1 true:s2]', got: '%v'", config.MapBoolStringVal)
	}

	if !reflect.DeepEqual(config.MapBoolUintVal, map[bool]uint{true: 1, false: 2}) {
		t.Errorf("config.MapBoolUintVal should be 'map[true:1 false:2]', got: '%v'", config.MapBoolUintVal)
	}
	if !reflect.DeepEqual(config.MapBoolUint8Val, map[bool]uint8{true: 1, false: 2}) {
		t.Errorf("config.MapBoolUint8Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolUint8Val)
	}
	if !reflect.DeepEqual(config.MapBoolUint16Val, map[bool]uint16{true: 1, false: 2}) {
		t.Errorf("config.MapBoolUint16Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolUint16Val)
	}
	if !reflect.DeepEqual(config.MapBoolUint32Val, map[bool]uint32{true: 1, false: 2}) {
		t.Errorf("config.MapBoolUint32Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolUint32Val)
	}
	if !reflect.DeepEqual(config.MapBoolUint64Val, map[bool]uint64{true: 1, false: 2}) {
		t.Errorf("config.MapBoolUint64Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolUint64Val)
	}

	if !reflect.DeepEqual(config.MapBoolIntVal, map[bool]int{true: 1, false: 2}) {
		t.Errorf("config.MapBoolIntVal should be 'map[true:1 false:2]', got: '%v'", config.MapBoolIntVal)
	}
	if !reflect.DeepEqual(config.MapBoolInt8Val, map[bool]int8{true: 1, false: 2}) {
		t.Errorf("config.MapBoolInt8Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolInt8Val)
	}
	if !reflect.DeepEqual(config.MapBoolInt16Val, map[bool]int16{true: 1, false: 2}) {
		t.Errorf("config.MapBoolInt16Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolInt16Val)
	}
	if !reflect.DeepEqual(config.MapBoolInt32Val, map[bool]int32{true: 1, false: 2}) {
		t.Errorf("config.MapBoolInt32Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolInt32Val)
	}
	if !reflect.DeepEqual(config.MapBoolInt64Val, map[bool]int64{true: 1, false: 2}) {
		t.Errorf("config.MapBoolInt64Val should be 'map[true:1 false:2]', got: '%v'", config.MapBoolInt64Val)
	}

	if !reflect.DeepEqual(config.MapBoolFloat32Val, map[bool]float32{true: 1.0, false: 2.0}) {
		t.Errorf("config.MapBoolFloat32Val should be 'map[true:1.0 false:2.0]', got: '%v'", config.MapBoolFloat32Val)
	}
	if !reflect.DeepEqual(config.MapBoolFloat64Val, map[bool]float64{true: 1.0, false: 2.0}) {
		t.Errorf("config.MapBoolFloat64Val should be 'map[true:1.0 false:2.0]', got: '%v'", config.MapBoolFloat64Val)
	}

	if !reflect.DeepEqual(config.MapStringBoolVal, map[string]bool{"s1": true, "s2": false}) {
		t.Errorf("config.MapStringBoolVal should be 'map[s1:true s2:false, got: '%v'", config.MapStringBoolVal)
	}

	if !reflect.DeepEqual(config.MapStringStringVal, map[string]string{"s1": "s1", "s2": "s2"}) {
		t.Errorf("config.MapStringStringVal should be 'map[s1:s1 s2:s2]', got: '%v'", config.MapStringStringVal)
	}

	if !reflect.DeepEqual(config.MapStringUintVal, map[string]uint{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringUintVal should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringUintVal)
	}
	if !reflect.DeepEqual(config.MapStringUint8Val, map[string]uint8{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringUint8Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringUint8Val)
	}
	if !reflect.DeepEqual(config.MapStringUint16Val, map[string]uint16{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringUint16Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringUint16Val)
	}
	if !reflect.DeepEqual(config.MapStringUint32Val, map[string]uint32{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringUint32Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringUint32Val)
	}
	if !reflect.DeepEqual(config.MapStringUint64Val, map[string]uint64{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringUint64Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringUint64Val)
	}

	if !reflect.DeepEqual(config.MapStringIntVal, map[string]int{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringIntVal should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringIntVal)
	}
	if !reflect.DeepEqual(config.MapStringInt8Val, map[string]int8{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringInt8Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringInt8Val)
	}
	if !reflect.DeepEqual(config.MapStringInt16Val, map[string]int16{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringInt16Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringInt16Val)
	}
	if !reflect.DeepEqual(config.MapStringInt32Val, map[string]int32{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringInt32Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringInt32Val)
	}
	if !reflect.DeepEqual(config.MapStringInt64Val, map[string]int64{"s1": 1, "s2": 2}) {
		t.Errorf("config.MapStringInt64Val should be 'map[s1:1 s2:2]', got: '%v'", config.MapStringInt64Val)
	}

	if !reflect.DeepEqual(config.MapStringFloat32Val, map[string]float32{"s1": 1.0, "s2": 2.0}) {
		t.Errorf("config.MapStringFloat32Val should be 'map[s1:1.0 s2:2.0]', got: '%v'", config.MapStringFloat32Val)
	}
	if !reflect.DeepEqual(config.MapStringFloat64Val, map[string]float64{"s1": 1.0, "s2": 2.0}) {
		t.Errorf("config.MapStringFloat64Val should be 'map[s1:1.0 s2:2.0]', got: '%v'", config.MapStringFloat64Val)
	}

	if !reflect.DeepEqual(config.MapUintBoolVal, map[uint]bool{1: true, 2: false}) {
		t.Errorf("config.MapUintBoolVal should be 'map[1:true 2:false, got: '%v'", config.MapUintBoolVal)
	}

	if !reflect.DeepEqual(config.MapUintStringVal, map[uint]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapUintStringVal should be 'map[1:1 2:2]', got: '%v'", config.MapUintStringVal)
	}

	if !reflect.DeepEqual(config.MapUintUintVal, map[uint]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapUintUintVal should be 'map[1:1 2:2]', got: '%v'", config.MapUintUintVal)
	}
	if !reflect.DeepEqual(config.MapUintUint8Val, map[uint]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapUintUint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintUint8Val)
	}
	if !reflect.DeepEqual(config.MapUintUint16Val, map[uint]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapUintUint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintUint16Val)
	}
	if !reflect.DeepEqual(config.MapUintUint32Val, map[uint]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapUintUint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintUint32Val)
	}
	if !reflect.DeepEqual(config.MapUintUint64Val, map[uint]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapUintUint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintUint64Val)
	}

	if !reflect.DeepEqual(config.MapUintIntVal, map[uint]int{1: 1, 2: 2}) {
		t.Errorf("config.MapUintIntVal should be 'map[1:1 2:2]', got: '%v'", config.MapUintIntVal)
	}
	if !reflect.DeepEqual(config.MapUintInt8Val, map[uint]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapUintInt8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintInt8Val)
	}
	if !reflect.DeepEqual(config.MapUintInt16Val, map[uint]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapUintInt16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintInt16Val)
	}
	if !reflect.DeepEqual(config.MapUintInt32Val, map[uint]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapUintInt32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintInt32Val)
	}
	if !reflect.DeepEqual(config.MapUintInt64Val, map[uint]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapUintInt64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUintInt64Val)
	}

	if !reflect.DeepEqual(config.MapUintFloat32Val, map[uint]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUintFloat32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUintFloat32Val)
	}
	if !reflect.DeepEqual(config.MapUintFloat64Val, map[uint]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUintFloat64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUintFloat64Val)
	}

	if !reflect.DeepEqual(config.MapUint8BoolVal, map[uint8]bool{1: true, 2: false}) {
		t.Errorf("config.MapUint8BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapUint8BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint8StringVal, map[uint8]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapUint8StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint8StringVal)
	}

	if !reflect.DeepEqual(config.MapUint8UintVal, map[uint8]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint8UintVal)
	}
	if !reflect.DeepEqual(config.MapUint8Uint8Val, map[uint8]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint8Uint16Val, map[uint8]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint8Uint32Val, map[uint8]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint8Uint64Val, map[uint8]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint8IntVal, map[uint8]int{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint8IntVal)
	}
	if !reflect.DeepEqual(config.MapUint8Int8Val, map[uint8]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint8Int16Val, map[uint8]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint8Int32Val, map[uint8]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint8Int64Val, map[uint8]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint8Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint8Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint8Float32Val, map[uint8]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint8Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint8Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint8Float64Val, map[uint8]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint8Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint8Float64Val)
	}

	if !reflect.DeepEqual(config.MapUint16BoolVal, map[uint16]bool{1: true, 2: false}) {
		t.Errorf("config.MapUint16BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapUint16BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint16StringVal, map[uint16]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapUint16StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint16StringVal)
	}

	if !reflect.DeepEqual(config.MapUint16UintVal, map[uint16]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint16UintVal)
	}
	if !reflect.DeepEqual(config.MapUint16Uint8Val, map[uint16]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint16Uint16Val, map[uint16]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint16Uint32Val, map[uint16]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint16Uint64Val, map[uint16]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint16IntVal, map[uint16]int{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint16IntVal)
	}
	if !reflect.DeepEqual(config.MapUint16Int8Val, map[uint16]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint16Int16Val, map[uint16]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint16Int32Val, map[uint16]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint16Int64Val, map[uint16]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint16Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint16Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint16Float32Val, map[uint16]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint16Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint16Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint16Float64Val, map[uint16]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint16Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint16Float64Val)
	}

	if !reflect.DeepEqual(config.MapUint32BoolVal, map[uint32]bool{1: true, 2: false}) {
		t.Errorf("config.MapUint32BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapUint32BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint32StringVal, map[uint32]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapUint32StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint32StringVal)
	}

	if !reflect.DeepEqual(config.MapUint32UintVal, map[uint32]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint32UintVal)
	}
	if !reflect.DeepEqual(config.MapUint32Uint8Val, map[uint32]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint32Uint16Val, map[uint32]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint32Uint32Val, map[uint32]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint32Uint64Val, map[uint32]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint32IntVal, map[uint32]int{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint32IntVal)
	}
	if !reflect.DeepEqual(config.MapUint32Int8Val, map[uint32]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint32Int16Val, map[uint32]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint32Int32Val, map[uint32]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint32Int64Val, map[uint32]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint32Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint32Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint32Float32Val, map[uint32]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint32Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint32Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint32Float64Val, map[uint32]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint32Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint32Float64Val)
	}

	if !reflect.DeepEqual(config.MapUint64BoolVal, map[uint64]bool{1: true, 2: false}) {
		t.Errorf("config.MapUint64BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapUint64BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint64StringVal, map[uint64]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapUint64StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint64StringVal)
	}

	if !reflect.DeepEqual(config.MapUint64UintVal, map[uint64]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint64UintVal)
	}
	if !reflect.DeepEqual(config.MapUint64Uint8Val, map[uint64]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint64Uint16Val, map[uint64]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint64Uint32Val, map[uint64]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint64Uint64Val, map[uint64]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint64IntVal, map[uint64]int{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapUint64IntVal)
	}
	if !reflect.DeepEqual(config.MapUint64Int8Val, map[uint64]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint64Int16Val, map[uint64]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint64Int32Val, map[uint64]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint64Int64Val, map[uint64]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapUint64Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapUint64Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint64Float32Val, map[uint64]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint64Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint64Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint64Float64Val, map[uint64]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapUint64Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapUint64Float64Val)
	}

	if !reflect.DeepEqual(config.MapIntBoolVal, map[int]bool{1: true, 2: false}) {
		t.Errorf("config.MapIntBoolVal should be 'map[1:true 2:false, got: '%v'", config.MapIntBoolVal)
	}

	if !reflect.DeepEqual(config.MapIntStringVal, map[int]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapIntStringVal should be 'map[1:1 2:2]', got: '%v'", config.MapIntStringVal)
	}

	if !reflect.DeepEqual(config.MapIntUintVal, map[int]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapIntUintVal should be 'map[1:1 2:2]', got: '%v'", config.MapIntUintVal)
	}
	if !reflect.DeepEqual(config.MapIntUint8Val, map[int]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapIntUint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntUint8Val)
	}
	if !reflect.DeepEqual(config.MapIntUint16Val, map[int]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapIntUint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntUint16Val)
	}
	if !reflect.DeepEqual(config.MapIntUint32Val, map[int]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapIntUint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntUint32Val)
	}
	if !reflect.DeepEqual(config.MapIntUint64Val, map[int]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapIntUint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntUint64Val)
	}

	if !reflect.DeepEqual(config.MapIntIntVal, map[int]int{1: 1, 2: 2}) {
		t.Errorf("config.MapIntIntVal should be 'map[1:1 2:2]', got: '%v'", config.MapIntIntVal)
	}
	if !reflect.DeepEqual(config.MapIntInt8Val, map[int]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapIntInt8Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntInt8Val)
	}
	if !reflect.DeepEqual(config.MapIntInt16Val, map[int]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapIntInt16Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntInt16Val)
	}
	if !reflect.DeepEqual(config.MapIntInt32Val, map[int]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapIntInt32Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntInt32Val)
	}
	if !reflect.DeepEqual(config.MapIntInt64Val, map[int]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapIntInt64Val should be 'map[1:1 2:2]', got: '%v'", config.MapIntInt64Val)
	}

	if !reflect.DeepEqual(config.MapIntFloat32Val, map[int]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapIntFloat32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapIntFloat32Val)
	}
	if !reflect.DeepEqual(config.MapIntFloat64Val, map[int]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapIntFloat64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapIntFloat64Val)
	}

	if !reflect.DeepEqual(config.MapInt8BoolVal, map[int8]bool{1: true, 2: false}) {
		t.Errorf("config.MapInt8BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapInt8BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt8StringVal, map[int8]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapInt8StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt8StringVal)
	}

	if !reflect.DeepEqual(config.MapInt8UintVal, map[int8]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt8UintVal)
	}
	if !reflect.DeepEqual(config.MapInt8Uint8Val, map[int8]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt8Uint16Val, map[int8]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt8Uint32Val, map[int8]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt8Uint64Val, map[int8]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt8IntVal, map[int8]int{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt8IntVal)
	}
	if !reflect.DeepEqual(config.MapInt8Int8Val, map[int8]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt8Int16Val, map[int8]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt8Int32Val, map[int8]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt8Int64Val, map[int8]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt8Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt8Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt8Float32Val, map[int8]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt8Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt8Float32Val)
	}
	if !reflect.DeepEqual(config.MapInt8Float64Val, map[int8]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt8Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt8Float64Val)
	}

	if !reflect.DeepEqual(config.MapInt16BoolVal, map[int16]bool{1: true, 2: false}) {
		t.Errorf("config.MapInt16BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapInt16BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt16StringVal, map[int16]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapInt16StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt16StringVal)
	}

	if !reflect.DeepEqual(config.MapInt16UintVal, map[int16]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt16UintVal)
	}
	if !reflect.DeepEqual(config.MapInt16Uint8Val, map[int16]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt16Uint16Val, map[int16]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt16Uint32Val, map[int16]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt16Uint64Val, map[int16]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt16IntVal, map[int16]int{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt16IntVal)
	}
	if !reflect.DeepEqual(config.MapInt16Int8Val, map[int16]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt16Int16Val, map[int16]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt16Int32Val, map[int16]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt16Int64Val, map[int16]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt16Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt16Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt16Float32Val, map[int16]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt16Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt16Float32Val)
	}
	if !reflect.DeepEqual(config.MapInt16Float64Val, map[int16]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt16Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt16Float64Val)
	}

	if !reflect.DeepEqual(config.MapInt32BoolVal, map[int32]bool{1: true, 2: false}) {
		t.Errorf("config.MapInt32BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapInt32BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt32StringVal, map[int32]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapInt32StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt32StringVal)
	}

	if !reflect.DeepEqual(config.MapInt32UintVal, map[int32]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt32UintVal)
	}
	if !reflect.DeepEqual(config.MapInt32Uint8Val, map[int32]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt32Uint16Val, map[int32]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt32Uint32Val, map[int32]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt32Uint64Val, map[int32]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt32IntVal, map[int32]int{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt32IntVal)
	}
	if !reflect.DeepEqual(config.MapInt32Int8Val, map[int32]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt32Int16Val, map[int32]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt32Int32Val, map[int32]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt32Int64Val, map[int32]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt32Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt32Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt32Float32Val, map[int32]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt32Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt32Float32Val)
	}
	if !reflect.DeepEqual(config.MapInt32Float64Val, map[int32]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt32Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt32Float64Val)
	}

	if !reflect.DeepEqual(config.MapInt64BoolVal, map[int64]bool{1: true, 2: false}) {
		t.Errorf("config.MapInt64BoolVal should be 'map[1:true 2:false, got: '%v'", config.MapInt64BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt64StringVal, map[int64]string{1: "s1", 2: "s2"}) {
		t.Errorf("config.MapInt64StringVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt64StringVal)
	}

	if !reflect.DeepEqual(config.MapInt64UintVal, map[int64]uint{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64UintVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt64UintVal)
	}
	if !reflect.DeepEqual(config.MapInt64Uint8Val, map[int64]uint8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Uint8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt64Uint16Val, map[int64]uint16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Uint16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt64Uint32Val, map[int64]uint32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Uint32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt64Uint64Val, map[int64]uint64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Uint64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt64IntVal, map[int64]int{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64IntVal should be 'map[1:1 2:2]', got: '%v'", config.MapInt64IntVal)
	}
	if !reflect.DeepEqual(config.MapInt64Int8Val, map[int64]int8{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Int8Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt64Int16Val, map[int64]int16{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Int16Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt64Int32Val, map[int64]int32{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Int32Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt64Int64Val, map[int64]int64{1: 1, 2: 2}) {
		t.Errorf("config.MapInt64Int64Val should be 'map[1:1 2:2]', got: '%v'", config.MapInt64Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt64Float32Val, map[int64]float32{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt64Float32Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt64Float32Val)
	}

	if !reflect.DeepEqual(config.MapInt64Float64Val, map[int64]float64{1: 1.0, 2: 2.0}) {
		t.Errorf("config.MapInt64Float64Val should be 'map[1:1.0 2:2.0]', got: '%v'", config.MapInt64Float64Val)
	}

	if !reflect.DeepEqual(config.MapFloat32BoolVal, map[float32]bool{1.0: true, 2.0: false}) {
		t.Errorf("config.MapFloat32BoolVal should be 'map[1.0:true 2.0:false, got: '%v'", config.MapFloat32BoolVal)
	}

	if !reflect.DeepEqual(config.MapFloat32StringVal, map[float32]string{1.0: "s1", 2.0: "s2"}) {
		t.Errorf("config.MapFloat32StringVal should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32StringVal)
	}

	if !reflect.DeepEqual(config.MapFloat32UintVal, map[float32]uint{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32UintVal should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32UintVal)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint8Val, map[float32]uint8{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Uint8Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Uint8Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint16Val, map[float32]uint16{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Uint16Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Uint16Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint32Val, map[float32]uint32{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Uint32Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Uint32Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint64Val, map[float32]uint64{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Uint64Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Uint64Val)
	}

	if !reflect.DeepEqual(config.MapFloat32IntVal, map[float32]int{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32IntVal should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32IntVal)
	}
	if !reflect.DeepEqual(config.MapFloat32Int8Val, map[float32]int8{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Int8Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Int8Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Int16Val, map[float32]int16{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Int16Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Int16Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Int32Val, map[float32]int32{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Int32Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Int32Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Int64Val, map[float32]int64{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat32Int64Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat32Int64Val)
	}

	if !reflect.DeepEqual(config.MapFloat32Float32Val, map[float32]float32{1.0: 1.0, 2.0: 2.0}) {
		t.Errorf("config.MapFloat32Float32Val should be 'map[1.0:1.0 2.0:2.0]', got: '%v'", config.MapFloat32Float32Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Float64Val, map[float32]float64{1.0: 1.0, 2.0: 2.0}) {
		t.Errorf("config.MapFloat32Float64Val should be 'map[1.0:1.0 2.0:2.0]', got: '%v'", config.MapFloat32Float64Val)
	}

	if !reflect.DeepEqual(config.MapFloat64BoolVal, map[float64]bool{1.0: true, 2.0: false}) {
		t.Errorf("config.MapFloat64BoolVal should be 'map[1.0:true 2.0:false, got: '%v'", config.MapFloat64BoolVal)
	}

	if !reflect.DeepEqual(config.MapFloat64StringVal, map[float64]string{1.0: "s1", 2.0: "s2"}) {
		t.Errorf("config.MapFloat64StringVal should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64StringVal)
	}

	if !reflect.DeepEqual(config.MapFloat64UintVal, map[float64]uint{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64UintVal should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64UintVal)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint8Val, map[float64]uint8{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Uint8Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Uint8Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint16Val, map[float64]uint16{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Uint16Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Uint16Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint32Val, map[float64]uint32{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Uint32Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Uint32Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint64Val, map[float64]uint64{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Uint64Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Uint64Val)
	}

	if !reflect.DeepEqual(config.MapFloat64IntVal, map[float64]int{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64IntVal should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64IntVal)
	}
	if !reflect.DeepEqual(config.MapFloat64Int8Val, map[float64]int8{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Int8Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Int8Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Int16Val, map[float64]int16{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Int16Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Int16Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Int32Val, map[float64]int32{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Int32Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Int32Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Int64Val, map[float64]int64{1.0: 1, 2.0: 2}) {
		t.Errorf("config.MapFloat64Int64Val should be 'map[1.0:1 2.0:2]', got: '%v'", config.MapFloat64Int64Val)
	}

	if !reflect.DeepEqual(config.MapFloat64Float32Val, map[float64]float32{1.0: 1.0, 2.0: 2.0}) {
		t.Errorf("config.MapFloat64Float32Val should be 'map[1.0:1.0 2.0:2.0]', got: '%v'", config.MapFloat64Float32Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Float64Val, map[float64]float64{1.0: 1.0, 2.0: 2.0}) {
		t.Errorf("config.MapFloat64Float64Val should be 'map[1.0:1.0 2.0:2.0]', got: '%v'", config.MapFloat64Float64Val)
	}
}

func TestMapTypes(t *testing.T) {
	os.Clearenv()
	envs := envValues{
		"ENV_MAP_BOOL_BOOL":    envValue{checkValue: map[bool]bool{false: true, true: true}},
		"ENV_MAP_BOOL_STRING":  envValue{checkValue: map[bool]string{false: "three", true: "four"}},
		"ENV_MAP_BOOL_UINT":    envValue{checkValue: map[bool]uint{false: 2, true: 3}},
		"ENV_MAP_BOOL_UINT8":   envValue{checkValue: map[bool]uint8{false: 2, true: 3}},
		"ENV_MAP_BOOL_UINT16":  envValue{checkValue: map[bool]uint16{false: 2, true: 3}},
		"ENV_MAP_BOOL_UINT32":  envValue{checkValue: map[bool]uint32{false: 2, true: 3}},
		"ENV_MAP_BOOL_UINT64":  envValue{checkValue: map[bool]uint64{false: 2, true: 3}},
		"ENV_MAP_BOOL_INT":     envValue{checkValue: map[bool]int{false: 2, true: 3}},
		"ENV_MAP_BOOL_INT8":    envValue{checkValue: map[bool]int8{false: 2, true: 3}},
		"ENV_MAP_BOOL_INT16":   envValue{checkValue: map[bool]int16{false: 2, true: 3}},
		"ENV_MAP_BOOL_INT32":   envValue{checkValue: map[bool]int32{false: 2, true: 3}},
		"ENV_MAP_BOOL_INT64":   envValue{checkValue: map[bool]int64{false: 2, true: 3}},
		"ENV_MAP_BOOL_FLOAT32": envValue{checkValue: map[bool]float32{false: 2.0, true: 3.0}},
		"ENV_MAP_BOOL_FLOAT64": envValue{checkValue: map[bool]float64{false: 2.0, true: 3.0}},

		"ENV_MAP_STRING_BOOL":    envValue{checkValue: map[string]bool{"fooval1": true, "fooval2": true}},
		"ENV_MAP_STRING_STRING":  envValue{checkValue: map[string]string{"fooval1": "three", "fooval2": "four"}},
		"ENV_MAP_STRING_UINT":    envValue{checkValue: map[string]uint{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_UINT8":   envValue{checkValue: map[string]uint8{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_UINT16":  envValue{checkValue: map[string]uint16{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_UINT32":  envValue{checkValue: map[string]uint32{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_UINT64":  envValue{checkValue: map[string]uint64{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_INT":     envValue{checkValue: map[string]int{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_INT8":    envValue{checkValue: map[string]int8{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_INT16":   envValue{checkValue: map[string]int16{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_INT32":   envValue{checkValue: map[string]int32{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_INT64":   envValue{checkValue: map[string]int64{"fooval1": 2, "fooval2": 3}},
		"ENV_MAP_STRING_FLOAT32": envValue{checkValue: map[string]float32{"fooval1": 2.0, "fooval2": 3.0}},
		"ENV_MAP_STRING_FLOAT64": envValue{checkValue: map[string]float64{"fooval1": 2.0, "fooval2": 3.0}},

		"ENV_MAP_UINT_BOOL":    envValue{checkValue: map[uint]bool{2: true, 3: true}},
		"ENV_MAP_UINT_STRING":  envValue{checkValue: map[uint]string{2: "three", 3: "four"}},
		"ENV_MAP_UINT_UINT":    envValue{checkValue: map[uint]uint{2: 2, 3: 3}},
		"ENV_MAP_UINT_UINT8":   envValue{checkValue: map[uint]uint8{2: 2, 3: 3}},
		"ENV_MAP_UINT_UINT16":  envValue{checkValue: map[uint]uint16{2: 2, 3: 3}},
		"ENV_MAP_UINT_UINT32":  envValue{checkValue: map[uint]uint32{2: 2, 3: 3}},
		"ENV_MAP_UINT_UINT64":  envValue{checkValue: map[uint]uint64{2: 2, 3: 3}},
		"ENV_MAP_UINT_INT":     envValue{checkValue: map[uint]int{2: 2, 3: 3}},
		"ENV_MAP_UINT_INT8":    envValue{checkValue: map[uint]int8{2: 2, 3: 3}},
		"ENV_MAP_UINT_INT16":   envValue{checkValue: map[uint]int16{2: 2, 3: 3}},
		"ENV_MAP_UINT_INT32":   envValue{checkValue: map[uint]int32{2: 2, 3: 3}},
		"ENV_MAP_UINT_INT64":   envValue{checkValue: map[uint]int64{2: 2, 3: 3}},
		"ENV_MAP_UINT_FLOAT32": envValue{checkValue: map[uint]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_UINT_FLOAT64": envValue{checkValue: map[uint]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_UINT8_BOOL":    envValue{checkValue: map[uint8]bool{2: true, 3: true}},
		"ENV_MAP_UINT8_STRING":  envValue{checkValue: map[uint8]string{2: "three", 3: "four"}},
		"ENV_MAP_UINT8_UINT":    envValue{checkValue: map[uint8]uint{2: 2, 3: 3}},
		"ENV_MAP_UINT8_UINT8":   envValue{checkValue: map[uint8]uint8{2: 2, 3: 3}},
		"ENV_MAP_UINT8_UINT16":  envValue{checkValue: map[uint8]uint16{2: 2, 3: 3}},
		"ENV_MAP_UINT8_UINT32":  envValue{checkValue: map[uint8]uint32{2: 2, 3: 3}},
		"ENV_MAP_UINT8_UINT64":  envValue{checkValue: map[uint8]uint64{2: 2, 3: 3}},
		"ENV_MAP_UINT8_INT":     envValue{checkValue: map[uint8]int{2: 2, 3: 3}},
		"ENV_MAP_UINT8_INT8":    envValue{checkValue: map[uint8]int8{2: 2, 3: 3}},
		"ENV_MAP_UINT8_INT16":   envValue{checkValue: map[uint8]int16{2: 2, 3: 3}},
		"ENV_MAP_UINT8_INT32":   envValue{checkValue: map[uint8]int32{2: 2, 3: 3}},
		"ENV_MAP_UINT8_INT64":   envValue{checkValue: map[uint8]int64{2: 2, 3: 3}},
		"ENV_MAP_UINT8_FLOAT32": envValue{checkValue: map[uint8]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_UINT8_FLOAT64": envValue{checkValue: map[uint8]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_UINT16_BOOL":    envValue{checkValue: map[uint16]bool{2: true, 3: true}},
		"ENV_MAP_UINT16_STRING":  envValue{checkValue: map[uint16]string{2: "three", 3: "four"}},
		"ENV_MAP_UINT16_UINT":    envValue{checkValue: map[uint16]uint{2: 2, 3: 3}},
		"ENV_MAP_UINT16_UINT8":   envValue{checkValue: map[uint16]uint8{2: 2, 3: 3}},
		"ENV_MAP_UINT16_UINT16":  envValue{checkValue: map[uint16]uint16{2: 2, 3: 3}},
		"ENV_MAP_UINT16_UINT32":  envValue{checkValue: map[uint16]uint32{2: 2, 3: 3}},
		"ENV_MAP_UINT16_UINT64":  envValue{checkValue: map[uint16]uint64{2: 2, 3: 3}},
		"ENV_MAP_UINT16_INT":     envValue{checkValue: map[uint16]int{2: 2, 3: 3}},
		"ENV_MAP_UINT16_INT8":    envValue{checkValue: map[uint16]int8{2: 2, 3: 3}},
		"ENV_MAP_UINT16_INT16":   envValue{checkValue: map[uint16]int16{2: 2, 3: 3}},
		"ENV_MAP_UINT16_INT32":   envValue{checkValue: map[uint16]int32{2: 2, 3: 3}},
		"ENV_MAP_UINT16_INT64":   envValue{checkValue: map[uint16]int64{2: 2, 3: 3}},
		"ENV_MAP_UINT16_FLOAT32": envValue{checkValue: map[uint16]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_UINT16_FLOAT64": envValue{checkValue: map[uint16]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_UINT32_BOOL":    envValue{checkValue: map[uint32]bool{2: true, 3: true}},
		"ENV_MAP_UINT32_STRING":  envValue{checkValue: map[uint32]string{2: "three", 3: "four"}},
		"ENV_MAP_UINT32_UINT":    envValue{checkValue: map[uint32]uint{2: 2, 3: 3}},
		"ENV_MAP_UINT32_UINT8":   envValue{checkValue: map[uint32]uint8{2: 2, 3: 3}},
		"ENV_MAP_UINT32_UINT16":  envValue{checkValue: map[uint32]uint16{2: 2, 3: 3}},
		"ENV_MAP_UINT32_UINT32":  envValue{checkValue: map[uint32]uint32{2: 2, 3: 3}},
		"ENV_MAP_UINT32_UINT64":  envValue{checkValue: map[uint32]uint64{2: 2, 3: 3}},
		"ENV_MAP_UINT32_INT":     envValue{checkValue: map[uint32]int{2: 2, 3: 3}},
		"ENV_MAP_UINT32_INT8":    envValue{checkValue: map[uint32]int8{2: 2, 3: 3}},
		"ENV_MAP_UINT32_INT16":   envValue{checkValue: map[uint32]int16{2: 2, 3: 3}},
		"ENV_MAP_UINT32_INT32":   envValue{checkValue: map[uint32]int32{2: 2, 3: 3}},
		"ENV_MAP_UINT32_INT64":   envValue{checkValue: map[uint32]int64{2: 2, 3: 3}},
		"ENV_MAP_UINT32_FLOAT32": envValue{checkValue: map[uint32]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_UINT32_FLOAT64": envValue{checkValue: map[uint32]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_UINT64_BOOL":    envValue{checkValue: map[uint64]bool{2: true, 3: true}},
		"ENV_MAP_UINT64_STRING":  envValue{checkValue: map[uint64]string{2: "three", 3: "four"}},
		"ENV_MAP_UINT64_UINT":    envValue{checkValue: map[uint64]uint{2: 2, 3: 3}},
		"ENV_MAP_UINT64_UINT8":   envValue{checkValue: map[uint64]uint8{2: 2, 3: 3}},
		"ENV_MAP_UINT64_UINT16":  envValue{checkValue: map[uint64]uint16{2: 2, 3: 3}},
		"ENV_MAP_UINT64_UINT32":  envValue{checkValue: map[uint64]uint32{2: 2, 3: 3}},
		"ENV_MAP_UINT64_UINT64":  envValue{checkValue: map[uint64]uint64{2: 2, 3: 3}},
		"ENV_MAP_UINT64_INT":     envValue{checkValue: map[uint64]int{2: 2, 3: 3}},
		"ENV_MAP_UINT64_INT8":    envValue{checkValue: map[uint64]int8{2: 2, 3: 3}},
		"ENV_MAP_UINT64_INT16":   envValue{checkValue: map[uint64]int16{2: 2, 3: 3}},
		"ENV_MAP_UINT64_INT32":   envValue{checkValue: map[uint64]int32{2: 2, 3: 3}},
		"ENV_MAP_UINT64_INT64":   envValue{checkValue: map[uint64]int64{2: 2, 3: 3}},
		"ENV_MAP_UINT64_FLOAT32": envValue{checkValue: map[uint64]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_UINT64_FLOAT64": envValue{checkValue: map[uint64]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_INT_BOOL":    envValue{checkValue: map[int]bool{2: true, 3: true}},
		"ENV_MAP_INT_STRING":  envValue{checkValue: map[int]string{2: "three", 3: "four"}},
		"ENV_MAP_INT_UINT":    envValue{checkValue: map[int]uint{2: 2, 3: 3}},
		"ENV_MAP_INT_UINT8":   envValue{checkValue: map[int]uint8{2: 2, 3: 3}},
		"ENV_MAP_INT_UINT16":  envValue{checkValue: map[int]uint16{2: 2, 3: 3}},
		"ENV_MAP_INT_UINT32":  envValue{checkValue: map[int]uint32{2: 2, 3: 3}},
		"ENV_MAP_INT_UINT64":  envValue{checkValue: map[int]uint64{2: 2, 3: 3}},
		"ENV_MAP_INT_INT":     envValue{checkValue: map[int]int{2: 2, 3: 3}},
		"ENV_MAP_INT_INT8":    envValue{checkValue: map[int]int8{2: 2, 3: 3}},
		"ENV_MAP_INT_INT16":   envValue{checkValue: map[int]int16{2: 2, 3: 3}},
		"ENV_MAP_INT_INT32":   envValue{checkValue: map[int]int32{2: 2, 3: 3}},
		"ENV_MAP_INT_INT64":   envValue{checkValue: map[int]int64{2: 2, 3: 3}},
		"ENV_MAP_INT_FLOAT32": envValue{checkValue: map[int]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_INT_FLOAT64": envValue{checkValue: map[int]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_INT8_BOOL":    envValue{checkValue: map[int8]bool{2: true, 3: true}},
		"ENV_MAP_INT8_STRING":  envValue{checkValue: map[int8]string{2: "three", 3: "four"}},
		"ENV_MAP_INT8_UINT":    envValue{checkValue: map[int8]uint{2: 2, 3: 3}},
		"ENV_MAP_INT8_UINT8":   envValue{checkValue: map[int8]uint8{2: 2, 3: 3}},
		"ENV_MAP_INT8_UINT16":  envValue{checkValue: map[int8]uint16{2: 2, 3: 3}},
		"ENV_MAP_INT8_UINT32":  envValue{checkValue: map[int8]uint32{2: 2, 3: 3}},
		"ENV_MAP_INT8_UINT64":  envValue{checkValue: map[int8]uint64{2: 2, 3: 3}},
		"ENV_MAP_INT8_INT":     envValue{checkValue: map[int8]int{2: 2, 3: 3}},
		"ENV_MAP_INT8_INT8":    envValue{checkValue: map[int8]int8{2: 2, 3: 3}},
		"ENV_MAP_INT8_INT16":   envValue{checkValue: map[int8]int16{2: 2, 3: 3}},
		"ENV_MAP_INT8_INT32":   envValue{checkValue: map[int8]int32{2: 2, 3: 3}},
		"ENV_MAP_INT8_INT64":   envValue{checkValue: map[int8]int64{2: 2, 3: 3}},
		"ENV_MAP_INT8_FLOAT32": envValue{checkValue: map[int8]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_INT8_FLOAT64": envValue{checkValue: map[int8]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_INT16_BOOL":    envValue{checkValue: map[int16]bool{2: true, 3: true}},
		"ENV_MAP_INT16_STRING":  envValue{checkValue: map[int16]string{2: "three", 3: "four"}},
		"ENV_MAP_INT16_UINT":    envValue{checkValue: map[int16]uint{2: 2, 3: 3}},
		"ENV_MAP_INT16_UINT8":   envValue{checkValue: map[int16]uint8{2: 2, 3: 3}},
		"ENV_MAP_INT16_UINT16":  envValue{checkValue: map[int16]uint16{2: 2, 3: 3}},
		"ENV_MAP_INT16_UINT32":  envValue{checkValue: map[int16]uint32{2: 2, 3: 3}},
		"ENV_MAP_INT16_UINT64":  envValue{checkValue: map[int16]uint64{2: 2, 3: 3}},
		"ENV_MAP_INT16_INT":     envValue{checkValue: map[int16]int{2: 2, 3: 3}},
		"ENV_MAP_INT16_INT8":    envValue{checkValue: map[int16]int8{2: 2, 3: 3}},
		"ENV_MAP_INT16_INT16":   envValue{checkValue: map[int16]int16{2: 2, 3: 3}},
		"ENV_MAP_INT16_INT32":   envValue{checkValue: map[int16]int32{2: 2, 3: 3}},
		"ENV_MAP_INT16_INT64":   envValue{checkValue: map[int16]int64{2: 2, 3: 3}},
		"ENV_MAP_INT16_FLOAT32": envValue{checkValue: map[int16]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_INT16_FLOAT64": envValue{checkValue: map[int16]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_INT32_BOOL":    envValue{checkValue: map[int32]bool{2: true, 3: true}},
		"ENV_MAP_INT32_STRING":  envValue{checkValue: map[int32]string{2: "three", 3: "four"}},
		"ENV_MAP_INT32_UINT":    envValue{checkValue: map[int32]uint{2: 2, 3: 3}},
		"ENV_MAP_INT32_UINT8":   envValue{checkValue: map[int32]uint8{2: 2, 3: 3}},
		"ENV_MAP_INT32_UINT16":  envValue{checkValue: map[int32]uint16{2: 2, 3: 3}},
		"ENV_MAP_INT32_UINT32":  envValue{checkValue: map[int32]uint32{2: 2, 3: 3}},
		"ENV_MAP_INT32_UINT64":  envValue{checkValue: map[int32]uint64{2: 2, 3: 3}},
		"ENV_MAP_INT32_INT":     envValue{checkValue: map[int32]int{2: 2, 3: 3}},
		"ENV_MAP_INT32_INT8":    envValue{checkValue: map[int32]int8{2: 2, 3: 3}},
		"ENV_MAP_INT32_INT16":   envValue{checkValue: map[int32]int16{2: 2, 3: 3}},
		"ENV_MAP_INT32_INT32":   envValue{checkValue: map[int32]int32{2: 2, 3: 3}},
		"ENV_MAP_INT32_INT64":   envValue{checkValue: map[int32]int64{2: 2, 3: 3}},
		"ENV_MAP_INT32_FLOAT32": envValue{checkValue: map[int32]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_INT32_FLOAT64": envValue{checkValue: map[int32]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_INT64_BOOL":    envValue{checkValue: map[int64]bool{2: true, 3: true}},
		"ENV_MAP_INT64_STRING":  envValue{checkValue: map[int64]string{2: "three", 3: "four"}},
		"ENV_MAP_INT64_UINT":    envValue{checkValue: map[int64]uint{2: 2, 3: 3}},
		"ENV_MAP_INT64_UINT8":   envValue{checkValue: map[int64]uint8{2: 2, 3: 3}},
		"ENV_MAP_INT64_UINT16":  envValue{checkValue: map[int64]uint16{2: 2, 3: 3}},
		"ENV_MAP_INT64_UINT32":  envValue{checkValue: map[int64]uint32{2: 2, 3: 3}},
		"ENV_MAP_INT64_UINT64":  envValue{checkValue: map[int64]uint64{2: 2, 3: 3}},
		"ENV_MAP_INT64_INT":     envValue{checkValue: map[int64]int{2: 2, 3: 3}},
		"ENV_MAP_INT64_INT8":    envValue{checkValue: map[int64]int8{2: 2, 3: 3}},
		"ENV_MAP_INT64_INT16":   envValue{checkValue: map[int64]int16{2: 2, 3: 3}},
		"ENV_MAP_INT64_INT32":   envValue{checkValue: map[int64]int32{2: 2, 3: 3}},
		"ENV_MAP_INT64_INT64":   envValue{checkValue: map[int64]int64{2: 2, 3: 3}},
		"ENV_MAP_INT64_FLOAT32": envValue{checkValue: map[int64]float32{2: 2.0, 3: 3.0}},
		"ENV_MAP_INT64_FLOAT64": envValue{checkValue: map[int64]float64{2: 2.0, 3: 3.0}},

		"ENV_MAP_FLOAT32_BOOL":    envValue{checkValue: map[float32]bool{2: true, 3: true}},
		"ENV_MAP_FLOAT32_STRING":  envValue{checkValue: map[float32]string{2: "three", 3: "four"}},
		"ENV_MAP_FLOAT32_UINT":    envValue{checkValue: map[float32]uint{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_UINT8":   envValue{checkValue: map[float32]uint8{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_UINT16":  envValue{checkValue: map[float32]uint16{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_UINT32":  envValue{checkValue: map[float32]uint32{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_UINT64":  envValue{checkValue: map[float32]uint64{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_INT":     envValue{checkValue: map[float32]int{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_INT8":    envValue{checkValue: map[float32]int8{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_INT16":   envValue{checkValue: map[float32]int16{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_INT32":   envValue{checkValue: map[float32]int32{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_INT64":   envValue{checkValue: map[float32]int64{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT32_FLOAT32": envValue{checkValue: map[float32]float32{2.0: 2.0, 3.0: 3.0}},
		"ENV_MAP_FLOAT32_FLOAT64": envValue{checkValue: map[float32]float64{2.0: 2.0, 3.0: 3.0}},

		"ENV_MAP_FLOAT64_BOOL":    envValue{checkValue: map[float64]bool{2.0: true, 3.0: true}},
		"ENV_MAP_FLOAT64_STRING":  envValue{checkValue: map[float64]string{2.0: "three", 3.0: "four"}},
		"ENV_MAP_FLOAT64_UINT":    envValue{checkValue: map[float64]uint{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_UINT8":   envValue{checkValue: map[float64]uint8{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_UINT16":  envValue{checkValue: map[float64]uint16{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_UINT32":  envValue{checkValue: map[float64]uint32{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_UINT64":  envValue{checkValue: map[float64]uint64{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_INT":     envValue{checkValue: map[float64]int{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_INT8":    envValue{checkValue: map[float64]int8{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_INT16":   envValue{checkValue: map[float64]int16{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_INT32":   envValue{checkValue: map[float64]int32{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_INT64":   envValue{checkValue: map[float64]int64{2.0: 2, 3.0: 3}},
		"ENV_MAP_FLOAT64_FLOAT32": envValue{checkValue: map[float64]float32{2.0: 2.0, 3.0: 3.0}},
		"ENV_MAP_FLOAT64_FLOAT64": envValue{checkValue: map[float64]float64{2.0: 2.0, 3.0: 3.0}},
	}
	if err := envs.set(); err != nil {
		t.Fatalf("envs.set() fail, raw: %v", err)
	}
	config := new(mapTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if !reflect.DeepEqual(config.MapBoolBoolVal, map[bool]bool{false: true, true: true}) {
		t.Errorf("config.MapBoolBoolVal should be 'map[false:true true:true, got: '%v'", config.MapBoolBoolVal)
	}

	if !reflect.DeepEqual(config.MapBoolStringVal, map[bool]string{false: "three", true: "four"}) {
		t.Errorf("config.MapBoolStringVal should be 'map[false:three true:four]', got: '%v'", config.MapBoolStringVal)
	}

	if !reflect.DeepEqual(config.MapBoolUintVal, map[bool]uint{false: 2, true: 3}) {
		t.Errorf("config.MapBoolUintVal should be 'map[false:2 true:3]', got: '%v'", config.MapBoolUintVal)
	}
	if !reflect.DeepEqual(config.MapBoolUint8Val, map[bool]uint8{false: 2, true: 3}) {
		t.Errorf("config.MapBoolUint8Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolUint8Val)
	}
	if !reflect.DeepEqual(config.MapBoolUint16Val, map[bool]uint16{false: 2, true: 3}) {
		t.Errorf("config.MapBoolUint16Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolUint16Val)
	}
	if !reflect.DeepEqual(config.MapBoolUint32Val, map[bool]uint32{false: 2, true: 3}) {
		t.Errorf("config.MapBoolUint32Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolUint32Val)
	}
	if !reflect.DeepEqual(config.MapBoolUint64Val, map[bool]uint64{false: 2, true: 3}) {
		t.Errorf("config.MapBoolUint64Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolUint64Val)
	}

	if !reflect.DeepEqual(config.MapBoolIntVal, map[bool]int{false: 2, true: 3}) {
		t.Errorf("config.MapBoolIntVal should be 'map[false:2 true:3]', got: '%v'", config.MapBoolIntVal)
	}
	if !reflect.DeepEqual(config.MapBoolInt8Val, map[bool]int8{false: 2, true: 3}) {
		t.Errorf("config.MapBoolInt8Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolInt8Val)
	}
	if !reflect.DeepEqual(config.MapBoolInt16Val, map[bool]int16{false: 2, true: 3}) {
		t.Errorf("config.MapBoolInt16Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolInt16Val)
	}
	if !reflect.DeepEqual(config.MapBoolInt32Val, map[bool]int32{false: 2, true: 3}) {
		t.Errorf("config.MapBoolInt32Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolInt32Val)
	}
	if !reflect.DeepEqual(config.MapBoolInt64Val, map[bool]int64{false: 2, true: 3}) {
		t.Errorf("config.MapBoolInt64Val should be 'map[false:2 true:3]', got: '%v'", config.MapBoolInt64Val)
	}

	if !reflect.DeepEqual(config.MapBoolFloat32Val, map[bool]float32{false: 2.0, true: 3.0}) {
		t.Errorf("config.MapBoolFloat32Val should be 'map[false:2.0 true:3.0]', got: '%v'", config.MapBoolFloat32Val)
	}
	if !reflect.DeepEqual(config.MapBoolFloat64Val, map[bool]float64{false: 2.0, true: 3.0}) {
		t.Errorf("config.MapBoolFloat64Val should be 'map[false:2.0 true:3.0]', got: '%v'", config.MapBoolFloat64Val)
	}

	if !reflect.DeepEqual(config.MapUintBoolVal, map[uint]bool{2: true, 3: true}) {
		t.Errorf("config.MapUintBoolVal should be 'map[2:true 3:true, got: '%v'", config.MapUintBoolVal)
	}

	if !reflect.DeepEqual(config.MapUintStringVal, map[uint]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapUintStringVal should be 'map[2:three 3:four]', got: '%v'", config.MapUintStringVal)
	}

	if !reflect.DeepEqual(config.MapUintUintVal, map[uint]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapUintUintVal should be 'map[2:2 3:3]', got: '%v'", config.MapUintUintVal)
	}
	if !reflect.DeepEqual(config.MapUintUint8Val, map[uint]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapUintUint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintUint8Val)
	}
	if !reflect.DeepEqual(config.MapUintUint16Val, map[uint]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapUintUint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintUint16Val)
	}
	if !reflect.DeepEqual(config.MapUintUint32Val, map[uint]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapUintUint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintUint32Val)
	}
	if !reflect.DeepEqual(config.MapUintUint64Val, map[uint]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapUintUint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintUint64Val)
	}

	if !reflect.DeepEqual(config.MapUintIntVal, map[uint]int{2: 2, 3: 3}) {
		t.Errorf("config.MapUintIntVal should be 'map[2:2 3:3]', got: '%v'", config.MapUintIntVal)
	}
	if !reflect.DeepEqual(config.MapUintInt8Val, map[uint]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapUintInt8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintInt8Val)
	}
	if !reflect.DeepEqual(config.MapUintInt16Val, map[uint]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapUintInt16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintInt16Val)
	}
	if !reflect.DeepEqual(config.MapUintInt32Val, map[uint]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapUintInt32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintInt32Val)
	}
	if !reflect.DeepEqual(config.MapUintInt64Val, map[uint]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapUintInt64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUintInt64Val)
	}

	if !reflect.DeepEqual(config.MapUintFloat32Val, map[uint]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUintFloat32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUintFloat32Val)
	}
	if !reflect.DeepEqual(config.MapUintFloat64Val, map[uint]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUintFloat64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUintFloat64Val)
	}

	if !reflect.DeepEqual(config.MapUint8BoolVal, map[uint8]bool{2: true, 3: true}) {
		t.Errorf("config.MapUint8BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapUint8BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint8StringVal, map[uint8]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapUint8StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapUint8StringVal)
	}

	if !reflect.DeepEqual(config.MapUint8UintVal, map[uint8]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint8UintVal)
	}
	if !reflect.DeepEqual(config.MapUint8Uint8Val, map[uint8]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint8Uint16Val, map[uint8]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint8Uint32Val, map[uint8]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint8Uint64Val, map[uint8]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint8IntVal, map[uint8]int{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint8IntVal)
	}
	if !reflect.DeepEqual(config.MapUint8Int8Val, map[uint8]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint8Int16Val, map[uint8]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint8Int32Val, map[uint8]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint8Int64Val, map[uint8]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint8Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint8Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint8Float32Val, map[uint8]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint8Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint8Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint8Float64Val, map[uint8]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint8Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint8Float64Val)
	}

	if !reflect.DeepEqual(config.MapUint16BoolVal, map[uint16]bool{2: true, 3: true}) {
		t.Errorf("config.MapUint16BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapUint16BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint16StringVal, map[uint16]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapUint16StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapUint16StringVal)
	}

	if !reflect.DeepEqual(config.MapUint16UintVal, map[uint16]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint16UintVal)
	}
	if !reflect.DeepEqual(config.MapUint16Uint8Val, map[uint16]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint16Uint16Val, map[uint16]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint16Uint32Val, map[uint16]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint16Uint64Val, map[uint16]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint16IntVal, map[uint16]int{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint16IntVal)
	}
	if !reflect.DeepEqual(config.MapUint16Int8Val, map[uint16]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint16Int16Val, map[uint16]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint16Int32Val, map[uint16]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint16Int64Val, map[uint16]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint16Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint16Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint16Float32Val, map[uint16]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint16Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint16Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint16Float64Val, map[uint16]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint16Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint16Float64Val)
	}

	if !reflect.DeepEqual(config.MapUint32BoolVal, map[uint32]bool{2: true, 3: true}) {
		t.Errorf("config.MapUint32BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapUint32BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint32StringVal, map[uint32]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapUint32StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapUint32StringVal)
	}

	if !reflect.DeepEqual(config.MapUint32UintVal, map[uint32]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint32UintVal)
	}
	if !reflect.DeepEqual(config.MapUint32Uint8Val, map[uint32]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint32Uint16Val, map[uint32]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint32Uint32Val, map[uint32]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint32Uint64Val, map[uint32]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint32IntVal, map[uint32]int{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint32IntVal)
	}
	if !reflect.DeepEqual(config.MapUint32Int8Val, map[uint32]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint32Int16Val, map[uint32]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint32Int32Val, map[uint32]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint32Int64Val, map[uint32]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint32Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint32Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint32Float32Val, map[uint32]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint32Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint32Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint32Float64Val, map[uint32]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint32Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint32Float64Val)
	}

	if !reflect.DeepEqual(config.MapUint64BoolVal, map[uint64]bool{2: true, 3: true}) {
		t.Errorf("config.MapUint64BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapUint64BoolVal)
	}

	if !reflect.DeepEqual(config.MapUint64StringVal, map[uint64]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapUint64StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapUint64StringVal)
	}

	if !reflect.DeepEqual(config.MapUint64UintVal, map[uint64]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint64UintVal)
	}
	if !reflect.DeepEqual(config.MapUint64Uint8Val, map[uint64]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Uint8Val)
	}
	if !reflect.DeepEqual(config.MapUint64Uint16Val, map[uint64]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Uint16Val)
	}
	if !reflect.DeepEqual(config.MapUint64Uint32Val, map[uint64]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Uint32Val)
	}
	if !reflect.DeepEqual(config.MapUint64Uint64Val, map[uint64]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Uint64Val)
	}

	if !reflect.DeepEqual(config.MapUint64IntVal, map[uint64]int{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapUint64IntVal)
	}
	if !reflect.DeepEqual(config.MapUint64Int8Val, map[uint64]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Int8Val)
	}
	if !reflect.DeepEqual(config.MapUint64Int16Val, map[uint64]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Int16Val)
	}
	if !reflect.DeepEqual(config.MapUint64Int32Val, map[uint64]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Int32Val)
	}
	if !reflect.DeepEqual(config.MapUint64Int64Val, map[uint64]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapUint64Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapUint64Int64Val)
	}

	if !reflect.DeepEqual(config.MapUint64Float32Val, map[uint64]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint64Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint64Float32Val)
	}
	if !reflect.DeepEqual(config.MapUint64Float64Val, map[uint64]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapUint64Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapUint64Float64Val)
	}

	if !reflect.DeepEqual(config.MapIntBoolVal, map[int]bool{2: true, 3: true}) {
		t.Errorf("config.MapIntBoolVal should be 'map[2:true 3:true], got: '%v'", config.MapIntBoolVal)
	}

	if !reflect.DeepEqual(config.MapIntStringVal, map[int]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapIntStringVal should be 'map[2:three 3:four]', got: '%v'", config.MapIntStringVal)
	}

	if !reflect.DeepEqual(config.MapIntUintVal, map[int]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapIntUintVal should be 'map[2:2 3:3]', got: '%v'", config.MapIntUintVal)
	}
	if !reflect.DeepEqual(config.MapIntUint8Val, map[int]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapIntUint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntUint8Val)
	}
	if !reflect.DeepEqual(config.MapIntUint16Val, map[int]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapIntUint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntUint16Val)
	}
	if !reflect.DeepEqual(config.MapIntUint32Val, map[int]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapIntUint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntUint32Val)
	}
	if !reflect.DeepEqual(config.MapIntUint64Val, map[int]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapIntUint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntUint64Val)
	}

	if !reflect.DeepEqual(config.MapIntIntVal, map[int]int{2: 2, 3: 3}) {
		t.Errorf("config.MapIntIntVal should be 'map[2:2 3:3]', got: '%v'", config.MapIntIntVal)
	}
	if !reflect.DeepEqual(config.MapIntInt8Val, map[int]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapIntInt8Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntInt8Val)
	}
	if !reflect.DeepEqual(config.MapIntInt16Val, map[int]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapIntInt16Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntInt16Val)
	}
	if !reflect.DeepEqual(config.MapIntInt32Val, map[int]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapIntInt32Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntInt32Val)
	}
	if !reflect.DeepEqual(config.MapIntInt64Val, map[int]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapIntInt64Val should be 'map[2:2 3:3]', got: '%v'", config.MapIntInt64Val)
	}

	if !reflect.DeepEqual(config.MapIntFloat32Val, map[int]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapIntFloat32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapIntFloat32Val)
	}
	if !reflect.DeepEqual(config.MapIntFloat64Val, map[int]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapIntFloat64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapIntFloat64Val)
	}

	if !reflect.DeepEqual(config.MapInt8BoolVal, map[int8]bool{2: true, 3: true}) {
		t.Errorf("config.MapInt8BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapInt8BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt8StringVal, map[int8]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapInt8StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapInt8StringVal)
	}

	if !reflect.DeepEqual(config.MapInt8UintVal, map[int8]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt8UintVal)
	}
	if !reflect.DeepEqual(config.MapInt8Uint8Val, map[int8]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt8Uint16Val, map[int8]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt8Uint32Val, map[int8]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt8Uint64Val, map[int8]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt8IntVal, map[int8]int{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt8IntVal)
	}
	if !reflect.DeepEqual(config.MapInt8Int8Val, map[int8]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt8Int16Val, map[int8]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt8Int32Val, map[int8]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt8Int64Val, map[int8]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt8Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt8Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt8Float32Val, map[int8]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt8Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt8Float32Val)
	}
	if !reflect.DeepEqual(config.MapInt8Float64Val, map[int8]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt8Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt8Float64Val)
	}

	if !reflect.DeepEqual(config.MapInt16BoolVal, map[int16]bool{2: true, 3: true}) {
		t.Errorf("config.MapInt16BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapInt16BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt16StringVal, map[int16]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapInt16StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapInt16StringVal)
	}

	if !reflect.DeepEqual(config.MapInt16UintVal, map[int16]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt16UintVal)
	}
	if !reflect.DeepEqual(config.MapInt16Uint8Val, map[int16]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt16Uint16Val, map[int16]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt16Uint32Val, map[int16]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt16Uint64Val, map[int16]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt16IntVal, map[int16]int{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt16IntVal)
	}
	if !reflect.DeepEqual(config.MapInt16Int8Val, map[int16]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt16Int16Val, map[int16]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt16Int32Val, map[int16]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt16Int64Val, map[int16]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt16Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt16Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt16Float32Val, map[int16]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt16Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt16Float32Val)
	}
	if !reflect.DeepEqual(config.MapInt16Float64Val, map[int16]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt16Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt16Float64Val)
	}

	if !reflect.DeepEqual(config.MapInt32BoolVal, map[int32]bool{2: true, 3: true}) {
		t.Errorf("config.MapInt32BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapInt32BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt32StringVal, map[int32]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapInt32StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapInt32StringVal)
	}

	if !reflect.DeepEqual(config.MapInt32UintVal, map[int32]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt32UintVal)
	}
	if !reflect.DeepEqual(config.MapInt32Uint8Val, map[int32]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt32Uint16Val, map[int32]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt32Uint32Val, map[int32]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt32Uint64Val, map[int32]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt32IntVal, map[int32]int{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt32IntVal)
	}
	if !reflect.DeepEqual(config.MapInt32Int8Val, map[int32]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt32Int16Val, map[int32]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt32Int32Val, map[int32]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt32Int64Val, map[int32]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt32Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt32Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt32Float32Val, map[int32]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt32Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt32Float32Val)
	}
	if !reflect.DeepEqual(config.MapInt32Float64Val, map[int32]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt32Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt32Float64Val)
	}

	if !reflect.DeepEqual(config.MapInt64BoolVal, map[int64]bool{2: true, 3: true}) {
		t.Errorf("config.MapInt64BoolVal should be 'map[2:true 3:true], got: '%v'", config.MapInt64BoolVal)
	}

	if !reflect.DeepEqual(config.MapInt64StringVal, map[int64]string{2: "three", 3: "four"}) {
		t.Errorf("config.MapInt64StringVal should be 'map[2:three 3:four]', got: '%v'", config.MapInt64StringVal)
	}

	if !reflect.DeepEqual(config.MapInt64UintVal, map[int64]uint{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64UintVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt64UintVal)
	}
	if !reflect.DeepEqual(config.MapInt64Uint8Val, map[int64]uint8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Uint8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Uint8Val)
	}
	if !reflect.DeepEqual(config.MapInt64Uint16Val, map[int64]uint16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Uint16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Uint16Val)
	}
	if !reflect.DeepEqual(config.MapInt64Uint32Val, map[int64]uint32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Uint32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Uint32Val)
	}
	if !reflect.DeepEqual(config.MapInt64Uint64Val, map[int64]uint64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Uint64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Uint64Val)
	}

	if !reflect.DeepEqual(config.MapInt64IntVal, map[int64]int{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64IntVal should be 'map[2:2 3:3]', got: '%v'", config.MapInt64IntVal)
	}
	if !reflect.DeepEqual(config.MapInt64Int8Val, map[int64]int8{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Int8Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Int8Val)
	}
	if !reflect.DeepEqual(config.MapInt64Int16Val, map[int64]int16{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Int16Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Int16Val)
	}
	if !reflect.DeepEqual(config.MapInt64Int32Val, map[int64]int32{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Int32Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Int32Val)
	}
	if !reflect.DeepEqual(config.MapInt64Int64Val, map[int64]int64{2: 2, 3: 3}) {
		t.Errorf("config.MapInt64Int64Val should be 'map[2:2 3:3]', got: '%v'", config.MapInt64Int64Val)
	}

	if !reflect.DeepEqual(config.MapInt64Float32Val, map[int64]float32{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt64Float32Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt64Float32Val)
	}
	if !reflect.DeepEqual(config.MapInt64Float64Val, map[int64]float64{2: 2.0, 3: 3.0}) {
		t.Errorf("config.MapInt64Float64Val should be 'map[2:2.0 3:3.0]', got: '%v'", config.MapInt64Float64Val)
	}

	if !reflect.DeepEqual(config.MapFloat32BoolVal, map[float32]bool{2.0: true, 3.0: true}) {
		t.Errorf("config.MapFloat32BoolVal should be 'map[2.0:true 3.0:true], got: '%v'", config.MapFloat32BoolVal)
	}

	if !reflect.DeepEqual(config.MapFloat32StringVal, map[float32]string{2.0: "three", 3.0: "four"}) {
		t.Errorf("config.MapFloat32StringVal should be 'map[2.0:three 3.0:four]', got: '%v'", config.MapFloat32StringVal)
	}

	if !reflect.DeepEqual(config.MapFloat32UintVal, map[float32]uint{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32UintVal should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32UintVal)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint8Val, map[float32]uint8{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Uint8Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Uint8Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint16Val, map[float32]uint16{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Uint16Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Uint16Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint32Val, map[float32]uint32{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Uint32Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Uint32Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Uint64Val, map[float32]uint64{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Uint64Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Uint64Val)
	}

	if !reflect.DeepEqual(config.MapFloat32IntVal, map[float32]int{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32IntVal should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32IntVal)
	}
	if !reflect.DeepEqual(config.MapFloat32Int8Val, map[float32]int8{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Int8Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Int8Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Int16Val, map[float32]int16{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Int16Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Int16Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Int32Val, map[float32]int32{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Int32Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Int32Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Int64Val, map[float32]int64{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat32Int64Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat32Int64Val)
	}

	if !reflect.DeepEqual(config.MapFloat32Float32Val, map[float32]float32{2.0: 2.0, 3.0: 3.0}) {
		t.Errorf("config.MapFloat32Float32Val should be 'map[2.0:2.0 3.0:3.0]', got: '%v'", config.MapFloat32Float32Val)
	}
	if !reflect.DeepEqual(config.MapFloat32Float64Val, map[float32]float64{2.0: 2.0, 3.0: 3.0}) {
		t.Errorf("config.MapFloat32Float64Val should be 'map[2.0:2.0 3.0:3.0]', got: '%v'", config.MapFloat32Float64Val)
	}

	if !reflect.DeepEqual(config.MapFloat64BoolVal, map[float64]bool{2.0: true, 3.0: true}) {
		t.Errorf("config.MapFloat64BoolVal should be 'map[2.0:true 3.0:true], got: '%v'", config.MapFloat64BoolVal)
	}

	if !reflect.DeepEqual(config.MapFloat64StringVal, map[float64]string{2.0: "three", 3.0: "four"}) {
		t.Errorf("config.MapFloat64StringVal should be 'map[2.0:three 3.0:four]', got: '%v'", config.MapFloat64StringVal)
	}

	if !reflect.DeepEqual(config.MapFloat64UintVal, map[float64]uint{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64UintVal should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64UintVal)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint8Val, map[float64]uint8{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Uint8Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Uint8Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint16Val, map[float64]uint16{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Uint16Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Uint16Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint32Val, map[float64]uint32{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Uint32Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Uint32Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Uint64Val, map[float64]uint64{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Uint64Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Uint64Val)
	}

	if !reflect.DeepEqual(config.MapFloat64IntVal, map[float64]int{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64IntVal should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64IntVal)
	}
	if !reflect.DeepEqual(config.MapFloat64Int8Val, map[float64]int8{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Int8Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Int8Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Int16Val, map[float64]int16{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Int16Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Int16Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Int32Val, map[float64]int32{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Int32Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Int32Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Int64Val, map[float64]int64{2.0: 2, 3.0: 3}) {
		t.Errorf("config.MapFloat64Int64Val should be 'map[2.0:2 3.0:3]', got: '%v'", config.MapFloat64Int64Val)
	}

	if !reflect.DeepEqual(config.MapFloat64Float32Val, map[float64]float32{2.0: 2.0, 3.0: 3.0}) {
		t.Errorf("config.MapFloat64Float32Val should be 'map[2.0:2.0 3.0:3.0]', got: '%v'", config.MapFloat64Float32Val)
	}
	if !reflect.DeepEqual(config.MapFloat64Float64Val, map[float64]float64{2.0: 2.0, 3.0: 3.0}) {
		t.Errorf("config.MapFloat64Float64Val should be 'map[2.0:2.0 3.0:3.0]', got: '%v'", config.MapFloat64Float64Val)
	}
}

// mapTypes is an test struct that holds parameters of different map types
type mapTypes struct {
	MapBoolBoolVal    map[bool]bool    `env:"ENV_MAP_BOOL_BOOL" default:"true:true,false:false"`
	MapBoolStringVal  map[bool]string  `env:"ENV_MAP_BOOL_STRING" default:"true:s1,false:s2"`
	MapBoolUintVal    map[bool]uint    `env:"ENV_MAP_BOOL_UINT" default:"true:1,false:2"`
	MapBoolUint8Val   map[bool]uint8   `env:"ENV_MAP_BOOL_UINT8" default:"true:1,false:2"`
	MapBoolUint16Val  map[bool]uint16  `env:"ENV_MAP_BOOL_UINT16" default:"true:1,false:2"`
	MapBoolUint32Val  map[bool]uint32  `env:"ENV_MAP_BOOL_UINT32" default:"true:1,false:2"`
	MapBoolUint64Val  map[bool]uint64  `env:"ENV_MAP_BOOL_UINT64" default:"true:1,false:2"`
	MapBoolIntVal     map[bool]int     `env:"ENV_MAP_BOOL_INT" default:"true:1,false:2"`
	MapBoolInt8Val    map[bool]int8    `env:"ENV_MAP_BOOL_INT8" default:"true:1,false:2"`
	MapBoolInt16Val   map[bool]int16   `env:"ENV_MAP_BOOL_INT16" default:"true:1,false:2"`
	MapBoolInt32Val   map[bool]int32   `env:"ENV_MAP_BOOL_INT32" default:"true:1,false:2"`
	MapBoolInt64Val   map[bool]int64   `env:"ENV_MAP_BOOL_INT64" default:"true:1,false:2"`
	MapBoolFloat32Val map[bool]float32 `env:"ENV_MAP_BOOL_FLOAT32" default:"true:1.0,false:2.0"`
	MapBoolFloat64Val map[bool]float64 `env:"ENV_MAP_BOOL_FLOAT64" default:"true:1.0,false:2.0"`

	MapStringBoolVal    map[string]bool    `env:"ENV_MAP_STRING_BOOL" default:"s1:true,s2:false"`
	MapStringStringVal  map[string]string  `env:"ENV_MAP_STRING_STRING" default:"s1:s1,s2:s2"`
	MapStringUintVal    map[string]uint    `env:"ENV_MAP_STRING_UINT" default:"s1:1,s2:2"`
	MapStringUint8Val   map[string]uint8   `env:"ENV_MAP_STRING_UINT8" default:"s1:1,s2:2"`
	MapStringUint16Val  map[string]uint16  `env:"ENV_MAP_STRING_UINT16" default:"s1:1,s2:2"`
	MapStringUint32Val  map[string]uint32  `env:"ENV_MAP_STRING_UINT32" default:"s1:1,s2:2"`
	MapStringUint64Val  map[string]uint64  `env:"ENV_MAP_STRING_UINT64" default:"s1:1,s2:2"`
	MapStringIntVal     map[string]int     `env:"ENV_MAP_STRING_INT" default:"s1:1,s2:2"`
	MapStringInt8Val    map[string]int8    `env:"ENV_MAP_STRING_INT8" default:"s1:1,s2:2"`
	MapStringInt16Val   map[string]int16   `env:"ENV_MAP_STRING_INT16" default:"s1:1,s2:2"`
	MapStringInt32Val   map[string]int32   `env:"ENV_MAP_STRING_INT32" default:"s1:1,s2:2"`
	MapStringInt64Val   map[string]int64   `env:"ENV_MAP_STRING_INT64" default:"s1:1,s2:2"`
	MapStringFloat32Val map[string]float32 `env:"ENV_MAP_STRING_FLOAT32" default:"s1:1.0,s2:2.0"`
	MapStringFloat64Val map[string]float64 `env:"ENV_MAP_STRING_FLOAT64" default:"s1:1.0,s2:2.0"`

	MapUintBoolVal    map[uint]bool    `env:"ENV_MAP_UINT_BOOL" default:"1:true,2:false"`
	MapUintStringVal  map[uint]string  `env:"ENV_MAP_UINT_STRING" default:"1:s1,2:s2"`
	MapUintUintVal    map[uint]uint    `env:"ENV_MAP_UINT_UINT" default:"1:1,2:2"`
	MapUintUint8Val   map[uint]uint8   `env:"ENV_MAP_UINT_UINT8" default:"1:1,2:2"`
	MapUintUint16Val  map[uint]uint16  `env:"ENV_MAP_UINT_UINT16" default:"1:1,2:2"`
	MapUintUint32Val  map[uint]uint32  `env:"ENV_MAP_UINT_UINT32" default:"1:1,2:2"`
	MapUintUint64Val  map[uint]uint64  `env:"ENV_MAP_UINT_UINT64" default:"1:1,2:2"`
	MapUintIntVal     map[uint]int     `env:"ENV_MAP_UINT_INT" default:"1:1,2:2"`
	MapUintInt8Val    map[uint]int8    `env:"ENV_MAP_UINT_INT8" default:"1:1,2:2"`
	MapUintInt16Val   map[uint]int16   `env:"ENV_MAP_UINT_INT16" default:"1:1,2:2"`
	MapUintInt32Val   map[uint]int32   `env:"ENV_MAP_UINT_INT32" default:"1:1,2:2"`
	MapUintInt64Val   map[uint]int64   `env:"ENV_MAP_UINT_INT64" default:"1:1,2:2"`
	MapUintFloat32Val map[uint]float32 `env:"ENV_MAP_UINT_FLOAT32" default:"1:1.0,2:2.0"`
	MapUintFloat64Val map[uint]float64 `env:"ENV_MAP_UINT_FLOAT64" default:"1:1.0,2:2.0"`

	MapUint8BoolVal    map[uint8]bool    `env:"ENV_MAP_UINT8_BOOL" default:"1:true,2:false"`
	MapUint8StringVal  map[uint8]string  `env:"ENV_MAP_UINT8_STRING" default:"1:s1,2:s2"`
	MapUint8UintVal    map[uint8]uint    `env:"ENV_MAP_UINT8_UINT" default:"1:1,2:2"`
	MapUint8Uint8Val   map[uint8]uint8   `env:"ENV_MAP_UINT8_UINT8" default:"1:1,2:2"`
	MapUint8Uint16Val  map[uint8]uint16  `env:"ENV_MAP_UINT8_UINT16" default:"1:1,2:2"`
	MapUint8Uint32Val  map[uint8]uint32  `env:"ENV_MAP_UINT8_UINT32" default:"1:1,2:2"`
	MapUint8Uint64Val  map[uint8]uint64  `env:"ENV_MAP_UINT8_UINT64" default:"1:1,2:2"`
	MapUint8IntVal     map[uint8]int     `env:"ENV_MAP_UINT8_INT" default:"1:1,2:2"`
	MapUint8Int8Val    map[uint8]int8    `env:"ENV_MAP_UINT8_INT8" default:"1:1,2:2"`
	MapUint8Int16Val   map[uint8]int16   `env:"ENV_MAP_UINT8_INT16" default:"1:1,2:2"`
	MapUint8Int32Val   map[uint8]int32   `env:"ENV_MAP_UINT8_INT32" default:"1:1,2:2"`
	MapUint8Int64Val   map[uint8]int64   `env:"ENV_MAP_UINT8_INT64" default:"1:1,2:2"`
	MapUint8Float32Val map[uint8]float32 `env:"ENV_MAP_UINT8_FLOAT32" default:"1:1.0,2:2.0"`
	MapUint8Float64Val map[uint8]float64 `env:"ENV_MAP_UINT8_FLOAT64" default:"1:1.0,2:2.0"`

	MapUint16BoolVal    map[uint16]bool    `env:"ENV_MAP_UINT16_BOOL" default:"1:true,2:false"`
	MapUint16StringVal  map[uint16]string  `env:"ENV_MAP_UINT16_STRING" default:"1:s1,2:s2"`
	MapUint16UintVal    map[uint16]uint    `env:"ENV_MAP_UINT16_UINT" default:"1:1,2:2"`
	MapUint16Uint8Val   map[uint16]uint8   `env:"ENV_MAP_UINT16_UINT8" default:"1:1,2:2"`
	MapUint16Uint16Val  map[uint16]uint16  `env:"ENV_MAP_UINT16_UINT16" default:"1:1,2:2"`
	MapUint16Uint32Val  map[uint16]uint32  `env:"ENV_MAP_UINT16_UINT32" default:"1:1,2:2"`
	MapUint16Uint64Val  map[uint16]uint64  `env:"ENV_MAP_UINT16_UINT64" default:"1:1,2:2"`
	MapUint16IntVal     map[uint16]int     `env:"ENV_MAP_UINT16_INT" default:"1:1,2:2"`
	MapUint16Int8Val    map[uint16]int8    `env:"ENV_MAP_UINT16_INT8" default:"1:1,2:2"`
	MapUint16Int16Val   map[uint16]int16   `env:"ENV_MAP_UINT16_INT16" default:"1:1,2:2"`
	MapUint16Int32Val   map[uint16]int32   `env:"ENV_MAP_UINT16_INT32" default:"1:1,2:2"`
	MapUint16Int64Val   map[uint16]int64   `env:"ENV_MAP_UINT16_INT64" default:"1:1,2:2"`
	MapUint16Float32Val map[uint16]float32 `env:"ENV_MAP_UINT16_FLOAT32" default:"1:1.0,2:2.0"`
	MapUint16Float64Val map[uint16]float64 `env:"ENV_MAP_UINT16_FLOAT64" default:"1:1.0,2:2.0"`

	MapUint32BoolVal    map[uint32]bool    `env:"ENV_MAP_UINT32_BOOL" default:"1:true,2:false"`
	MapUint32StringVal  map[uint32]string  `env:"ENV_MAP_UINT32_STRING" default:"1:s1,2:s2"`
	MapUint32UintVal    map[uint32]uint    `env:"ENV_MAP_UINT32_UINT" default:"1:1,2:2"`
	MapUint32Uint8Val   map[uint32]uint8   `env:"ENV_MAP_UINT32_UINT8" default:"1:1,2:2"`
	MapUint32Uint16Val  map[uint32]uint16  `env:"ENV_MAP_UINT32_UINT16" default:"1:1,2:2"`
	MapUint32Uint32Val  map[uint32]uint32  `env:"ENV_MAP_UINT32_UINT32" default:"1:1,2:2"`
	MapUint32Uint64Val  map[uint32]uint64  `env:"ENV_MAP_UINT32_UINT64" default:"1:1,2:2"`
	MapUint32IntVal     map[uint32]int     `env:"ENV_MAP_UINT32_INT" default:"1:1,2:2"`
	MapUint32Int8Val    map[uint32]int8    `env:"ENV_MAP_UINT32_INT8" default:"1:1,2:2"`
	MapUint32Int16Val   map[uint32]int16   `env:"ENV_MAP_UINT32_INT16" default:"1:1,2:2"`
	MapUint32Int32Val   map[uint32]int32   `env:"ENV_MAP_UINT32_INT32" default:"1:1,2:2"`
	MapUint32Int64Val   map[uint32]int64   `env:"ENV_MAP_UINT32_INT64" default:"1:1,2:2"`
	MapUint32Float32Val map[uint32]float32 `env:"ENV_MAP_UINT32_FLOAT32" default:"1:1.0,2:2.0"`
	MapUint32Float64Val map[uint32]float64 `env:"ENV_MAP_UINT32_FLOAT64" default:"1:1.0,2:2.0"`

	MapUint64BoolVal    map[uint64]bool    `env:"ENV_MAP_UINT64_BOOL" default:"1:true,2:false"`
	MapUint64StringVal  map[uint64]string  `env:"ENV_MAP_UINT64_STRING" default:"1:s1,2:s2"`
	MapUint64UintVal    map[uint64]uint    `env:"ENV_MAP_UINT64_UINT" default:"1:1,2:2"`
	MapUint64Uint8Val   map[uint64]uint8   `env:"ENV_MAP_UINT64_UINT8" default:"1:1,2:2"`
	MapUint64Uint16Val  map[uint64]uint16  `env:"ENV_MAP_UINT64_UINT16" default:"1:1,2:2"`
	MapUint64Uint32Val  map[uint64]uint32  `env:"ENV_MAP_UINT64_UINT32" default:"1:1,2:2"`
	MapUint64Uint64Val  map[uint64]uint64  `env:"ENV_MAP_UINT64_UINT64" default:"1:1,2:2"`
	MapUint64IntVal     map[uint64]int     `env:"ENV_MAP_UINT64_INT" default:"1:1,2:2"`
	MapUint64Int8Val    map[uint64]int8    `env:"ENV_MAP_UINT64_INT8" default:"1:1,2:2"`
	MapUint64Int16Val   map[uint64]int16   `env:"ENV_MAP_UINT64_INT16" default:"1:1,2:2"`
	MapUint64Int32Val   map[uint64]int32   `env:"ENV_MAP_UINT64_INT32" default:"1:1,2:2"`
	MapUint64Int64Val   map[uint64]int64   `env:"ENV_MAP_UINT64_INT64" default:"1:1,2:2"`
	MapUint64Float32Val map[uint64]float32 `env:"ENV_MAP_UINT64_FLOAT32" default:"1:1.0,2:2.0"`
	MapUint64Float64Val map[uint64]float64 `env:"ENV_MAP_UINT64_FLOAT64" default:"1:1.0,2:2.0"`

	MapIntBoolVal    map[int]bool    `env:"ENV_MAP_INT_BOOL" default:"1:true,2:false"`
	MapIntStringVal  map[int]string  `env:"ENV_MAP_INT_STRING" default:"1:s1,2:s2"`
	MapIntUintVal    map[int]uint    `env:"ENV_MAP_INT_UINT" default:"1:1,2:2"`
	MapIntUint8Val   map[int]uint8   `env:"ENV_MAP_INT_UINT8" default:"1:1,2:2"`
	MapIntUint16Val  map[int]uint16  `env:"ENV_MAP_INT_UINT16" default:"1:1,2:2"`
	MapIntUint32Val  map[int]uint32  `env:"ENV_MAP_INT_UINT32" default:"1:1,2:2"`
	MapIntUint64Val  map[int]uint64  `env:"ENV_MAP_INT_UINT64" default:"1:1,2:2"`
	MapIntIntVal     map[int]int     `env:"ENV_MAP_INT_INT" default:"1:1,2:2"`
	MapIntInt8Val    map[int]int8    `env:"ENV_MAP_INT_INT8" default:"1:1,2:2"`
	MapIntInt16Val   map[int]int16   `env:"ENV_MAP_INT_INT16" default:"1:1,2:2"`
	MapIntInt32Val   map[int]int32   `env:"ENV_MAP_INT_INT32" default:"1:1,2:2"`
	MapIntInt64Val   map[int]int64   `env:"ENV_MAP_INT_INT64" default:"1:1,2:2"`
	MapIntFloat32Val map[int]float32 `env:"ENV_MAP_INT_FLOAT32" default:"1:1.0,2:2.0"`
	MapIntFloat64Val map[int]float64 `env:"ENV_MAP_INT_FLOAT64" default:"1:1.0,2:2.0"`

	MapInt8BoolVal    map[int8]bool    `env:"ENV_MAP_INT8_BOOL" default:"1:true,2:false"`
	MapInt8StringVal  map[int8]string  `env:"ENV_MAP_INT8_STRING" default:"1:s1,2:s2"`
	MapInt8UintVal    map[int8]uint    `env:"ENV_MAP_INT8_UINT" default:"1:1,2:2"`
	MapInt8Uint8Val   map[int8]uint8   `env:"ENV_MAP_INT8_UINT8" default:"1:1,2:2"`
	MapInt8Uint16Val  map[int8]uint16  `env:"ENV_MAP_INT8_UINT16" default:"1:1,2:2"`
	MapInt8Uint32Val  map[int8]uint32  `env:"ENV_MAP_INT8_UINT32" default:"1:1,2:2"`
	MapInt8Uint64Val  map[int8]uint64  `env:"ENV_MAP_INT8_UINT64" default:"1:1,2:2"`
	MapInt8IntVal     map[int8]int     `env:"ENV_MAP_INT8_INT" default:"1:1,2:2"`
	MapInt8Int8Val    map[int8]int8    `env:"ENV_MAP_INT8_INT8" default:"1:1,2:2"`
	MapInt8Int16Val   map[int8]int16   `env:"ENV_MAP_INT8_INT16" default:"1:1,2:2"`
	MapInt8Int32Val   map[int8]int32   `env:"ENV_MAP_INT8_INT32" default:"1:1,2:2"`
	MapInt8Int64Val   map[int8]int64   `env:"ENV_MAP_INT8_INT64" default:"1:1,2:2"`
	MapInt8Float32Val map[int8]float32 `env:"ENV_MAP_INT8_FLOAT32" default:"1:1.0,2:2.0"`
	MapInt8Float64Val map[int8]float64 `env:"ENV_MAP_INT8_FLOAT64" default:"1:1.0,2:2.0"`

	MapInt16BoolVal    map[int16]bool    `env:"ENV_MAP_INT16_BOOL" default:"1:true,2:false"`
	MapInt16StringVal  map[int16]string  `env:"ENV_MAP_INT16_STRING" default:"1:s1,2:s2"`
	MapInt16UintVal    map[int16]uint    `env:"ENV_MAP_INT16_UINT" default:"1:1,2:2"`
	MapInt16Uint8Val   map[int16]uint8   `env:"ENV_MAP_INT16_UINT8" default:"1:1,2:2"`
	MapInt16Uint16Val  map[int16]uint16  `env:"ENV_MAP_INT16_UINT16" default:"1:1,2:2"`
	MapInt16Uint32Val  map[int16]uint32  `env:"ENV_MAP_INT16_UINT32" default:"1:1,2:2"`
	MapInt16Uint64Val  map[int16]uint64  `env:"ENV_MAP_INT16_UINT64" default:"1:1,2:2"`
	MapInt16IntVal     map[int16]int     `env:"ENV_MAP_INT16_INT" default:"1:1,2:2"`
	MapInt16Int8Val    map[int16]int8    `env:"ENV_MAP_INT16_INT8" default:"1:1,2:2"`
	MapInt16Int16Val   map[int16]int16   `env:"ENV_MAP_INT16_INT16" default:"1:1,2:2"`
	MapInt16Int32Val   map[int16]int32   `env:"ENV_MAP_INT16_INT32" default:"1:1,2:2"`
	MapInt16Int64Val   map[int16]int64   `env:"ENV_MAP_INT16_INT64" default:"1:1,2:2"`
	MapInt16Float32Val map[int16]float32 `env:"ENV_MAP_INT16_FLOAT32" default:"1:1.0,2:2.0"`
	MapInt16Float64Val map[int16]float64 `env:"ENV_MAP_INT16_FLOAT64" default:"1:1.0,2:2.0"`

	MapInt32BoolVal    map[int32]bool    `env:"ENV_MAP_INT32_BOOL" default:"1:true,2:false"`
	MapInt32StringVal  map[int32]string  `env:"ENV_MAP_INT32_STRING" default:"1:s1,2:s2"`
	MapInt32UintVal    map[int32]uint    `env:"ENV_MAP_INT32_UINT" default:"1:1,2:2"`
	MapInt32Uint8Val   map[int32]uint8   `env:"ENV_MAP_INT32_UINT8" default:"1:1,2:2"`
	MapInt32Uint16Val  map[int32]uint16  `env:"ENV_MAP_INT32_UINT16" default:"1:1,2:2"`
	MapInt32Uint32Val  map[int32]uint32  `env:"ENV_MAP_INT32_UINT32" default:"1:1,2:2"`
	MapInt32Uint64Val  map[int32]uint64  `env:"ENV_MAP_INT32_UINT64" default:"1:1,2:2"`
	MapInt32IntVal     map[int32]int     `env:"ENV_MAP_INT32_INT" default:"1:1,2:2"`
	MapInt32Int8Val    map[int32]int8    `env:"ENV_MAP_INT32_INT8" default:"1:1,2:2"`
	MapInt32Int16Val   map[int32]int16   `env:"ENV_MAP_INT32_INT16" default:"1:1,2:2"`
	MapInt32Int32Val   map[int32]int32   `env:"ENV_MAP_INT32_INT32" default:"1:1,2:2"`
	MapInt32Int64Val   map[int32]int64   `env:"ENV_MAP_INT32_INT64" default:"1:1,2:2"`
	MapInt32Float32Val map[int32]float32 `env:"ENV_MAP_INT32_FLOAT32" default:"1:1.0,2:2.0"`
	MapInt32Float64Val map[int32]float64 `env:"ENV_MAP_INT32_FLOAT64" default:"1:1.0,2:2.0"`

	MapInt64BoolVal    map[int64]bool    `env:"ENV_MAP_INT64_BOOL" default:"1:true,2:false"`
	MapInt64StringVal  map[int64]string  `env:"ENV_MAP_INT64_STRING" default:"1:s1,2:s2"`
	MapInt64UintVal    map[int64]uint    `env:"ENV_MAP_INT64_UINT" default:"1:1,2:2"`
	MapInt64Uint8Val   map[int64]uint8   `env:"ENV_MAP_INT64_UINT8" default:"1:1,2:2"`
	MapInt64Uint16Val  map[int64]uint16  `env:"ENV_MAP_INT64_UINT16" default:"1:1,2:2"`
	MapInt64Uint32Val  map[int64]uint32  `env:"ENV_MAP_INT64_UINT32" default:"1:1,2:2"`
	MapInt64Uint64Val  map[int64]uint64  `env:"ENV_MAP_INT64_UINT64" default:"1:1,2:2"`
	MapInt64IntVal     map[int64]int     `env:"ENV_MAP_INT64_INT" default:"1:1,2:2"`
	MapInt64Int8Val    map[int64]int8    `env:"ENV_MAP_INT64_INT8" default:"1:1,2:2"`
	MapInt64Int16Val   map[int64]int16   `env:"ENV_MAP_INT64_INT16" default:"1:1,2:2"`
	MapInt64Int32Val   map[int64]int32   `env:"ENV_MAP_INT64_INT32" default:"1:1,2:2"`
	MapInt64Int64Val   map[int64]int64   `env:"ENV_MAP_INT64_INT64" default:"1:1,2:2"`
	MapInt64Float32Val map[int64]float32 `env:"ENV_MAP_INT64_FLOAT32" default:"1:1.0,2:2.0"`
	MapInt64Float64Val map[int64]float64 `env:"ENV_MAP_INT64_FLOAT64" default:"1:1.0,2:2.0"`

	MapFloat32BoolVal    map[float32]bool    `env:"ENV_MAP_FLOAT32_BOOL" default:"1.0:true,2.0:false"`
	MapFloat32StringVal  map[float32]string  `env:"ENV_MAP_FLOAT32_STRING" default:"1.0:s1,2.0:s2"`
	MapFloat32UintVal    map[float32]uint    `env:"ENV_MAP_FLOAT32_UINT" default:"1.0:1,2.0:2"`
	MapFloat32Uint8Val   map[float32]uint8   `env:"ENV_MAP_FLOAT32_UINT8" default:"1.0:1,2.0:2"`
	MapFloat32Uint16Val  map[float32]uint16  `env:"ENV_MAP_FLOAT32_UINT16" default:"1.0:1,2.0:2"`
	MapFloat32Uint32Val  map[float32]uint32  `env:"ENV_MAP_FLOAT32_UINT32" default:"1.0:1,2.0:2"`
	MapFloat32Uint64Val  map[float32]uint64  `env:"ENV_MAP_FLOAT32_UINT64" default:"1.0:1,2.0:2"`
	MapFloat32IntVal     map[float32]int     `env:"ENV_MAP_FLOAT32_INT" default:"1.0:1,2.0:2"`
	MapFloat32Int8Val    map[float32]int8    `env:"ENV_MAP_FLOAT32_INT8" default:"1.0:1,2.0:2"`
	MapFloat32Int16Val   map[float32]int16   `env:"ENV_MAP_FLOAT32_INT16" default:"1.0:1,2.0:2"`
	MapFloat32Int32Val   map[float32]int32   `env:"ENV_MAP_FLOAT32_INT32" default:"1.0:1,2.0:2"`
	MapFloat32Int64Val   map[float32]int64   `env:"ENV_MAP_FLOAT32_INT64" default:"1.0:1,2.0:2"`
	MapFloat32Float32Val map[float32]float32 `env:"ENV_MAP_FLOAT32_FLOAT32" default:"1.0:1.0,2.0:2.0"`
	MapFloat32Float64Val map[float32]float64 `env:"ENV_MAP_FLOAT32_FLOAT64" default:"1.0:1.0,2.0:2.0"`

	MapFloat64BoolVal    map[float64]bool    `env:"ENV_MAP_FLOAT64_BOOL" default:"1.0:true,2.0:false"`
	MapFloat64StringVal  map[float64]string  `env:"ENV_MAP_FLOAT64_STRING" default:"1.0:s1,2.0:s2"`
	MapFloat64UintVal    map[float64]uint    `env:"ENV_MAP_FLOAT64_UINT" default:"1.0:1,2.0:2"`
	MapFloat64Uint8Val   map[float64]uint8   `env:"ENV_MAP_FLOAT64_UINT8" default:"1.0:1,2.0:2"`
	MapFloat64Uint16Val  map[float64]uint16  `env:"ENV_MAP_FLOAT64_UINT16" default:"1.0:1,2.0:2"`
	MapFloat64Uint32Val  map[float64]uint32  `env:"ENV_MAP_FLOAT64_UINT32" default:"1.0:1,2.0:2"`
	MapFloat64Uint64Val  map[float64]uint64  `env:"ENV_MAP_FLOAT64_UINT64" default:"1.0:1,2.0:2"`
	MapFloat64IntVal     map[float64]int     `env:"ENV_MAP_FLOAT64_INT" default:"1.0:1,2.0:2"`
	MapFloat64Int8Val    map[float64]int8    `env:"ENV_MAP_FLOAT64_INT8" default:"1.0:1,2.0:2"`
	MapFloat64Int16Val   map[float64]int16   `env:"ENV_MAP_FLOAT64_INT16" default:"1.0:1,2.0:2"`
	MapFloat64Int32Val   map[float64]int32   `env:"ENV_MAP_FLOAT64_INT32" default:"1.0:1,2.0:2"`
	MapFloat64Int64Val   map[float64]int64   `env:"ENV_MAP_FLOAT64_INT64" default:"1.0:1,2.0:2"`
	MapFloat64Float32Val map[float64]float32 `env:"ENV_MAP_FLOAT64_FLOAT32" default:"1.0:1.0,2.0:2.0"`
	MapFloat64Float64Val map[float64]float64 `env:"ENV_MAP_FLOAT64_FLOAT64" default:"1.0:1.0,2.0:2.0"`
}
