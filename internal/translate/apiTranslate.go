package translate

import (
	"fmt"
	"google.golang.org/api/translate/v2"
)

type Translator struct {
	client *translate.Translator
}

// NewTranslator initializes the go-free-translate client
func NewTranslator() (*Translator, error) {
	client := translate.New()
	return &Translator{client: client}, nil
}

// TranslateText translates the given text into the target language
func (t *Translator) TranslateText(text, targetLang string) (string, error) {
	translation, err := t.client.Translate(text, "auto", targetLang)
	if err != nil {
		return "", fmt.Errorf("translation error: %v", err)
	}
	return translation, nil
}
