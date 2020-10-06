package translator

// Translatable interface for struct
type Translatable interface {
	GetTranslation(locale string, key string) string
}

func resolveTranslatable(s interface{}) Translatable {
	if v, ok := s.(Translatable); ok {
		return v
	}
	return nil
}
