package env

import (
	"os"
	"strconv"
)

// GetEnvAsString returns env with the provided key.
// If env value is empty, it will returns default value.
func GetEnvAsString(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// GetEnvAsInt returns env with the provided key.
// If env value is empty, it will returns default value.
// Or if env value parse to int error, it will returns default value.
func GetEnvAsInt(key string, defaultValue int) (int, error) {
	if v := os.Getenv(key); v != "" {
		value, err := strconv.Atoi(v)
		if err != nil {
			return defaultValue, err
		}
		return value, nil
	}
	return defaultValue, nil
}

// IsEnvExist returns env exist
func IsEnvExist(key string) bool {
	_, exist := os.LookupEnv(key)
	return exist
}
