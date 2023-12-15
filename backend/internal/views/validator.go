package views

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	AppValidator        validator.Validate
	defaultTrans        ut.Translator
	univarsalTranslator ut.UniversalTranslator
)

type Validatable interface {
	Validate() []InvalidParam
}

func init() {
	v := validator.New()
	// locale
	en := en.New()
	// Set default
	univarsalTranslator = *ut.New(en, en)
	enTrans, _ := univarsalTranslator.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(v, enTrans)

	defaultTrans = enTrans
	// TODO 日付とか

	AppValidator = *v
}

func ToInvalidParams(errs validator.ValidationErrors) []InvalidParam {
	invalidParams := []InvalidParam{}
	for _, err := range errs {
		msg := err.Translate(defaultTrans)
		invalidParams = append(invalidParams, InvalidParam{
			Reason: msg,
			Field:  err.Field(),
		})
	}
	return invalidParams
}
