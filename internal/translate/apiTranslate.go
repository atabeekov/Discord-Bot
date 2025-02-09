package translate

import (
	"fmt"

	"github.com/bregydoc/gtranslate"
	"golang.org/x/text/language"
)

type Translator struct{}

// NewTranslator initializes the translator
func NewTranslator() *Translator {
	return &Translator{}
}

// TranslateText translates the given text into the target language
func (t *Translator) TranslateText(text, targetLang string) (string, error) {
	destLang, err := language.Parse(targetLang)
	if err != nil {
		return "", fmt.Errorf("invalid target language: %v", err)
	}

	translatedText, err := gtranslate.Translate(text, language.English, destLang)
	if err != nil {
		return "", fmt.Errorf("translation error: %v", err)
	}

	return translatedText, nil
}
