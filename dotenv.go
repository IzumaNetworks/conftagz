package conftagz

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// ParseDotEnv parses a .env formatted reader into a map of key-value pairs.
// Supports comments (#), blank lines, export prefix, and single/double quoted values.
func ParseDotEnv(r io.Reader) (map[string]string, error) {
	envMap := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Strip optional "export " prefix
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(line[7:])
		}

		// Split on first '='
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Strip matching quotes
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		envMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envMap, nil
}

// LoadDotEnvFile loads a .env file from the given path.
// Returns an empty map and nil error if the file does not exist.
// Returns an error for permission problems or parse errors.
func LoadDotEnvFile(filepath string) (map[string]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]string), nil
		}
		return nil, err
	}
	defer f.Close()

	return ParseDotEnv(f)
}

// MergeMaps merges two maps. Values in the override map take priority.
func MergeMaps(base, override map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(override))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range override {
		merged[k] = v
	}
	return merged
}
