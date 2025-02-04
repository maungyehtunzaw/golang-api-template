package middlewares

import (
	"strings"

	"golang-api-template/internal/i18n"

	"github.com/gin-gonic/gin"
)

// LocaleMiddleware extracts the desired language from Accept-Language header
func LocaleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		langHeader := c.GetHeader("Accept-Language")
		if langHeader == "" {
			langHeader = "en" // default language
		} else {
			// e.g., parse if "en-US" => "en"
			parts := strings.Split(langHeader, ",")
			if len(parts) > 0 {
				langHeader = strings.Split(parts[0], "-")[0]
			}
		}
		i18n.SetLocale(c, langHeader)
		c.Next()
	}
}
