package config

import (
	"io/ioutil"
	"reflect"
	"testing"
	"time"

	"github.com/iyidan/goutils/mise"
)

var testConfig = `

# 井号注释

{
# 井号注释
// 单行注释
/* 区块注释 */
/**
 * 区块注释 
 */
"StringKey":"stringVal",
"StringTimeKey": "2017-07-15 09:00:00",
"StringDurationKey": "1m30s50ms",
"IntKey": 32,
"IntKey2": "32",
"Int64Key": 123456789, # 井号注释
"FloatKey": 321.234,
"BoolKey": true, // 单行注释
"BoolKey2": "On",
"BoolKey3":false,
"BoolKey4":"false",

"SliceStringKey": ["a","b", "c"], /* 区块注释 */
"SliceIntKey": ["-1",2,3],
"SliceFloatKey": [-1,0.3,"123.55"],
"SliceBoolKey": ["false","true","on","off", false, true, 0, 1],

"MapStringKey": {"a":"a", "b":"b"},
"MapStringIntKey":{"a":3, "b":2, "c":"3", "d":-3},
"MapStringFloatKey":{"a":3.1, "b":2, "c":"3.1", "d":"-3.1"},
"MapStringBoolKey":{"a":true, "b":"True", "c":"false"},

"MapStringSliceKey": {"a":["a","b","c"]},
"MapStringSliceFloatKey":{"a":["32.1",32,1.1,"-3.0"]},
"MapStringSliceIntKey":{"a":["32",-32,1,"-3"]},
"MapStringSliceBoolKey":{"a":["false","true","on","off", false, true, 0, 1]}
}`

var testConfigCorrect = map[string]interface{}{
	"StringKey":         "stringVal",
	"StringTimeKey":     "2017-07-15 09:00:00",
	"StringDurationKey": "1m30s50ms",
	"IntKey":            int(32),
	"IntKey2":           int(32),
	"Int64Key":          int64(123456789),
	"FloatKey":          float64(321.234),
	"BoolKey":           true,
	"BoolKey2":          true,
	"BoolKey3":          false,
	"BoolKey4":          false,

	"SliceStringKey": []string{"a", "b", "c"},
	"SliceIntKey":    []int{-1, 2, 3},
	"SliceFloatKey":  []float64{-1, 0.3, 123.55},
	"SliceBoolKey":   []bool{false, true, true, false, false, true, false, true},

	"MapStringKey":      map[string]string{"a": "a", "b": "b"},
	"MapStringIntKey":   map[string]int{"a": 3, "b": 2, "c": 3, "d": -3},
	"MapStringFloatKey": map[string]float64{"a": 3.1, "b": 2, "c": 3.1, "d": -3.1},
	"MapStringBoolKey":  map[string]bool{"a": true, "b": true, "c": false},

	"MapStringSliceKey":      map[string][]string{"a": []string{"a", "b", "c"}},
	"MapStringSliceFloatKey": map[string][]float64{"a": []float64{32.1, 32, 1.1, -3.0}},
	"MapStringSliceIntKey":   map[string][]int{"a": []int{32, -32, 1, -3}},
	"MapStringSliceBoolKey":  map[string][]bool{"a": []bool{false, true, true, false, false, true, false, true}},
}

func init() {
	tmval, err := mise.StrToLocalTime(testConfigCorrect["StringTimeKey"].(string))
	if err != nil {
		panic(err)
	}
	testConfigCorrect["StringTimeKey"] = tmval

	duval, err := time.ParseDuration(testConfigCorrect["StringDurationKey"].(string))
	if err != nil {
		panic(err)
	}
	testConfigCorrect["StringDurationKey"] = duval

}

func TestConfigFuncs(t *testing.T) {
	tmpfile, err := getTempfileWithJSON([]byte(testConfig))
	if err != nil {
		t.Fatal(err)
	}
	conf, err := ParseFromFile(tmpfile)
	if err != nil {
		t.Fatal(err)
	}

	testFunc := func(msg string) {
		t.Logf("=== %s === \n", msg)
		for k, v := range testConfigCorrect {
			var parsed interface{}
			switch v.(type) {
			case string:
				parsed = conf.String(k)
			case int:
				parsed = conf.Int(k)
			case int64:
				parsed = conf.Int64(k)
			case float64:
				parsed = conf.Float(k)
			case bool:
				parsed = conf.Bool(k)

			case []string:
				parsed = conf.SliceString(k)
			case []int:
				parsed = conf.SliceInt(k)
			case []float64:
				parsed = conf.SliceFloat(k)
			case []bool:
				parsed = conf.SliceBool(k)

			case map[string]string:
				parsed = conf.MapStringString(k)
			case map[string]int:
				parsed = conf.MapStringInt(k)
			case map[string]float64:
				parsed = conf.MapStringFloat(k)
			case map[string]bool:
				parsed = conf.MapStringBool(k)

			case map[string][]string:
				parsed = conf.MapStringSliceString(k)
			case map[string][]int:
				parsed = conf.MapStringSliceInt(k)
			case map[string][]float64:
				parsed = conf.MapStringSliceFloat(k)
			case map[string][]bool:
				parsed = conf.MapStringSliceBool(k)

			case time.Time:
				parsed = conf.GetTime(k)
				t.Logf("GetTime: %s => time: %s\n", k, parsed)
			case time.Duration:
				parsed = conf.GetDuration(k)
				t.Logf("GetDuration: %s => time-duration: %s\n", k, parsed)
			default:
				t.Fatalf("unsupported config type: %s => %T(%#v)\n", k, v, v)
			}
			if !reflect.DeepEqual(parsed, v) {
				t.Fatalf("config not Equal! parsed: %T(%#v), right: %T(%#v)\n", parsed, parsed, v, v)
			}
			t.Logf("parsed: %s => %T(%#v), right: %T(%#v) --ok\n", k, parsed, parsed, v, v)
		}
	}

	testFunc("1")
	testFunc("2")
	testFunc("3")
}

func BenchmarkConfigGet(b *testing.B) {
	b.StopTimer()
	conf, err := ParseFromData([]byte(testConfig))
	if err != nil {
		b.Fatal(err)
	}
	key := "MapStringSliceBoolKey"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = conf.MapStringSliceBool(key)
	}
}

func BenchmarkMapGet(b *testing.B) {
	key := "MapStringSliceBoolKey"
	for i := 0; i < b.N; i++ {
		_ = testConfigCorrect[key].(map[string][]bool)
	}
}

func getTempfileWithJSON(data []byte) (string, error) {
	tmpfile, err := ioutil.TempFile("", "getTempfileWithJSON")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(data)
	if err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}
