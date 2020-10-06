package translator

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/gobardofw/utils"
	"github.com/tidwall/gjson"
)

type JsonDriver struct {
	fallback string
	dir      string
	jsonData string
}

func (t *JsonDriver) init(fallbackLocale string, dir string) error {
	t.fallback = fallbackLocale
	t.dir = dir
	return t.Load()
}

// Load load translations file to memory
func (t *JsonDriver) Load() error {
	var resolveFiles = func(dir string) (map[string]string, error) {
		dir = filepath.Dir(path.Join(dir, "some.txt"))
		res := make(map[string]string)
		files := utils.FindFile(dir, ".json")
		for _, f := range files {
			if filepath.Dir(f) != dir {
				continue
			}
			// get file info
			filePath := f
			fileName := filepath.Base(f)
			fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName))

			// read file
			bytes, err := ioutil.ReadFile(filePath)
			if err != nil {
				return nil, err
			}

			// validate json
			content := string(bytes)
			if !gjson.Valid(content) {
				return nil, errors.New("Invalid json for " + filePath)
			}

			// append to files list
			res[fileName] = content
		}
		return res, nil
	}

	var unwrapJson = func(jsonText string) (string, error) {
		var res map[string]interface{}
		if err := json.Unmarshal([]byte(jsonText), &res); err != nil {
			return "", err
		}

		if byts, err := json.Marshal(res); err != nil {
			return "", err
		} else {
			jsonStr := string(byts)
			return strings.TrimSuffix(strings.TrimPrefix(jsonStr, "{"), "}"), nil
		}
	}

	locales, err := utils.GetSubDirectory(t.dir)
	if err != nil {
		return err
	}

	locales = append(locales, "")
	contents := make([]string, 0)

	for _, locale := range locales {
		files, err := resolveFiles(path.Join(t.dir, locale))
		if err != nil {
			return err
		}

		if len(files) == 1 {
			for _, cnt := range files {
				if locale == "" {
					unwrappedJson, err := unwrapJson(cnt)
					if err != nil {
						return err
					}
					contents = append(contents, unwrappedJson)
				} else {
					contents = append(contents, `"`+locale+`":`+cnt)
				}
			}
		} else {
			if locale == "" {
				for file, cnt := range files {
					contents = append(contents, `"`+file+`":`+cnt)
				}
			} else {
				subContent := make([]string, 0)
				for file, cnt := range files {
					subContent = append(subContent, `"`+file+`":`+cnt)
				}
				contents = append(contents, `"`+locale+`":{`+strings.Join(subContent, ",")+"}")
			}
		}
	}

	t.jsonData = "{" + strings.Join(contents, ",") + "}"
	return nil
}

// Register new translation message for locale
// Use placeholder in message for field name
// @example:
// t.Register("en", "welcome", "Hello {name}, welcome!")
func (t *JsonDriver) Register(locale string, key string, message string) {
	// Do nothing
	// Json driver not support register
}

// Resolve find translation for locale
// if no translation found for locale return fallback translation or nil
func (t *JsonDriver) Resolve(locale string, key string) string {
	value := gjson.Get(t.jsonData, locale+"."+key)
	if !value.Exists() {
		value = gjson.Get(t.jsonData, t.fallback+"."+key)
	}
	return value.String()
}

// ResolveStruct find translation from translatable
// if empty string returned from translatable or struct not translatable, default translation will resolved
func (t *JsonDriver) ResolveStruct(s interface{}, locale string, key string) string {
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
func (t *JsonDriver) Translate(locale string, key string, placeholders map[string]string) string {
	message := t.Resolve(locale, key)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}

// TranslateStruct translate using translatable interface
// if empty string returned from translatable or struct not translatable, default translation will resolved
// Caution: use non-pointer implemantation for struct
func (t *JsonDriver) TranslateStruct(s interface{}, locale string, key string, placeholders map[string]string) string {
	message := t.ResolveStruct(s, locale, key)
	for p, v := range placeholders {
		message = strings.ReplaceAll(message, "{"+p+"}", v)
	}
	return message
}
