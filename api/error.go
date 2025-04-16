package api

type ErrorCode string

const (
	ERR_UNKNOWN   ErrorCode = "ERR_UNKNOWN"
	ERR_INVALID   ErrorCode = "ERR_INVALID"
	ERR_NOT_FOUND ErrorCode = "ERR_NOT_FOUND"
	ERR_CONFLICT  ErrorCode = "ERR_CONFLICT"
)

type ErrorDto struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
	Details []string  `json:"details,omitempty"`
}

type Error interface {
	error
	ToErrorDto() ErrorDto
}

func (e ErrorDto) Error() string {
	return e.Message
}

func (e ErrorDto) ToErrorDto() ErrorDto {
	return e
}

func NewError(message string, code ErrorCode) ErrorDto {
	return ErrorDto{
		Message: message,
		Code:    code,
	}
}

func NewValidationError(message string, details []error) ErrorDto {
	stringDetails := make([]string, len(details))
	for i, err := range details {
		stringDetails[i] = err.Error()
	}
	return ErrorDto{
		Message: message,
		Code:    ERR_INVALID,
		Details: stringDetails,
	}
}
