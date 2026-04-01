package utils

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var data map[string]interface{}

func init() {
	f, err := os.ReadFile("assets/strings.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to read strings.yaml: %v", err))
	}
	if err := yaml.Unmarshal(f, &data); err != nil {
		panic(fmt.Sprintf("failed to unmarshal strings.yaml: %v", err))
	}
}

func T(lang string, path string) string {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else {
			panic(fmt.Sprintf("key %s not found in strings.yaml", path))
		}
	}

	if translations, ok := current.(map[string]interface{}); ok {
		if val, ok := translations[lang].(string); ok {
			return val
		}
		return translations["en"].(string)
	}

	if val, ok := current.(string); ok {
		return val
	}

	return ""
}
