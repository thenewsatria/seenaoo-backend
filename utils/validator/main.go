package validator

import (
	"errors"
	"fmt"
	"mime/multipart"
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

func ValidateContentType(file *multipart.FileHeader, allowedContentType []string) error {
	for _, contentType := range allowedContentType {
		if contentType == file.Header.Get("Content-Type") {
			return nil
		}
	}
	errOut := fmt.Sprintf("%s with the type of %s is not allowed", file.Filename, file.Header.Get("Content-Type"))
	return errors.New(errOut)
}

func ValidateFiles(files []*multipart.FileHeader, maxSize int64, allowedContentType []string) error {
	errorList := []string{}
	for _, file := range files {
		err := ValidateContentType(file, allowedContentType)
		if err != nil {
			errorList = append(errorList, err.Error())
		}
		if file.Size >= maxSize {
			errorStr := fmt.Sprintf("%s is too big, please upload a file less than %d kb", file.Filename, maxSize/1024)
			errorList = append(errorList, errorStr)
		}
	}
	if len(errorList) == 0 {
		return nil
	}

	errOut := strings.Join(errorList, ", ")
	return errors.New(errOut)
}
