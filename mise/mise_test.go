package mise

import (
	"fmt"
	"testing"
)

func TestParseBool(t *testing.T) {

	trueCases := []interface{}{
		true,
		"1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On",
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint32(1),
		1.0,
		float32(1.0),
	}

	falseCases := []interface{}{
		false,
		"0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off",
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint32(0),
		0.0,
		float32(0.0),
	}

	for _, val := range trueCases {
		b, err := ParseBool(val)
		if err != nil {
			t.Fatal(val, err)
		}
		if b != true {
			t.Fatal(val, "test failed")
		}
		t.Logf("test-true-case: %T(%#v) ok", val, val)
	}

	for _, val := range falseCases {
		b, err := ParseBool(val)
		if err != nil {
			t.Fatal(val, err)
		}
		if b != false {
			t.Fatal(val, "test failed")
		}
		t.Logf("test-false-case: %T(%#v) ok", val, val)
	}

	_, err := ParseBool("fawef")
	if err == nil {
		t.Fatal(`ParseBool("fawef") test fail`)
	}
}

func TestParseFloat(t *testing.T) {
	v := float64(33.0)
	numberCases := []interface{}{
		fmt.Sprintf("%v", v),
		int(v),
		int8(v),
		int16(v),
		int32(v),
		int64(v),
		uint(v),
		uint32(v),
		v,
		float32(v),
	}

	for _, val := range numberCases {
		number, err := ParseFloat(val)
		if err != nil {
			t.Fatal(val, err)
		}
		if number != v {
			t.Fatal(val, "test failed")
		}
		t.Logf("test-number-case: %T(%#v) ok", val, val)
	}

	_, err := ParseFloat("fawef")
	if err == nil {
		t.Fatal(`ParseFloat("fawef") test fail`)
	}
}

func TestParseInt64(t *testing.T) {
	v := int64(33)
	numberCases := []interface{}{
		fmt.Sprintf("%v", v),
		int(v),
		int8(v),
		int16(v),
		int32(v),
		int64(v),
		uint(v),
		uint32(v),
		v,
	}

	for _, val := range numberCases {
		number, err := ParseInt64(val)
		if err != nil {
			t.Fatal(val, err)
		}
		if number != v {
			t.Fatal(val, "test failed")
		}
		t.Logf("test-number-case: %T(%#v) ok", val, val)
	}

	_, err := ParseInt64("fawef")
	if err == nil {
		t.Fatal(`ParseInt64("fawef") test fail`)
	}
}
