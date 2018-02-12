package config

import "fmt"

type Configer interface {
	String(key string) string
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
	Set(key, value string) error
}

// Config is the adapter interface for parsing config file
type Config interface {
	Parse(filename string) (Configer, error)
	//ParseData(data []byte) (Configer, error)
}

var adapters = make(map[string]Config)

func Register(name string, adapter Config) {
	if adapter == nil {
		panic("config: Register adapter is nil")
	}

	if _, ok := adapters[name]; ok {
		panic("config: Register called twice for adapter " + name)
	}

	adapters[name] = adapter
}

// NewConfig adapterName is ini/xml/json
// filename is the config file path
func NewConfig(adapterName, filename string) (Configer, error) {
	adapter, ok := adapters[adapterName]
	if !ok {
		return nil, fmt.Errorf("config: Unkown adaptername %q", adapterName)
	}
	return adapter.Parse(filename)
}

func ParseBool(val interface{}) (bool, error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%d", v)
			if strV == "1" {
				return true, nil
			} else {
				return false, nil
			}
		case float64:
			if v == 1.0 {
				return true, nil
			} else {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}

	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}
