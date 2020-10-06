package translator

import "strings"

type record struct {
	locale  string
	key     string
	message string
}

type memTranslator struct {
	fallback string
	records  []record
}

func (t *memTranslator) init(fallbackLocale string) {
	t.fallback = fallbackLocale
	t.records = make([]record, 0)
}

// Register new translation message for locale
// Use placeholder in message for field name
// @example:
// t.Register("en", "welcome", "Hello {name}, welcome!")
func (t *memTranslator) Register(locale string, key string, message string) {
	t.records = append(t.records, record{
		locale:  locale,
		key:     key,
		message: message,
	})
}

// Resolve find translation for locale
// if no translation found for locale return fallback translation or nil
func (t *memTranslator) Resolve(locale string, key string) string {
	for _, r := range t.records {
		if r.locale == locale && r.key == key {
			return r.message
		}
	}

	if locale != t.fallback {
		return t.Resolve(t.fallback, key)
	}

	return ""
}

// ResolveStruct find translation from translatable
// if empty string returned from translatable or struct not translatable, default translation will resolved
func (t *memTranslator) ResolveStruct(s interface{}, locale string, key string) string {
	if tr := resolveTranslatable(s); tr != nil {
		tr := tr.GetTranslation(locale, key)
		if tr != "" {
			return tr
		}
	}
	return t.Resolve(locale, key)
}

// Translate get translation for locale
// @example:
// t.Translate("en", "welcome", map[string]string{ "name": "John" })
func (t *memTranslator) Translate(locale string, key string, placeholders map[string]string) string {
	message := t.Resolve(locale, key)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}

// TranslateStruct translate using translatable interface
// if empty string returned from translatable or struct not translatable, default translation will resolved
// Caution: use non-pointer implemantation for struct
func (t *memTranslator) TranslateStruct(s interface{}, locale string, key string, placeholders map[string]string) string {
	message := t.ResolveStruct(s, locale, key)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}
