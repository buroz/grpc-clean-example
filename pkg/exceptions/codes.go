package exceptions

type ErrorCode int32

const (
	ErrCodeUnauthorized   ErrorCode = iota + 1
	ErrCodeValidation               // 2
	ErrCodeJSONParse                // 3
	ErrCodeInternalServer           // 4
	ErrCodeDuplicateEntry           // 5
	ErrCodeNotFound                 // 6

	// USER ERRORS
	ErrCodeWrongPassword // 7
)
