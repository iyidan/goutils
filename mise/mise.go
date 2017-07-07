package mise

import (
	"fmt"
	"reflect"
	"strconv"
)

// GetValueKind get the given value's kind
func GetValueKind(val interface{}) (reflect.Value, reflect.Kind) {
	v := reflect.ValueOf(val)
	kd := v.Kind()
	if kd == reflect.Ptr {
		return GetValueKind(v.Elem().Interface())
	}
	return v, kd
}

// ParseFloat parse any float-like-value into float64 type
func ParseFloat(val interface{}) (float64, error) {
	v, kd := GetValueKind(val)
	switch kd {
	case reflect.Float64, reflect.Float32:
		return v.Float(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), nil
	case reflect.String:
		number, err := strconv.ParseFloat(v.String(), 64)
		if err != nil {
			return 0, err
		}
		return number, nil
	default:
		return 0, fmt.Errorf("%#v parse float-like-value failed", val)
	}
}

// ParseInt64 parse any int-like-value into int64 type
func ParseInt64(val interface{}) (int64, error) {
	v, kd := GetValueKind(val)
	switch kd {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(v.Uint()), nil
	case reflect.String:
		number, err := strconv.ParseInt(v.String(), 10, 64)
		if err != nil {
			return 0, err
		}
		return number, nil
	default:
		return 0, fmt.Errorf("%#v parse int-like-value failed", val)
	}
}

// ParseInt similar as ParseInt64() function
func ParseInt(val interface{}) (int, error) {
	v, err := ParseInt64(val)
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

// ParseBool parse any bool-like-value into bool type
func ParseBool(val interface{}) (bool, error) {
	v, kd := GetValueKind(val)
	switch kd {
	case reflect.Bool:
		return v.Bool(), nil
	case reflect.String:
		switch v.String() {
		case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
			return true, nil
		case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
			return false, nil
		}
	case reflect.Float64, reflect.Float32:
		if v.Float() == 1 {
			return true, nil
		} else if v.Float() == 0 {
			return false, nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if v.Int() == 1 {
			return true, nil
		} else if v.Int() == 0 {
			return false, nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v.Uint() == 1 {
			return true, nil
		} else if v.Uint() == 0 {
			return false, nil
		}
	}
	return false, fmt.Errorf("%#v parse bool failed", val)
}
