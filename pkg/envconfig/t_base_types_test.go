package envconfig

import (
	"os"
	"testing"
)

func TestBaseTypesDefaults(t *testing.T) {
	os.Clearenv()
	config := new(baseTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if !config.BoolVal {
		t.Errorf("config.BoolVal should be 'true', got: %t", config.BoolVal)
	}
	if config.StringVal != "string_value" {
		t.Errorf("config.StringValue should be 'string_value', got: %q", config.StringVal)
	}

	if config.UintVal != uint(1) {
		t.Errorf("config.UintValue should be '1', got: %q", config.UintVal)
	}
	if config.Uint8Val != uint8(1) {
		t.Errorf("config.Uint8Value should be '1', got: %q", config.Uint8Val)
	}
	if config.Uint16Val != uint16(1) {
		t.Errorf("config.Uint16Value should be '1', got: %q", config.Uint16Val)
	}
	if config.Uint32Val != uint32(1) {
		t.Errorf("config.Uint32Value should be '1', got: %q", config.Uint32Val)
	}
	if config.Uint64Val != uint64(1) {
		t.Errorf("config.Uint64Value should be '1', got: %q", config.Uint64Val)
	}

	if config.IntVal != 1 {
		t.Errorf("config.IntValue should be '1', got: %q", config.IntVal)
	}
	if config.Int8Val != 1 {
		t.Errorf("config.Int8Value should be '1', got: %q", config.Int8Val)
	}
	if config.Int16Val != 1 {
		t.Errorf("config.Int16Value should be '1', got: %q", config.Int16Val)
	}
	if config.Int32Val != 1 {
		t.Errorf("config.Int32Value should be '1', got: %q", config.Int32Val)
	}
	if config.Int64Val != 1 {
		t.Errorf("config.Int64Value should be '1', got: %q", config.Int64Val)
	}

	if config.Float32Val != float32(1) {
		t.Errorf("config.Float32Value should be '1.0', got: '%.1f'", config.Float32Val)
	}
	if config.Float64Val != float64(1) {
		t.Errorf("config.Float64Value should be '1.0', got: '%.1f'", config.Float64Val)
	}
}

func TestBaseTypes(t *testing.T) {
	os.Clearenv()
	envs := envValues{
		"ENV_BOOL":    envValue{checkValue: true},
		"ENV_STRING":  envValue{checkValue: "set string"},
		"ENV_UINT":    envValue{checkValue: uint(2)},
		"ENV_UINT8":   envValue{checkValue: uint(2)},
		"ENV_UINT16":  envValue{checkValue: uint(2)},
		"ENV_UINT32":  envValue{checkValue: uint(2)},
		"ENV_UINT64":  envValue{checkValue: uint(2)},
		"ENV_INT":     envValue{checkValue: int(2)},
		"ENV_INT8":    envValue{checkValue: int8(2)},
		"ENV_INT16":   envValue{checkValue: int16(2)},
		"ENV_INT32":   envValue{checkValue: int32(2)},
		"ENV_INT64":   envValue{checkValue: int64(2)},
		"ENV_FLOAT32": envValue{checkValue: float32(2.0)},
		"ENV_FLOAT64": envValue{checkValue: float64(2.0)},
	}
	if err := envs.set(); err != nil {
		t.Fatalf("envs.set() fail, raw: %v", err)
	}
	config := new(baseTypes)
	// envs.getDefaultValues(config)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if !config.BoolVal {
		t.Errorf("config.BoolVal should be 'true', got: %t", config.BoolVal)
	}
	if config.StringVal != "set string" {
		t.Errorf("config.StringValue should be 'set string', got: %q", config.StringVal)
	}

	if config.UintVal != uint(2) {
		t.Errorf("config.UintValue should be '2', got: '%d'", config.UintVal)
	}
	if config.Uint8Val != uint8(2) {
		t.Errorf("config.Uint8Value should be '2', got: '%d'", config.Uint8Val)
	}
	if config.Uint16Val != uint16(2) {
		t.Errorf("config.Uint16Value should be '2', got: '%d'", config.Uint16Val)
	}
	if config.Uint32Val != uint32(2) {
		t.Errorf("config.Uint32Value should be '2', got: '%d'", config.Uint32Val)
	}
	if config.Uint64Val != uint64(2) {
		t.Errorf("config.Uint64Value should be '2', got: '%d'", config.Uint64Val)
	}

	if config.IntVal != 2 {
		t.Errorf("config.IntValue should be '2', got: %q", config.IntVal)
	}
	if config.Int8Val != 2 {
		t.Errorf("config.Int8Value should be '2', got: %q", config.Int8Val)
	}
	if config.Int16Val != 2 {
		t.Errorf("config.Int16Value should be '2', got: %q", config.Int16Val)
	}
	if config.Int32Val != 2 {
		t.Errorf("config.Int32Value should be '2', got: %q", config.Int32Val)
	}
	if config.Int64Val != 2 {
		t.Errorf("config.Int64Value should be '2', got: %q", config.Int64Val)
	}

	if config.Float32Val != float32(2.0) {
		t.Errorf("config.Float32Value should be '2.0', got: '%.1f'", config.Float32Val)
	}
	if config.Float64Val != float64(2.0) {
		t.Errorf("config.Float64Value should be '2.0', got: '%.1f'", config.Float64Val)
	}
}

// baseTypes is an test struct that holds parameters of different base types
// such as bool, string, integers (included unsigned integers), floats
type baseTypes struct {
	// base types
	BoolVal    bool    `env:"ENV_BOOL" default:"true"`
	StringVal  string  `env:"ENV_STRING" default:"string_value"`
	UintVal    uint    `env:"ENV_UINT" default:"1"`
	Uint8Val   uint8   `env:"ENV_UINT8" default:"1"`
	Uint16Val  uint16  `env:"ENV_UINT16" default:"1"`
	Uint32Val  uint32  `env:"ENV_UINT32" default:"1"`
	Uint64Val  uint64  `env:"ENV_UINT64" default:"1"`
	IntVal     int     `env:"ENV_INT" default:"1"`
	Int8Val    int8    `env:"ENV_INT8" default:"1"`
	Int16Val   int16   `env:"ENV_INT16" default:"1"`
	Int32Val   int32   `env:"ENV_INT32" default:"1"`
	Int64Val   int64   `env:"ENV_INT64" default:"1"`
	Float32Val float32 `env:"ENV_FLOAT32" default:"1.0"`
	Float64Val float64 `env:"ENV_FLOAT64" default:"1.0"`
}
