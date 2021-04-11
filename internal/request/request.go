package request

import (
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

func ExtractPreferLanguage(r *http.Request) string {
	s := r.Header.Get("Accept-Language")
	if tags, _, _ := language.ParseAcceptLanguage(s); len(tags) > 0 {
		return tags[0].String()
	}

	return ""
}

func ExtractBearerToken(r *http.Request) string {
	s := r.Header.Get("Authorization")
	return strings.TrimPrefix(s, "Bearer ")
}
