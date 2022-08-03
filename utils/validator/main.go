package validator

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func IsEmail(email string) error {
	errs := validate.Var(email, "required,email")
	return errs
}
