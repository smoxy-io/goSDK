package env

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type IEnvVal interface {
	string | bool | int | uint | []string | []bool | []int | []uint | map[string]any | []byte
}

func Set[T IEnvVal](key string, value T) error {
	val := ""

	switch v := any(value).(type) {
	case bool:
		val = strconv.FormatBool(v)
	case int:
		val = strconv.Itoa(v)
	case uint:
		val = strconv.FormatUint(uint64(v), 10)
	case string:
		val = v
	case []byte:
		val = string(v)
	case []string, []bool, []int, []uint, map[string]any:
		if b, err := json.Marshal(value); err == nil {
			val = string(b)
		}
	default:
		// should be impossible, but just in case
		return fmt.Errorf("invalid type for environment variable value: %T", value)
	}

	return os.Setenv(key, val)
}

func Get(key string, defVal ...string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	if len(defVal) == 0 {
		return ""
	}

	return defVal[0]
}

// GetAs Gets the value of an environment variable as the desired type.
// Returns the default value if the environment variable is not set or if there is an error converting the value
// to the desired type.
func GetAs[T IEnvVal](key string, defVal T) T {
	if val := os.Getenv(key); val != "" {
		switch any(defVal).(type) {
		case bool:
			b, err := strconv.ParseBool(val)

			if err != nil {
				return defVal
			}

			return any(b).(T)

		case int:
			i, err := strconv.Atoi(val)

			if err != nil {
				return defVal
			}

			return any(i).(T)

		case uint:
			u, err := strconv.ParseUint(val, 10, 64)

			if err != nil {
				return defVal
			}

			return any(u).(T)

		case string, []byte:
			return any(val).(T)

		case []string, []bool, []int, []uint, map[string]any:
			var d T

			if err := json.Unmarshal([]byte(val), &d); err != nil {
				return defVal
			}

			return d

		default:
			// should be impossible, but just in case
			return defVal
		}
	}

	return defVal
}
