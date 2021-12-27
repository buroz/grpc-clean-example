package validator

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/buroz/grpc-clean-example/pkg/exceptions"
	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

func init() {
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	validate.RegisterValidation("weak_password", func(fl validator.FieldLevel) bool {
		return validPassword(fl.Field().String())
	})
}

func validPassword(s string) bool {
next:
	for _, classes := range map[string][]*unicode.RangeTable{
		"upper case": {unicode.Upper, unicode.Title},
		"lower case": {unicode.Lower},
		"numeric":    {unicode.Number, unicode.Digit},
		"special":    {unicode.Space, unicode.Symbol, unicode.Punct, unicode.Mark},
	} {
		for _, r := range s {
			if unicode.IsOneOf(classes, r) {
				continue next
			}
		}
		return false
	}

	return true
}

func ValidateStruct(model interface{}) *exceptions.ErrorResponse {
	var errs []exceptions.ValidationError

	err := validate.Struct(model)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			element := exceptions.ValidationError{
				FieldName: strings.Join(strings.Split(err.Namespace(), ".")[1:], "."),
				Cause:     err.Tag(),
			}

			errs = append(errs, element)
		}
	}

	if len(errs) > 0 {
		return &exceptions.ErrorResponse{
			Code:          exceptions.ErrCodeValidation,
			InvalidFields: errs,
		}
	}

	return nil
}
