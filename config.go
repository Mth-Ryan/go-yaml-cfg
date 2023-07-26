package goyamlcfg

import (
	"fmt"
	"os"
	"reflect"
	"sync"

	"gopkg.in/yaml.v3"
)

type Instance interface{}

type config struct {
	mu   sync.Mutex
	ty   reflect.Type
	data interface{}
}

var singleton config

// InitializeConfigSingleton will create a instance of the configuration
// to be fetched with GetConfigFromSingleton or MustGetConfigFromSingleton.
// Use this when you are not using a DI system.
func InitializeConfigSingleton[T Instance](path string) error {
	if singleton.data == nil {
		singleton.mu.Lock()
		defer singleton.mu.Unlock()

		target, err := GetConfig[T](path)
		if err != nil {
			return err
		}

		singleton.ty = reflect.TypeOf(target)
		singleton.data = interface{}(target)
	}

	return nil
}

// MustInitializeConfigSingleton has the same behaviour of InitializeConfigSingleton,
// but will panic if an error occurs.
func MustInitializeConfigSingleton[T Instance](path string) {
	err := InitializeConfigSingleton[T](path)
	if err != nil {
		panic(err)
	}
}

// GetConfigFromSingleton will only return a configuration with type T if:
// 1) The singleton has been already initialized.
// 2) The type of the generic is the same of the singleton initialization.
func GetConfigFromSingleton[T Instance]() (T, error) {
	if singleton.data != nil {
		if singleton.ty != reflect.TypeOf(*new(T)) {
			return *new(T), fmt.Errorf("mismatch generic type: the given type is diferent from the singleton data")
		}
		return (singleton.data).(T), nil
	} else {
		return *new(T), fmt.Errorf("trying to access an unitialized config data")
	}
}

// MustGetConfigFromSingleton has the same behaviour of GetConfigFromSingleton
// but will panic if an error occurs.
func MustGetConfigFromSingleton[T Instance]() T {
	conf, err := GetConfigFromSingleton[T]()
	if err != nil {
		panic(err)
	}
	return conf
}

// GetConfig will always return a new config instance with the type T if the
// read file and the parsing stages was successful.
func GetConfig[T Instance](path string) (T, error) {
	target := new(T)
	err := loadAndParseConfig(path, target)
	if err != nil {
		return *target, err
	}
	replaceEnvOnTarget(target)
	return *target, nil
}

// MustGetConfig has the same behaviour of GetConfig but panics
// when an error happens.
func MustGetConfig[T Instance](path string) T {
	target, err := GetConfig[T](path)
	if err != nil {
		panic(err)
	}
	return target
}

func loadAndParseConfig[T Instance](path string, target *T) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, target)
}
