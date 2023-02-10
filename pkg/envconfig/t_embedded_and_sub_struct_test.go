package envconfig

import (
	"os"
	"testing"
)

func TestEmbeddedAndSubStructsDefaults(t *testing.T) {
	os.Clearenv()
	config := new(embeddedAndSubStruct)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}
	if config.StructVal.StructValStringVal != "sub_string" {
		t.Errorf("config.StructVal.StructValStringVal should be 'sub_string', got: %q", config.StructVal.StructValStringVal)
	}
	if config.EmbeddedStructValStringVal != "embedded_string" {
		t.Errorf("config.EmbeddedStructValStringVal  should be 'embedded_string', got: %q", config.EmbeddedStructValStringVal)
	}
}

func TestEmbeddeAndSubStructs(t *testing.T) {
	os.Clearenv()
	envs := envValues{
		"ENV_SUB_STRING_VAL":      envValue{checkValue: "set sub string"},
		"ENV_EMBEDDED_STRING_VAL": envValue{checkValue: "set embedded string"},
	}
	if err := envs.set(); err != nil {
		t.Fatalf("envs.set() fail, raw: %v", err)
	}
	config := new(embeddedAndSubStruct)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}
	if config.StructVal.StructValStringVal != "set sub string" {
		t.Errorf("config.StructVal.StructValStringVal should be 'set sub string', got: %q", config.StructVal.StructValStringVal)
	}
	if config.EmbeddedStructValStringVal != "set embedded string" {
		t.Errorf("config.EmbeddedStructValStringVal  should be 'set embedded string', got: %q", config.EmbeddedStructValStringVal)
	}
}

type embeddedStruct struct {
	EmbeddedStructValStringVal string `env:"ENV_EMBEDDED_STRING_VAL" default:"embedded_string"`
}

// embeddedAndSubStruct is an test struct that holds parameters of embedded and sub struct
type embeddedAndSubStruct struct {
	// Sub struct
	StructVal struct {
		StructValStringVal string `env:"ENV_SUB_STRING_VAL" default:"sub_string"`
	}

	// Embedded struct
	embeddedStruct
}
