package exceptions

type ValidationError struct {
	FieldName string `json:"field_name,omitempty"`
	Cause     string `json:"cause,omitempty"`
}

type ErrorResponse struct {
	Code             ErrorCode         `json:"code"`
	InvalidFields    []ValidationError `json:"invalid_fields,omitempty"`
	DuplicatedFields []string          `json:"duplicated_fields,omitempty"`
}
