package validator

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var enLoc = en.New()
var universalTranslator = ut.New(enLoc, enLoc)
var translator, _ = universalTranslator.GetTranslator("en")

var validate = validator.New()
var _ = en_translations.RegisterDefaultTranslations(validate, translator)

func ValidateStruct(anyStruct interface{}) error {
	return validate.Struct(anyStruct)
}

func IsValidationError(err error) bool {
	_, ok := err.(validator.ValidationErrors)
	return ok
}

func TranslateError(err error) error {
	errs := err.(validator.ValidationErrors)
	errorList := []string{}
	for _, e := range errs {
		errorList = append(errorList, e.Translate(translator))
	}
	errOut := strings.Join(errorList, ", ")
	return errors.New(errOut)
}

func IsEmail(email string) error {
	errs := validate.Var(email, "required,email")
	return errs
}
