package config

import (
	"bytes"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BxdConfig struct {
	folder string
	keyDelim string

	envMaps map[string]string
	confMaps map[string]interface{}
	confRows map[string][]byte
}

func NewBxdConfig(params ...interface{}) (interface{}, error) {
	if len(params) > 2 {
		return nil, errors.New("NewBxdConfig params error")
	}

	folder := params[0].(string)
	var envMaps map[string]string
	if len(params) >= 2 {
		envMaps = params[1].(map[string]string)
	}

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return nil, errors.New("folder " + folder + " not exist: " + err.Error())
	}

	bxdConf := &BxdConfig{
		folder:folder,
		envMaps:envMaps,
		confMaps:map[string]interface{}{},
		confRows:map[string][]byte{},
		keyDelim:".",
	}

	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, file := range files {
		s := strings.Split(file.Name(), ".")
		if len(s) == 2 && (s[1] == "yaml" || s[1] == "yml") {
			name := s[0]

			bf, err := ioutil.ReadFile(filepath.Join(folder, file.Name()))
			if err != nil {
				continue
			}
			bxdConf.confRows[name] = bf
			bf = replace(bf, envMaps)

			c := map[string]interface{}{}
			if err := yaml.Unmarshal(bf, &c); err != nil {
				continue
			}

			bxdConf.confMaps[name] = c
		}
	}

	return bxdConf, nil
}

func replace(content []byte, maps map[string]string) []byte {
	if maps == nil {
		return content
	}

	for key, val := range maps {
		reKey := "env(" + key + ")"
		content = bytes.ReplaceAll(content, []byte(reKey), []byte(val))
	}

	return content
}

func searchMap(source map[string]interface{}, path []string) interface{} {
	if len(path) == 0 {
		return source
	}

	next, ok := source[path[0]]
	if ok {
		if len(path) == 1 {
			return next
		}
		switch next.(type) {
		case map[interface{}]interface{}:
			return searchMap(cast.ToStringMap(next), path[1:])
		case map[string]interface{}:
			return searchMap(next.(map[string]interface{}), path[:1])
		default:
			return nil
		}
	}

	return nil
}

func (conf *BxdConfig) find(key string) interface{} {
	return searchMap(conf.confMaps, strings.Split(key, conf.keyDelim))
}

func (conf *BxdConfig) IsExist(key string) bool {
	return conf.find(key) != nil
}

func (conf *BxdConfig) Get(key string) interface{} {
	return conf.find(key)
}

func (conf *BxdConfig) GetBool(key string) bool {
	return cast.ToBool(conf.find(key))
}

func (conf *BxdConfig) GetInt(key string) int {
	return cast.ToInt(conf.find(key))
}

// GetFloat64 get float64
func (conf *BxdConfig) GetFloat64(key string) float64 {
	return cast.ToFloat64(conf.find(key))
}

// GetTime get time type
func (conf *BxdConfig) GetTime(key string) time.Time {
	return cast.ToTime(conf.find(key))
}

// GetString get string typen
func (conf *BxdConfig) GetString(key string) string {
	return cast.ToString(conf.find(key))
}

// GetIntSlice get int slice type
func (conf *BxdConfig) GetIntSlice(key string) []int {
	return cast.ToIntSlice(conf.find(key))
}

// GetStringSlice get string slice type
func (conf *BxdConfig) GetStringSlice(key string) []string {
	return cast.ToStringSlice(conf.find(key))
}

// GetStringMap get map which key is string, value is interface
func (conf *BxdConfig) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(conf.find(key))
}

// GetStringMapString get map which key is string, value is string
func (conf *BxdConfig) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(conf.find(key))
}

// GetStringMapStringSlice get map which key is string, value is string slice
func (conf *BxdConfig) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(conf.find(key))
}

func (conf *BxdConfig) Load(key string, val interface{}) error {
	return mapstructure.Decode(conf.find(key), val)
}