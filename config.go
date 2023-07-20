package goyamlcfg

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Instance interface{}

type config struct {
	mu   sync.Mutex
	data interface{}
}

var singleton config

func InitializeConfigSingleton[T Instance](path string) error {
	singleton.mu.Lock()
	defer singleton.mu.Unlock()

	target := new(T)
	loadAndParseConfig(path, target)
	replaceEnvOnTarget(target)
	singleton.data = interface{}(*target)

	return nil
}

func GetConfig[T Instance]() T {
	if singleton.data != nil {
		return (singleton.data).(T)
	} else {
		panic("trying to access an unitialized config data")
	}
}

func loadAndParseConfig[T Instance](path string, target *T) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, target)
}
