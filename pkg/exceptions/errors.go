package exceptions

func NewJSONParseError() *ErrorResponse {
	return &ErrorResponse{
		Code: ErrCodeJSONParse,
	}
}

func NewDuplicateError(fields []string) *ErrorResponse {
	return &ErrorResponse{
		Code:             ErrCodeDuplicateEntry,
		DuplicatedFields: fields,
	}
}
