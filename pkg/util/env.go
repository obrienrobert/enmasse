/*
 * Copyright 2019, EnMasse authors.
 * License: Apache License 2.0 (see the file LICENSE or http://apache.org/licenses/LICENSE-2.0.html).
 */

package util

import (
	"fmt"
	"os"
	"strings"
)

type ApplyEnvFn func(key string, value string, ok bool) error

// ApplyEnv retrieves the value of the environment variable named
// by the key and applies the key, value (which may be empty) and
// existence flag to the consumer function. The consumer function
// is called even if the key is not present in the environment,
// in which case the existence flag will be false.
func ApplyEnv(key string, consumer ApplyEnvFn) error {
	value, ok := os.LookupEnv(key)
	return consumer(key, value, ok)
}

func GetEnvOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	} else {
		return defaultValue
	}
}

func GetBooleanEnvOrDefault(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	} else {
		return "true" == strings.ToLower(value)
	}
}

func GetEnvOrError(name string) (string, error) {
	result := os.Getenv(name)
	if result != "" {
		return result, nil
	} else {
		return "", fmt.Errorf("'%s' is not set", name)
	}
}

func GetBooleanEnv(key string) bool {
	return GetBooleanEnvOrDefault(key, false)
}

type EnvironmentProvider interface {
	LookupEnv(string) (string, bool)
	Get(string) string
}

// OS environment provider

type OSEnvironmentProvider struct {
}

func (p OSEnvironmentProvider) LookupEnv(name string) (string, bool) {
	return os.LookupEnv(name)
}

func (p OSEnvironmentProvider) Get(name string) string {
	return os.Getenv(name)
}

var _ EnvironmentProvider = OSEnvironmentProvider{}

// mock environment provider

type MockEnvironmentProvider struct {
	Environment map[string]string
}

func (p MockEnvironmentProvider) LookupEnv(name string) (string, bool) {
	val, ok := p.Environment[name]
	return val, ok
}

func (p MockEnvironmentProvider) Get(name string) string {
	return p.Environment[name]
}

var _ EnvironmentProvider = MockEnvironmentProvider{}
