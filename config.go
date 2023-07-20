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

func InitializeConfigSingleton[T Instance](path string) error {
	singleton.mu.Lock()
	defer singleton.mu.Unlock()

	target := new(T)
	loadAndParseConfig(path, target)
	replaceEnvOnTarget(target)
	singleton.ty = reflect.TypeOf(target)
	singleton.data = interface{}(*target)

	return nil
}

func GetConfig[T Instance]() (T, error) {
	if singleton.data != nil {
		if singleton.ty != reflect.TypeOf(*new(T)) {
			return *new(T), fmt.Errorf("mismatch generic type: the given type is diferent from the singleton data")
		}
		return (singleton.data).(T), nil
	} else {
		return *new(T), fmt.Errorf("trying to access an unitialized config data")
	}
}

func loadAndParseConfig[T Instance](path string, target *T) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, target)
}
