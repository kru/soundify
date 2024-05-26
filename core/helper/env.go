package helper

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		keyValue := strings.SplitN(line, "=", 2)
		if len(keyValue) != 2 {
			return fmt.Errorf("invalid line: %s", line)
		}

		key := strings.TrimSpace(keyValue[0])
		value := strings.TrimSpace(keyValue[1])

		err := os.Setenv(key, value)
		if err != nil {
			return fmt.Errorf("error setting variable: %s", err)
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return nil
}
