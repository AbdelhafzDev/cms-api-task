package i18nutil

import (
	"database/sql"
	"strings"
)

const (
	LangEnglish = "en"
	LangArabic  = "ar"
)

func NormalizeLanguage(lang string) string {
	switch lang {
	case "arabic", "ar":
		return LangArabic
	case "english", "en", "":
		return LangEnglish
	default:
		return LangEnglish
	}
}

func IsArabic(lang string) bool {
	return NormalizeLanguage(lang) == LangArabic
}

func LocalizedString(isArabic bool, ar, en string) string {
	if isArabic {
		return ar
	}
	return en
}

func LocalizedNullString(isArabic bool, ar, en sql.NullString) string {
	if isArabic && ar.Valid {
		return ar.String
	}
	if en.Valid {
		return en.String
	}
	return ""
}

func LocalizedNullStringPtr(isArabic bool, ar, en sql.NullString) *string {
	if isArabic && ar.Valid && ar.String != "" {
		return &ar.String
	}
	if !isArabic && en.Valid && en.String != "" {
		return &en.String
	}
	return nil
}

// ParseAcceptLanguage extracts app language from an Accept-Language header value.
// It supports locale forms like "ar-SA" and "en-US", defaulting to English.
func ParseAcceptLanguage(acceptLanguage string) string {
	if strings.TrimSpace(acceptLanguage) == "" {
		return LangEnglish
	}

	first := strings.TrimSpace(strings.Split(acceptLanguage, ",")[0])
	if idx := strings.Index(first, ";"); idx != -1 {
		first = first[:idx]
	}

	first = strings.ToLower(strings.TrimSpace(first))
	if strings.HasPrefix(first, "ar") {
		return LangArabic
	}

	return LangEnglish
}
