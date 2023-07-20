package goyamlcfg

import (
	"os"
	"reflect"
	"regexp"
	"strings"
)

var envRegex, _ = regexp.Compile(`^\$\{.+\}$`)
var defaultRegex, _ = regexp.Compile(`^\$\{.+\:\-.+\}$`)

func replaceEnvOnTarget(target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			if envRegex.MatchString(field.String()) {
				field.SetString(replaceEnv(field.String()))
			}
		case reflect.Struct:
			replaceEnvOnTarget(field.Addr().Interface())
		}
	}
}

func replaceEnv(field string) string {
	env := stripEnvPattern(field)
	if defaultRegex.MatchString(field) {
		env, def := splitEnvDefault(env)
		return envOrDefault(env, def)
	}
	return envOrDefault(env, "")
}

func stripEnvPattern(env string) string {
	return strings.TrimSuffix(strings.TrimPrefix(env, "${"), "}")
}

func splitEnvDefault(env string) (string, string) {
	list := strings.Split(env, ":-")
	if len(list) > 2 {
		return list[0], strings.Join(list[1:], ":-")
	}
	return list[0], list[1]
}

func envOrDefault(env, defaultVal string) string {
	value := os.Getenv(env)
	if value == "" {
		return defaultVal
	}
	return value
}
