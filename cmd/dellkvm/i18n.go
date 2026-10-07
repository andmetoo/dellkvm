//go:build linux || windows

package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

const defaultLanguage = "en"

var supportedLanguages = []string{"en", "ru", "fr", "de", "zh"}

var (
	bundleOnce sync.Once
	bundle     *i18n.Bundle
	bundleErr  error
)

type appLocalizer struct {
	language  string
	localizer *i18n.Localizer
}

func defaultLocalizer() appLocalizer {
	l, err := newLocalizer(defaultLanguage)
	if err != nil {
		return appLocalizer{language: defaultLanguage}
	}
	return l
}

func newLocalizer(lang string) (appLocalizer, error) {
	normalized, err := normalizeLanguage(lang)
	if err != nil {
		return appLocalizer{}, err
	}

	b, err := i18nBundle()
	if err != nil {
		return appLocalizer{}, err
	}
	return appLocalizer{
		language:  normalized,
		localizer: i18n.NewLocalizer(b, normalized),
	}, nil
}

func (l appLocalizer) T(messageID string, data map[string]any) string {
	if l.localizer == nil {
		l = defaultLocalizer()
		if l.localizer == nil {
			return messageID
		}
	}
	if data == nil {
		data = map[string]any{}
	}

	message, err := l.localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return message
}

func i18nBundle() (*i18n.Bundle, error) {
	bundleOnce.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
		for _, lang := range supportedLanguages {
			path := fmt.Sprintf("locales/active.%s.toml", lang)
			if _, err := bundle.ParseMessageFileBytes([]byte(localeData[lang]), path); err != nil {
				bundleErr = err
				return
			}
		}
	})
	return bundle, bundleErr
}

func normalizeLanguage(lang string) (string, error) {
	lang = strings.TrimSpace(lang)
	if lang == "" {
		return defaultLanguage, nil
	}
	for _, supported := range supportedLanguages {
		if lang == supported {
			return lang, nil
		}
	}
	return "", fmt.Errorf("unsupported language %q; supported: %s", lang, supportedLanguageList())
}

func supportedLanguageList() string {
	return strings.Join(supportedLanguages, ", ")
}
