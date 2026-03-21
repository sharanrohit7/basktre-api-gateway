package config

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
)

var (
	once     sync.Once
	instance *Config
)

type Config struct {
	data map[string]interface{}
}

func Init(filePath string) error {
	var initErr error
	once.Do(func() {
		// Load .env file if it exists (local development)
		_ = godotenv.Load()

		baseData, err := loadFile(filePath)
		if err != nil {
			initErr = err
			return
		}

		// Look for .local version (e.g., config.yaml -> config.local.yaml)
		localPath := getLocalPath(filePath)
		if _, err := os.Stat(localPath); err == nil {
			localData, err := loadFile(localPath)
			if err == nil {
				mergeMaps(baseData, localData)
			}
		}

		instance = &Config{data: baseData}
	})
	return initErr
}

func loadFile(path string) (map[string]interface{}, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: failed to read file %s: %w", path, err)
	}

	data := make(map[string]interface{})
	if err := yaml.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("config: failed to parse YAML %s: %w", path, err)
	}
	if data == nil {
		data = make(map[string]interface{})
	}
	return data, nil
}

func getLocalPath(path string) string {
	if len(path) > 5 && path[len(path)-5:] == ".yaml" {
		return path[:len(path)-5] + ".local.yaml"
	}
	return path + ".local"
}

func mergeMaps(base, override map[string]interface{}) {
	for k, v := range override {
		if v == nil {
			delete(base, k)
			continue
		}
		if baseVal, ok := base[k]; ok {
			if baseMap, ok1 := baseVal.(map[string]interface{}); ok1 {
				if overrideMap, ok2 := v.(map[string]interface{}); ok2 {
					mergeMaps(baseMap, overrideMap)
					continue
				}
			}
		}
		base[k] = v
	}
}

func MustInit(filePath string) {
	if err := Init(filePath); err != nil {
		panic(err)
	}
}

func (c *Config) lookup(key string) (interface{}, bool) {
	if c == nil || c.data == nil {
		return nil, false
	}
	curr := c.data
	parts := strings.Split(key, ".")
	for idx, part := range parts {
		v, ok := curr[part]
		if !ok {
			return nil, false
		}
		if idx == len(parts)-1 {
			return v, true
		}
		next, ok := v.(map[string]interface{})
		if !ok {
			return nil, false
		}
		curr = next
	}
	return nil, false
}

func GetString(key string) string { return GetStringDefault(key, "") }

func GetStringDefault(key, def string) string {
	if env := os.Getenv(toEnvKey(key)); env != "" {
		return env
	}
	if instance == nil {
		return def
	}
	v, ok := instance.lookup(key)
	if !ok || v == nil {
		return def
	}
	return fmt.Sprintf("%v", v)
}

func GetBool(key string) bool {
	if strings.EqualFold(os.Getenv(toEnvKey(key)), "true") {
		return true
	}
	if instance == nil {
		return false
	}
	v, ok := instance.lookup(key)
	if !ok || v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}

func GetIntDefault(key string, def int) int {
	if instance == nil {
		return def
	}
	v, ok := instance.lookup(key)
	if !ok || v == nil {
		return def
	}
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	}
	return def
}

func GetStringSlice(key string) []string {
	if instance == nil {
		return nil
	}
	v, ok := instance.lookup(key)
	if !ok || v == nil {
		return nil
	}
	raw, ok := v.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		out = append(out, fmt.Sprintf("%v", item))
	}
	return out
}

func toEnvKey(key string) string {
	replacer := strings.NewReplacer(".", "_")
	return strings.ToUpper(replacer.Replace(key))
}
