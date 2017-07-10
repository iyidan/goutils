package config

import (
	"encoding/json"

	"fmt"

	"sync"

	"time"

	"github.com/iyidan/goutils/cmtjson"
	"github.com/iyidan/goutils/mise"
)

const (
	cKeyInt = iota
	cKeyInt64
	cKeyFloat
	cKeyBool
	cKeySliceString
	cKeySliceInt
	cKeySliceFloat
	cKeySliceBool
	cKeyMapString
	cKeyMapInt
	cKeyMapFloat
	cKeyMapBool
	cKeyMapSliceString
	cKeyMapSliceInt
	cKeyMapSliceFloat
	cKeyMapSliceBool
)

// Config store the json.Unmarshal data
type Config struct {
	origin map[string]interface{}
	cached map[int]map[string]interface{}
	cLock  sync.RWMutex
}

func newConf() *Config {
	return &Config{
		origin: make(map[string]interface{}),
		cached: map[int]map[string]interface{}{
			cKeyInt:            make(map[string]interface{}),
			cKeyInt64:          make(map[string]interface{}),
			cKeyFloat:          make(map[string]interface{}),
			cKeyBool:           make(map[string]interface{}),
			cKeySliceString:    make(map[string]interface{}),
			cKeySliceInt:       make(map[string]interface{}),
			cKeySliceFloat:     make(map[string]interface{}),
			cKeySliceBool:      make(map[string]interface{}),
			cKeyMapString:      make(map[string]interface{}),
			cKeyMapInt:         make(map[string]interface{}),
			cKeyMapFloat:       make(map[string]interface{}),
			cKeyMapBool:        make(map[string]interface{}),
			cKeyMapSliceString: make(map[string]interface{}),
			cKeyMapSliceInt:    make(map[string]interface{}),
			cKeyMapSliceFloat:  make(map[string]interface{}),
			cKeyMapSliceBool:   make(map[string]interface{}),
		},
	}
}

func (conf *Config) cacheGet(cKey int, k string) (interface{}, bool) {
	conf.cLock.RLock()
	v, ok := conf.cached[cKey][k]
	conf.cLock.RUnlock()
	return v, ok
}
func (conf *Config) cacheSet(cKey int, k string, v interface{}) {
	conf.cLock.Lock()
	conf.cached[cKey][k] = v
	conf.cLock.Unlock()
}

