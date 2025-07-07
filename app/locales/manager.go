package locales

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var text map[string]map[string]string

func Load() error {
	text = make(map[string]map[string]string)
	files, err := filepath.Glob("locales/*.json")
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no language files found in locales/ directory")
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		var langMap map[string]string
		if err := json.Unmarshal(data, &langMap); err != nil {
			return fmt.Errorf("failed to parse %s: %w", file, err)
		}

		langCode := filepath.Base(file)
		langCode = langCode[:len(langCode)-len(filepath.Ext(langCode))]
		text[langCode] = langMap
	}
	return nil
}

func Get(langCode, key string) string {
	if lang, ok := text[langCode]; ok {
		if value, ok := lang[key]; ok {
			return value
		}
	}
	// Fallback to English if the key is not found in the specified language
	if lang, ok := text["en"]; ok {
		if value, ok := lang[key]; ok {
			return value
		}
	}
	return key // Return the key itself if not found anywhere
}