package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// For simplicity, we embed translation JSON files
//
//go:embed en.json es.json my.json
var translationsFS embed.FS

var translations = map[string]map[string]string{}

// Initialize loads all translations into memory
func Initialize() error {
	files := []string{"en.json", "es.json", "my.json"}
	for _, file := range files {
		data, err := translationsFS.ReadFile(file)
		if err != nil {
			return err
		}
		var dict map[string]string
		if err := json.Unmarshal(data, &dict); err != nil {
			return err
		}
		lang := file[:2] // e.g. "en" from "en.json"
		translations[lang] = dict
	}
	return nil
}

func SetLocale(c *gin.Context, lang string) {
	c.Set("locale", lang)
}

func T(c *gin.Context, key string) string {
	lang, exists := c.Get("locale")
	if !exists {
		lang = "en"
	}
	loc := lang.(string)
	dict, ok := translations[loc]
	if !ok {
		// fallback to English
		dict = translations["en"]
	}
	val, ok := dict[key]
	if !ok {
		// fallback to key if missing
		return fmt.Sprintf("[Missing translation for %s in %s]", key, loc)
	}
	return val
}

func TT(lang, key string) string {
	if msgs, exists := translations[lang]; exists {
		if msg, exists := msgs[key]; exists {
			return msg
		}
	}
	return key // Fallback to the key itself if not found
}
