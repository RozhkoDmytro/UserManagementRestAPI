package apperrors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Message  string
	Code     string
	HTTPCode int
}

var (
	EnvConfigLoadError = AppError{
		Message: "Failed to load env file",
		Code:    "ENV_INIT_ERR",
	}

	EnvConfigVarError = AppError{
		Message: "CONFIG_PATH hasn't been found in environment variables",
		Code:    "ENV_CONFIG_VAR_ERR",
	}

	EnvConfigParseError = AppError{
		Message: "Failed to parse env file",
		Code:    "ENV_PARSE_ERR",
	}

	NilPostgresConfigError = AppError{
		Message: "Postgres config cannot be nil",
		Code:    "NIL_POSTGRES_ERR",
	}

	LoggerInitError = AppError{
		Message: "Cannot init logger",
		Code:    "LOGGER_INIT_ERR",
	}

	InsertionFailedErr = AppError{
		Message:  "Insertion operation has been failed",
		Code:     "INSERTION_ERR_FAILED",
		HTTPCode: http.StatusInternalServerError,
	}

	DeletionFailedErr = AppError{
		Message:  "Deletion failed",
		Code:     "DELETION_FAILED",
		HTTPCode: 500,
	}

	NoRecordFoundErr = AppError{
		Message:  "No record found",
		Code:     "NO_RECORD_FOUND",
		HTTPCode: 404,
	}
)

func (appError *AppError) Error() string {
	return appError.Code + ": " + appError.Message
}

func (appError *AppError) AppendMessage(anyErrs ...interface{}) *AppError {
	return &AppError{
		Message: fmt.Sprintf("%v : %v", appError.Message, anyErrs),
		Code:    appError.Code,
	}
}

func Is(err1 error, err2 *AppError) bool {
	err, ok := err1.(*AppError)
	if !ok {
		return false
	}

	return err.Code == err2.Code
}
