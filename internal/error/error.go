package error

import (
	"fmt"
	"net/http"
)

type ServerAnswer struct {
	Code    int
	Message string
}

func (se ServerAnswer) Error() string {
	return fmt.Sprintf("error with code %d happened: %s", se.Code, se.Message)
}

var (
	// определяем ошибки баз данных
	DatabaseError error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "database error",
	}

	InvalidQuery error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "invalid query",
	}

	EmptyQuery error = ServerAnswer{
		Code:    http.StatusNotFound,
		Message: "empty rows",
	}

	NotUpdated error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "not apdated",
	}

	// определяем ошибки уровня usecase
	NotExist error = ServerAnswer{
		Code:    http.StatusNotFound,
		Message: "not exist",
	}

	AlreadyExist error = ServerAnswer{
		Code:    http.StatusForbidden,
		Message: "already exist",
	}

	InternalError error = ServerAnswer{
		Code:    http.StatusInternalServerError,
		Message: "internal error",
	}

	PasswordMismatch error = ServerAnswer{
		Code:    http.StatusUnauthorized,
		Message: "password mismatch",
	}

	// определяем ошибки уровня http
	BadRequest error = ServerAnswer{
		Code:    http.StatusBadRequest,
		Message: "bad request",
	}

	Unauthorized error = ServerAnswer{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}
)

func ToMetaStatus(err error) (int, string) {
	answer, ok := err.(ServerAnswer)
	if ok {
		return answer.Code, answer.Message
	}
	return http.StatusInternalServerError, answer.Error()
}
