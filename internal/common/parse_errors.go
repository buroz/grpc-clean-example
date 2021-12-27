package common

import (
	"github.com/buroz/grpc-clean-example/pkg/exceptions"
)

func ParseErrors(err *exceptions.ErrorResponse) *ServiceError {
	var validationErrors []*ValidationError
	var duplicatedFields []string

	if len(err.InvalidFields) > 0 {
		for _, field := range err.InvalidFields {
			validationErrors = append(validationErrors, &ValidationError{
				FieldName: field.FieldName,
				Cause:     field.Cause,
			})
		}
	}

	if len(err.DuplicatedFields) > 0 {
		duplicatedFields = append(duplicatedFields, err.DuplicatedFields...)
	}

	return &ServiceError{
		Code:             int32(err.Code),
		ValidationErrors: validationErrors,
		DuplicatedFields: duplicatedFields,
	}
}
