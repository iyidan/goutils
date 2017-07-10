package config

import (
	"encoding/json"

	"fmt"

	"github.com/iyidan/goutils/cmtjson"
	"github.com/iyidan/goutils/mise"
)

// Config store the json.Unmarshal data
type Config map[string]interface{}

// ParseFromFile parse config from the given file
func ParseFromFile(filename string) (Config, error) {
	conf := make(Config)
	err := cmtjson.ParseFromFile(filename, &conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}

// ParseFromData parse config with the given data
func ParseFromData(data []byte) (Config, error) {
	conf := make(Config)
	err := cmtjson.ParseFromBytes(data, &conf)
	if err != nil {
		return nil, err
	}
	return conf, err
}

// Get get a config with original value
// if k not exists, return nil
func (conf Config) Get(k string) interface{} {
	if v, ok := conf[k]; ok {
		return v
	}
	return nil
}

// StrictGet get a config with original value
// if k not exists, panic
func (conf Config) StrictGet(k string) interface{} {
	v, ok := conf[k]
	if !ok {
		mise.PanicOnError(fmt.Errorf("config not exists:"+k), "config")
	}
	return v
}

// String get a config with k, if k not exists or parse error, panic
func (conf Config) String(k string) string {
	v := conf.StrictGet(k)
	sv, ok := v.(string)
	if !ok {
		mise.PanicOnError(fmt.Errorf("config not string type: %s => %#v", k, v), "config")
	}
	return sv
}

// Int same as conf.String method
func (conf Config) Int(k string) int {
	v := conf.StrictGet(k)
	iv, err := mise.ParseInt(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	return iv
}

// Int64 same as conf.String method
func (conf Config) Int64(k string) int64 {
	v := conf.StrictGet(k)
	iv, err := mise.ParseInt64(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	return iv
}

// Float same as conf.String method
func (conf Config) Float(k string) float64 {
	v := conf.StrictGet(k)
	fv, err := mise.ParseFloat(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	return fv
}

// Bool same as conf.String method
func (conf Config) Bool(k string) bool {
	v := conf.StrictGet(k)
	bv, err := mise.ParseBool(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	return bv
}

func sliceVal(id string, val interface{}) []interface{} {
	tmp, ok := val.([]interface{})
	if !ok {
		mise.PanicOnError(fmt.Errorf("config.sliceVal not []interface type: %s => %#v", id, val), "config")
	}
	return tmp
}
func sliceStringVal(id string, val interface{}) []string {
	tmp := sliceVal(id, val)
	ssv := make([]string, len(tmp))
	for i := range tmp {
		s, ok := tmp[i].(string)
		if !ok {
			mise.PanicOnError(fmt.Errorf("config.sliceStringVal not []string type: %s => %#v", id, tmp), "config")
		}
		ssv[i] = s
	}
	return ssv
}

func sliceIntVal(id string, val interface{}) []int {
	tmp := sliceVal(id, val)
	ssv := make([]int, len(tmp))
	for i := range tmp {
		iv, err := mise.ParseInt(tmp[i])
		if err != nil {
			mise.PanicOnError(err, "config")
		}
		ssv[i] = iv
	}
	return ssv
}

func sliceFloatVal(id string, val interface{}) []float64 {
	tmp := sliceVal(id, val)
	ssv := make([]float64, len(tmp))
	for i := range tmp {
		iv, err := mise.ParseFloat(tmp[i])
		if err != nil {
			mise.PanicOnError(err, "config")
		}
		ssv[i] = iv
	}
	return ssv
}

func sliceBoolVal(id string, val interface{}) []bool {
	tmp := sliceVal(id, val)
	ssv := make([]bool, len(tmp))
	for i := range tmp {
		iv, err := mise.ParseBool(tmp[i])
		if err != nil {
			mise.PanicOnError(err, "config")
		}
		ssv[i] = iv
	}
	return ssv
}

// Slice same as conf.String method
func (conf Config) Slice(k string) []interface{} {
	return sliceVal(k, conf.StrictGet(k))
}

// SliceString same as conf.String method
func (conf Config) SliceString(k string) []string {
	return sliceStringVal(k, conf.StrictGet(k))
}

// SliceInt same as conf.String method
func (conf Config) SliceInt(k string) []int {
	return sliceIntVal(k, conf.StrictGet(k))
}

// SliceFloat same as conf.String method
func (conf Config) SliceFloat(k string) []float64 {
	return sliceFloatVal(k, conf.StrictGet(k))
}

// SliceBool same as conf.String method
func (conf Config) SliceBool(k string) []bool {
	return sliceBoolVal(k, conf.StrictGet(k))
}

func mapStringVal(id string, val interface{}) map[string]interface{} {
	tmp, ok := val.(map[string]interface{})
	if !ok {
		mise.PanicOnError(fmt.Errorf("config not map[string]interface{} type: %s => %#v", id, val), "config")
	}
	return tmp
}

func mapStringStringVal(id string, val interface{}) map[string]string {
	tmp := mapStringVal(id, val)
	r := make(map[string]string)
	for tmpkey, tmpval := range tmp {
		sv, ok := tmpval.(string)
		if !ok {
			mise.PanicOnError(fmt.Errorf("config not map[string]string type: %s => %#v", id, tmp), "config")
		}
		r[tmpkey] = sv
	}
	return r
}

func mapStringIntVal(id string, val interface{}) map[string]int {
	tmp := mapStringVal(id, val)
	r := make(map[string]int)
	for tmpkey, tmpval := range tmp {
		iv, err := mise.ParseInt(tmpval)
		if err != nil {
			mise.PanicOnError(err, "config")
		}
		r[tmpkey] = iv
	}
	return r
}

func mapStringFloatVal(id string, val interface{}) map[string]float64 {
	tmp := mapStringVal(id, val)
	r := make(map[string]float64)
	for tmpkey, tmpval := range tmp {
		iv, err := mise.ParseFloat(tmpval)
		if err != nil {
			mise.PanicOnError(err, "config")
		}
		r[tmpkey] = iv
	}
	return r
}

func mapStringBoolVal(id string, val interface{}) map[string]bool {
	tmp := mapStringVal(id, val)
	r := make(map[string]bool)
	for tmpkey, tmpval := range tmp {
		iv, err := mise.ParseBool(tmpval)
		if err != nil {
			mise.PanicOnError(err, "config")
		}
		r[tmpkey] = iv
	}
	return r
}

func mapStringSliceStringVal(id string, val interface{}) map[string][]string {
	tmp := mapStringVal(id, val)
	r := make(map[string][]string)
	for tmpkey, tmpval := range tmp {
		r[tmpkey] = sliceStringVal(id+"."+tmpkey, tmpval)
	}
	return r
}

func mapStringSliceIntVal(id string, val interface{}) map[string][]int {
	tmp := mapStringVal(id, val)
	r := make(map[string][]int)
	for tmpkey, tmpval := range tmp {
		r[tmpkey] = sliceIntVal(id+"."+tmpkey, tmpval)
	}
	return r
}

func mapStringSliceFloatVal(id string, val interface{}) map[string][]float64 {
	tmp := mapStringVal(id, val)
	r := make(map[string][]float64)
	for tmpkey, tmpval := range tmp {
		r[tmpkey] = sliceFloatVal(id+"."+tmpkey, tmpval)
	}
	return r
}

func mapStringSliceBoolVal(id string, val interface{}) map[string][]bool {
	tmp := mapStringVal(id, val)
	r := make(map[string][]bool)
	for tmpkey, tmpval := range tmp {
		r[tmpkey] = sliceBoolVal(id+"."+tmpkey, tmpval)
	}
	return r
}

// MapString same as conf.String method
func (conf Config) MapString(k string) map[string]interface{} {
	return mapStringVal(k, conf.StrictGet(k))
}

// MapStringString same as conf.String method
func (conf Config) MapStringString(k string) map[string]string {
	return mapStringStringVal(k, conf.StrictGet(k))
}

// MapStringInt same as conf.String method
func (conf Config) MapStringInt(k string) map[string]int {
	return mapStringIntVal(k, conf.StrictGet(k))
}

// MapStringFloat same as conf.String method
func (conf Config) MapStringFloat(k string) map[string]float64 {
	return mapStringFloatVal(k, conf.StrictGet(k))
}

// MapStringBool same as conf.String method
func (conf Config) MapStringBool(k string) map[string]bool {
	return mapStringBoolVal(k, conf.StrictGet(k))
}

// MapStringSliceString same as conf.String method
func (conf Config) MapStringSliceString(k string) map[string][]string {
	return mapStringSliceStringVal(k, conf.StrictGet(k))
}

// MapStringSliceInt same as conf.String method
func (conf Config) MapStringSliceInt(k string) map[string][]int {
	return mapStringSliceIntVal(k, conf.StrictGet(k))
}

// MapStringSliceFloat same as conf.String method
func (conf Config) MapStringSliceFloat(k string) map[string][]float64 {
	return mapStringSliceFloatVal(k, conf.StrictGet(k))
}

// MapStringSliceBool same as conf.String method
func (conf Config) MapStringSliceBool(k string) map[string][]bool {
	return mapStringSliceBoolVal(k, conf.StrictGet(k))
}

// Unmarshal config k into v
func (conf Config) Unmarshal(k string, v interface{}) error {
	tmp := conf.MapString(k)
	data, err := json.Marshal(tmp)
	if err != nil {
		return mise.WrapErrorMsg(err, fmt.Sprintf("config.Unmarshal(%s) => %#v", k, tmp))
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return mise.WrapErrorMsg(err, fmt.Sprintf("config.Unmarshal(%s) => %#v", k, tmp))
	}
	return nil
}

// MustUnmarshal config k into v, if error, panic
func (conf Config) MustUnmarshal(k string, v interface{}) {
	err := conf.Unmarshal(k, v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
}
