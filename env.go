package goyamlcfg

import (
	"os"
	"reflect"
	"regexp"
	"strings"
)

var envRegex, _ = regexp.Compile(`\$\{[^\}]+\}`)
var defaultRegex, _ = regexp.Compile(`^\$\{.+\:\-.+\}$`)

func replaceEnvOnTarget(target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(replaceEnv(v.String()))

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			replaceEnvOnTarget(field.Addr().Interface())
		}

	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			replaceEnvOnTarget(elem.Addr().Interface())
		}

	default:
		return
	}

}

func replaceEnv(field string) string {
	return envRegex.ReplaceAllStringFunc(field, func(match string) string {
		env := stripEnvPattern(match)
		if defaultRegex.MatchString(match) {
			key, def := splitEnvDefault(env)
			return envOrDefault(key, def)
		}
		return envOrDefault(env, "")
	})
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
