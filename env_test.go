package goyamlcfg

import (
	"fmt"
	"strings"
	"testing"
)

func TestEnvRegexPass(t *testing.T) {
	envs := []string{
		"${HELLO}",
		"${HELLO:-hello}",
		"${HELLO:-hello world}",
		"${HELLO:-hello:-world}",
		"inner env ${HELLO:-hello}",
		"${HELLO:-hello} ${WORLD:-world}",
	}

	for _, env := range envs {
		if !envRegex.MatchString(env) {
			t.Errorf("Unable to match %s with the env regex", env)
		}
	}
}

func TestEnvRegexFail(t *testing.T) {
	envs := []string{
		"hello",
		"",
		"${}",
	}

	for _, env := range envs {
		if envRegex.MatchString(env) {
			t.Errorf("Wrong match %s with the env regex", env)
		}
	}
}

func TestStripEnv(t *testing.T) {
	envKeys := []string{
		"HELLO",
		"HELLO:-world",
		"HELLO:-hello world",
		"HELLO:-hello:-world",
	}

	for _, key := range envKeys {
		env := fmt.Sprintf("${%s}", key)
		result := stripEnvPattern(env)
		if result != key {
			t.Errorf("Unexpected strip env pattern. Expected: %s, found: %s", key, result)
		}
	}
}

func TestEnvDefaultSplit(t *testing.T) {
	defaults := []string{
		"hello",
		"hello:-world",
		"hello world",
	}

	for _, d := range defaults {
		env, def := splitEnvDefault(fmt.Sprintf("HELLO:-%s", d))
		if env != "HELLO" {
			t.Errorf("The default env do not match. Expected: HELLO, found: %s", env)
		}
		if def != d {
			t.Errorf("The default value do not match. Expected: %s, found: %s", d, def)
		}
	}
}

func TestEnvStringReplace(t *testing.T) {
	envs := []string{
		"HELLO",
		"WORLD",
	}

	// setting env variables
	for _, env := range envs {
		t.Setenv(env, strings.ToLower(env))
	}

	for _, env := range envs {
		value := replaceEnv(fmt.Sprintf("${%s}", env))
		if value != strings.ToLower(env) {
			t.Errorf("The env replace do not match. Expected: %s, found: %s", env, value)
		}
	}
}

func TestEnvStringReplaceInner(t *testing.T) {
	t.Setenv("SOME_VAR", "someVar")

	input := "some thing ${SOME_VAR}"
	expected := "some thing someVar"

	result := replaceEnv(input)
	if result != expected {
		t.Errorf("Mismatch inner replacement. Expected: %s found: %s", expected, result)
	}
}

func TestEnvStringReplaceMutiple(t *testing.T) {
	t.Setenv("HELLO", "hello")
	t.Setenv("WORLD", "world")

	input := "${HELLO} ${WORLD}"
	expected := "hello world"

	result := replaceEnv(input)
	if result != expected {
		t.Errorf("Mismatch inner replacement. Expected: %s found: %s", expected, result)
	}
}

func TestEnvStringReplaceDafault(t *testing.T) {
	envs := []struct {
		key string
		set bool
	}{
		{"HELLO:-hello", true},
		{"WORLD:-world", false},
	}

	// setting env variables
	for _, env := range envs {
		key, val := splitEnvDefault(env.key)
		if env.set {
			t.Setenv(key, val)
		}
	}

	for _, env := range envs {
		_, val := splitEnvDefault(env.key)
		result := replaceEnv(fmt.Sprintf("${%s}", env.key))
		if result != val {
			t.Errorf("Env replace dont match for '${%s}'. Expected: %s, found: %s", env.key, val, result)
		}
	}
}

func TestReplaceOnPlainTarget(t *testing.T) {
	// setting env variables
	envs := map[string]string{
		"PASSWD": "12345",
	}

	for key, val := range envs {
		t.Setenv(key, val)
	}

	type Target struct {
		Username string
		Password string
	}

	// Replace
	target := Target{
		Username: "John-Doe",
		Password: "${PASSWD}",
	}

	replaceEnvOnTarget(&target)

	if target.Password != envs["PASSWD"] {
		t.Errorf("Target field dont match with the env. Expected: %s, found: %s", envs["PASSWD"], target.Password)
	}
}

func TestReplaceOnNestedTarget(t *testing.T) {
	// setting env variables
	envs := map[string]string{
		"PASSWD": "12345",
	}

	for key, val := range envs {
		t.Setenv(key, val)
	}

	type Target struct {
		Name string
		Auth struct {
			Username string
			Password string
		}
	}

	// Replace
	target := Target{
		Name: "John Doe",
		Auth: struct {
			Username string
			Password string
		}{
			Username: "John-Doe",
			Password: "${PASSWD}",
		},
	}

	replaceEnvOnTarget(&target)

	if target.Auth.Password != envs["PASSWD"] {
		t.Errorf("Target field dont match with the env. Expected: %s, found: %s", envs["PASSWD"], target.Auth.Password)
	}
}