// ParseFromFile parse config from the given file
func ParseFromFile(filename string) (*Config, error) {
	conf := newConf()
	err := cmtjson.ParseFromFile(filename, &conf.origin)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// ParseFromData parse config with the given data
func ParseFromData(data []byte) (*Config, error) {
	conf := newConf()
	err := cmtjson.ParseFromBytes(data, &conf.origin)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// Get get a config with original value
// if k not exists, return nil
func (conf *Config) Get(k string) interface{} {
	if v, ok := conf.origin[k]; ok {
		return v
	}
	return nil
}

// StrictGet get a config with original value
// if k not exists, panic
func (conf *Config) StrictGet(k string) interface{} {
	v, ok := conf.origin[k]
	if !ok {
		mise.PanicOnError(fmt.Errorf("config not exists:"+k), "config")
	}
	return v
}

// String get a config with k, if k not exists or parse error, panic
func (conf *Config) String(k string) string {
	v := conf.StrictGet(k)
	sv, ok := v.(string)
	if !ok {
		mise.PanicOnError(fmt.Errorf("config not string type: %s => %#v", k, v), "config")
	}
	return sv
}

// GetTime same as conf.String method
func (conf *Config) GetTime(k string) time.Time {
	s := conf.String(k)
	tm, err := mise.StrToLocalTime(s)
	if err != nil {
		err = mise.WrapErrorMsg(err, fmt.Sprintf("config not time-string type: %s => %#v", k, s))
		mise.PanicOnError(err, "config.GetTime")
	}
	return tm
}

// GetDuration same as conf.String method
func (conf *Config) GetDuration(k string) time.Duration {
	s := conf.String(k)
	d, err := time.ParseDuration(s)
	if err != nil {
		err = mise.WrapErrorMsg(err, fmt.Sprintf("config not time-duration-string type: %s => %#v", k, s))
		mise.PanicOnError(err, "config.GetDuration")
	}
	return d
}

// Int same as conf.String method
func (conf *Config) Int(k string) int {
	if cachedv, ok := conf.cacheGet(cKeyInt, k); ok {
		return cachedv.(int)
	}
	v := conf.StrictGet(k)
	iv, err := mise.ParseInt(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	conf.cacheSet(cKeyInt, k, iv)
	return iv
}

// Int64 same as conf.String method
func (conf *Config) Int64(k string) int64 {
	if cachedv, ok := conf.cacheGet(cKeyInt64, k); ok {
		return cachedv.(int64)
	}
	v := conf.StrictGet(k)
	iv, err := mise.ParseInt64(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	conf.cacheSet(cKeyInt64, k, iv)
	return iv
}

// Float same as conf.String method
func (conf *Config) Float(k string) float64 {
	if cachedv, ok := conf.cacheGet(cKeyFloat, k); ok {
		return cachedv.(float64)
	}
	v := conf.StrictGet(k)
	fv, err := mise.ParseFloat(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	conf.cacheSet(cKeyFloat, k, fv)
	return fv
}

// Bool same as conf.String method
func (conf *Config) Bool(k string) bool {
	if cachedv, ok := conf.cacheGet(cKeyBool, k); ok {
		return cachedv.(bool)
	}
	v := conf.StrictGet(k)
	bv, err := mise.ParseBool(v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
	conf.cacheSet(cKeyBool, k, bv)
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
func (conf *Config) Slice(k string) []interface{} {
	return sliceVal(k, conf.StrictGet(k))
}

// SliceString same as conf.String method
func (conf *Config) SliceString(k string) []string {
	if cachedv, ok := conf.cacheGet(cKeySliceString, k); ok {
		return cachedv.([]string)
	}
	v := sliceStringVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeySliceString, k, v)
	return v
}

// SliceInt same as conf.String method
func (conf *Config) SliceInt(k string) []int {
	if cachedv, ok := conf.cacheGet(cKeySliceInt, k); ok {
		return cachedv.([]int)
	}
	v := sliceIntVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeySliceInt, k, v)
	return v
}

// SliceFloat same as conf.String method
func (conf *Config) SliceFloat(k string) []float64 {
	if cachedv, ok := conf.cacheGet(cKeySliceFloat, k); ok {
		return cachedv.([]float64)
	}
	v := sliceFloatVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeySliceFloat, k, v)
	return v
}

// SliceBool same as conf.String method
func (conf *Config) SliceBool(k string) []bool {
	if cachedv, ok := conf.cacheGet(cKeySliceBool, k); ok {
		return cachedv.([]bool)
	}
	v := sliceBoolVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeySliceBool, k, v)
	return v
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
func (conf *Config) MapString(k string) map[string]interface{} {
	return mapStringVal(k, conf.StrictGet(k))
}

// MapStringString same as conf.String method
func (conf *Config) MapStringString(k string) map[string]string {
	if cachedv, ok := conf.cacheGet(cKeyMapString, k); ok {
		return cachedv.(map[string]string)
	}
	v := mapStringStringVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapString, k, v)
	return v
}

// MapStringInt same as conf.String method
func (conf *Config) MapStringInt(k string) map[string]int {
	if cachedv, ok := conf.cacheGet(cKeyMapInt, k); ok {
		return cachedv.(map[string]int)
	}
	v := mapStringIntVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapInt, k, v)
	return v
}

// MapStringFloat same as conf.String method
func (conf *Config) MapStringFloat(k string) map[string]float64 {
	if cachedv, ok := conf.cacheGet(cKeyMapFloat, k); ok {
		return cachedv.(map[string]float64)
	}
	v := mapStringFloatVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapFloat, k, v)
	return v
}

// MapStringBool same as conf.String method
func (conf *Config) MapStringBool(k string) map[string]bool {
	if cachedv, ok := conf.cacheGet(cKeyMapBool, k); ok {
		return cachedv.(map[string]bool)
	}
	v := mapStringBoolVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapBool, k, v)
	return v
}

// MapStringSliceString same as conf.String method
func (conf *Config) MapStringSliceString(k string) map[string][]string {
	if cachedv, ok := conf.cacheGet(cKeyMapSliceString, k); ok {
		return cachedv.(map[string][]string)
	}
	v := mapStringSliceStringVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapSliceString, k, v)
	return v
}

// MapStringSliceInt same as conf.String method
func (conf *Config) MapStringSliceInt(k string) map[string][]int {
	if cachedv, ok := conf.cacheGet(cKeyMapSliceInt, k); ok {
		return cachedv.(map[string][]int)
	}
	v := mapStringSliceIntVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapSliceInt, k, v)
	return v
}

// MapStringSliceFloat same as conf.String method
func (conf *Config) MapStringSliceFloat(k string) map[string][]float64 {
	if cachedv, ok := conf.cacheGet(cKeyMapSliceFloat, k); ok {
		return cachedv.(map[string][]float64)
	}
	v := mapStringSliceFloatVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapSliceFloat, k, v)
	return v
}

// MapStringSliceBool same as conf.String method
func (conf *Config) MapStringSliceBool(k string) map[string][]bool {
	if cachedv, ok := conf.cacheGet(cKeyMapSliceBool, k); ok {
		return cachedv.(map[string][]bool)
	}
	v := mapStringSliceBoolVal(k, conf.StrictGet(k))
	conf.cacheSet(cKeyMapSliceBool, k, v)
	return v
}

// Unmarshal config k into v
func (conf *Config) Unmarshal(k string, v interface{}) error {
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
func (conf *Config) MustUnmarshal(k string, v interface{}) {
	err := conf.Unmarshal(k, v)
	if err != nil {
		mise.PanicOnError(err, "config")
	}
}
