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
	v := ToFixed(33.1, 1)
	vtruncated := ToFixed(33.0, 1)
	numberCases := [][2]interface{}{
		[2]interface{}{fmt.Sprintf("%v", v), v},
		[2]interface{}{int(v), vtruncated},
		[2]interface{}{int(-v), -vtruncated},
		[2]interface{}{int8(v), vtruncated},
		[2]interface{}{int32(v), vtruncated},
		[2]interface{}{int64(v), vtruncated},
		[2]interface{}{uint(v), vtruncated},
		[2]interface{}{v, v},
		[2]interface{}{float32(v), v},
	}

	for _, val := range numberCases {
		number, err := ParseFloat(val[0])
		if err != nil {
			t.Fatal(val, err)
		}
		if ToFixed(number, 1) != val[1] {
			t.Fatalf("%T(%#v) %#v test failed", number, number, val)
		}
		t.Logf("test-number-case: %#v ok", val)
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
		float64(v),
		float32(v),
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
	t.Log(err)
	_, err = ParseInt64(33.2)
	if err == nil {
		t.Fatal(`ParseInt64(33.2) test fail`)
	}
	t.Log(err)
}

func TestRound(t *testing.T) {
	if Round(33.1) != 33 {
		t.Fatal(`Round(33.1) != 33`)
	}
	if Round(33.5) != 34 {
		t.Fatal(`Round(33.5) != 34`)
	}
	if Round(33.0) != 33 {
		t.Fatal(`Round(33.0) != 33`)
	}
}

func TestToFixed(t *testing.T) {
	if ToFixed(1.0/3.0, 2) != 0.33 {
		t.Fatal(`ToFixed(1.0/3.0) != 0.3`)
	}

	if ToFixed(2.0/3.0, 2) != 0.67 {
		t.Fatal(`ToFixed(2.0/3.0, 2) != 0.67`)
	}
}
