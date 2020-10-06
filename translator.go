package translator

// Translator interface
type Translator interface {
	// Register new translation message for locale
	// Use placeholder in message for field name
	// @example:
	// t.Register("en", "welcome", "Hello {name}, welcome!")
	Register(locale string, key string, message string)
	// Resolve find translation for locale
	// if no translation found for locale return fallback translation or nil
	Resolve(locale string, key string) string
	// ResolveStruct find translation from translatable
	// if empty string returned from translatable or struct not translatable, default translation will resolved
	ResolveStruct(s interface{}, locale string, key string) string
	// Translate get translation for locale
	// @example:
	// t.Translate("en", "welcome", map[string]string{ "name": "John" })
	Translate(locale string, key string, placeholders map[string]string) string
	// TranslateStruct translate using translatable interface
	// if empty string returned from translatable or struct not translatable, default translation will resolved
	// Caution: use non-pointer implemantation for struct
	TranslateStruct(s interface{}, locale string, key string, placeholders map[string]string) string
}
