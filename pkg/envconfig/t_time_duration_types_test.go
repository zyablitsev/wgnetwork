package envconfig

import (
	"os"
	"testing"
	"time"
)

func TestTimeDurationTypesDefaults(t *testing.T) {
	os.Clearenv()
	config := new(timeDurationTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if config.NanosecondVal != 1*time.Nanosecond {
		t.Errorf(
			"config.NanosecondVal should be %d, got: %d",
			1*time.Nanosecond, config.NanosecondVal,
		)
	}
	if config.MicrosecondVal != 1*time.Microsecond {
		t.Errorf(
			"config.MicrosecondVal should be %d, got: %d",
			1*time.Microsecond, config.MicrosecondVal,
		)
	}
	if config.MillisecondVal != 1*time.Millisecond {
		t.Errorf(
			"config.MillisecondVal should be %d, got: %d",
			1*time.Millisecond, config.MillisecondVal,
		)
	}
	if config.SecondVal != 1*time.Second {
		t.Errorf(
			"config.SecondVal should be %d, got: %d",
			1*time.Second, config.SecondVal,
		)
	}
	if config.MinuteVal != 1*time.Minute {
		t.Errorf(
			"config.MinuteVal should be %d, got: %d",
			1*time.Minute, config.MinuteVal,
		)
	}
	if config.HourVal != 1*time.Hour {
		t.Errorf(
			"config.HourVal should be %d, got: %d",
			1*time.Hour, config.HourVal,
		)
	}
}

func TestTimeDurationTypes(t *testing.T) {
	os.Clearenv()
	envs := envValues{
		"ENV_NANOSECOND":  envValue{checkValue: "2ns"},
		"ENV_MICROSECOND": envValue{checkValue: "2us"},
		"ENV_MILLISECOND": envValue{checkValue: "2ms"},
		"ENV_SECOND":      envValue{checkValue: "2s"},
		"ENV_MINUTE":      envValue{checkValue: "2m"},
		"ENV_HOUR":        envValue{checkValue: "2h"},
	}
	if err := envs.set(); err != nil {
		t.Fatalf("envs.set() fail, raw: %v", err)
	}
	config := new(timeDurationTypes)
	if err := ReadEnv(config); err != nil {
		t.Fatalf("ReadEnv fail, raw: %v", err)
	}

	if config.NanosecondVal != 2*time.Nanosecond {
		t.Errorf(
			"config.NanosecondVal should be %d, got: %d",
			2*time.Nanosecond, config.NanosecondVal,
		)
	}
	if config.MicrosecondVal != 2*time.Microsecond {
		t.Errorf(
			"config.MicrosecondVal should be %d, got: %d",
			2*time.Microsecond, config.MicrosecondVal,
		)
	}
	if config.MillisecondVal != 2*time.Millisecond {
		t.Errorf(
			"config.MillisecondVal should be %d, got: %d",
			2*time.Millisecond, config.MillisecondVal,
		)
	}
	if config.SecondVal != 2*time.Second {
		t.Errorf(
			"config.SecondVal should be %d, got: %d",
			2*time.Second, config.SecondVal,
		)
	}
	if config.MinuteVal != 2*time.Minute {
		t.Errorf(
			"config.MinuteVal should be %d, got: %d",
			2*time.Minute, config.MinuteVal,
		)
	}
	if config.HourVal != 2*time.Hour {
		t.Errorf(
			"config.HourVal should be %d, got: %d",
			2*time.Hour, config.HourVal,
		)
	}
}

// timeDurationTypes is an test struct
// that holds parameters of different time.Duration unit types
type timeDurationTypes struct {
	// time.Duration types
	NanosecondVal  time.Duration `env:"ENV_NANOSECOND" default:"1ns"`
	MicrosecondVal time.Duration `env:"ENV_MICROSECOND" default:"1us"`
	MillisecondVal time.Duration `env:"ENV_MILLISECOND" default:"1ms"`
	SecondVal      time.Duration `env:"ENV_SECOND" default:"1s"`
	MinuteVal      time.Duration `env:"ENV_MINUTE" default:"1m"`
	HourVal        time.Duration `env:"ENV_HOUR" default:"1h"`
}
