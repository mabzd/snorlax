package service

type ErrorCode string

const (
	ERR_UNKNOWN   ErrorCode = "unknown"
	ERR_INVALID   ErrorCode = "invalid"
	ERR_NOT_FOUND ErrorCode = "not_found"
	ERR_CONFLICT  ErrorCode = "conflict"
)

type ServiceError interface {
	error
	Code() ErrorCode
	Details() []error
}

func NewServiceError(message string, code ErrorCode) ServiceError {
	return &serviceErrorData{
		message: message,
		code:    code,
		details: nil,
	}
}

func NewValidationError(message string, details []error) ServiceError {
	return &serviceErrorData{
		message: message,
		code:    ERR_INVALID,
		details: details,
	}
}

type serviceErrorData struct {
	message string
	code    ErrorCode
	details []error
}

func (err *serviceErrorData) Error() string {
	return err.message
}

func (err *serviceErrorData) Code() ErrorCode {
	return err.code
}

func (err *serviceErrorData) Details() []error {
	return err.details
}
