package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadDotEnv(name string) error {
	path := name
	if !filepath.IsAbs(name) {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		found := false
		for {
			candidate := filepath.Join(wd, name)
			if _, err := os.Stat(candidate); err == nil {
				path = candidate
				found = true
				break
			}
			parent := filepath.Dir(wd)
			if parent == wd {
				break
			}
			wd = parent
		}
		if !found {
			return nil // .env not found – ignore silently
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open env file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
	return nil
}
