package util

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

type GoValidator struct {
	Validator *validator.Validate
}

func NewGoValidator() *GoValidator {
	validator := &GoValidator{
		Validator: validator.New(),
	}

	validator.Validator.RegisterValidation("alphaSpace", alphaSpaceValidator)
	validator.Validator.RegisterValidation("username", usernameValidator)

	return validator
}

func alphaSpaceValidator(fl validator.FieldLevel) bool {

	matches, err := regexp.MatchString("^[a-zA-ZšŠđĐžŽčČćĆ\\s]+$", fl.Field().String())
	if err != nil {
		fmt.Println(err)
	}

	if !matches {
		return false
	}

	return true
}

func usernameValidator(fl validator.FieldLevel) bool {

	matches, err := regexp.MatchString("^[a-zA-Z\\d_.]+$", fl.Field().String())
	if err != nil {
		fmt.Println(err)
	}

	if !matches {
		return false
	}

	return true
}
