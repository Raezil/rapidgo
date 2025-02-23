package rapidgo

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func LoadEnv() error {
	return ReadLines(".env")
}

func LoadCustomEnv(filenames ...string) error {
	if len(filenames) == 0 {
		return errors.New("at least one filename must be provided")
	}

	for _, filename := range filenames {
		if filename == "" {
			return errors.New("filename cannot be empty")
		}
		if err := ReadLines(filename); err != nil {
			return err
		}
	}
	return nil
}

func ReadLines(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	envMappping := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = strings.Trim(value, `"'`)
		envMappping[key] = value
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for key, value := range envMappping {
		os.Setenv(key, value)
	}

	return nil
}
