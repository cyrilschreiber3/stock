package locales

import (
	"unicode"
	"unicode/utf8"

	"github.com/invopop/ctxi18n/i18n"
)

func (l *Locale) T(key string, args ...any) string {
	ctx := l.Context
	return i18n.T(ctx, key, args...)
}

func (l *Locale) Tc(key string, args ...any) string {
	return capitalizeFirst(l.T(key, args...))
}

func (l *Locale) N(key string, count int, args ...any) string {
	return i18n.N(l.Context, key, count, args...)
}

func (l *Locale) Nc(key string, count int, args ...any) string {
	return capitalizeFirst(l.N(key, count, args...))
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}
