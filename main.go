package translator

// NewMemoryTranslator create a new memory based translator
func NewMemoryTranslator(fallbackLocale string) Translator {
	t := new(memTranslator)
	t.init(fallbackLocale)
	return t
}

// NewJSONTranslator create a new memory based translator
func NewJSONTranslator(fallbackLocale string, dir string) (Translator, error) {
	t := new(JSONDriver)
	if err := t.init(fallbackLocale, dir); err != nil {
		return nil, err
	}
	return t, nil
}
