package i18n

import (
	"fmt"
	"reflect"
	"strings"
)

type Lang struct {
	name string
	*Translation
}

type langStore struct {
	defaultLang Lang
	langs       map[string]Lang
}

var defaultStore = langStore{
	langs: make(map[string]Lang),
}

func SetDefaultLang(lang Lang) {
	defaultStore.defaultLang = lang
}

func AddLang(lang Lang) error {
	if _, ok := defaultStore.langs[strings.ToLower(lang.name)]; ok {
		return fmt.Errorf("the lang `%s` already exist", lang.name)
	}
	defaultStore.langs[strings.ToLower(lang.name)] = lang
	if defaultStore.defaultLang.name == "" {
		defaultStore.defaultLang = lang
	}
	return nil
}

func GetLang(name string) (Lang, bool) {
	if lang, ok := defaultStore.langs[strings.ToLower(name)]; ok {
		return lang, true
	}
	return Lang{}, false
}

func (l Lang) Tr(format string, args ...interface{}) string {
	value, ok := l.Get(format)
	if ok {
		format = value
	}

	if len(args) > 0 {
		params := make([]interface{}, 0, len(args))
		for _, arg := range args {
			if arg == nil {
				continue
			}

			val := reflect.ValueOf(arg)
			if val.Kind() == reflect.Slice {
				for i := 0; i < val.Len(); i++ {
					params = append(params, val.Index(i).Interface())
				}
			} else {
				params = append(params, arg)
			}
		}
		return fmt.Sprintf(format, params...)
	}
	return format
}

type Translation struct {
	sections []string
	words    map[string]string
}

func (t *Translation) clone() *Translation {
	secs := make([]string, len(t.sections))
	copy(secs, t.sections)
	return &Translation{
		sections: secs,
		words:    t.words,
	}
}

func (t *Translation) Section(sec string) *Translation {
	t0 := t.clone()
	t0.sections = append(t0.sections, sec)
	return t0
}

func (t *Translation) Get(name string) (string, bool) {
	sections := t.sections
	sections = append(sections, name)
	key := strings.Join(sections, ".")
	v, ok := t.words[key]
	return v, ok
}

// Tr translates content to target language.
func Tr(format string, args ...interface{}) string {
	defLang := defaultStore.defaultLang
	if defLang.name == "" {
		defLang = en_us
	}
	_, ok := defLang.Get(format)
	if ok {
		return defLang.Tr(format, args...)
	}
	return en_us.Tr(format, args...)
}
