package app

import (
	"net/http"

	"github.com/Nesquiko/swimlogs/pkg/data"
	"github.com/Nesquiko/swimlogs/pkg/openapi"
)

func New(db *data.PostgresDbConn) SwimLogsApp {
	return SwimLogsApp{db}
}

type SwimLogsApp struct {
	db *data.PostgresDbConn
}

type Result[T any] struct {
	code     int
	body     T
	error    any
	hasError bool
}

func (r *Result[T]) IfSuccess(f func()) {
	if !r.hasError {
		f()
	}
}

func (r *Result[T]) HasError() bool {
	return r.hasError
}

func (r *Result[T]) Code() int {
	return r.code
}

// Returns response body, or error detail if an error occurred
func (r *Result[T]) Body() any {
	if r.hasError {
		return r.error
	}
	return r.body
}

func resultWithoutBody(code int) Result[struct{}] {
	return Result[struct{}]{code: code}
}

func resultWithBody[T any](body T, code int) Result[T] {
	return Result[T]{
		code: code,
		body: body,
	}
}

func resultWithError[T any](
	title, detail string,
	code int,
	extensions *map[string]any,
) Result[T] {
	return Result[T]{
		code:     code,
		hasError: true,
		error:    openapi.ErrorDetail{Title: title, Detail: detail, Extensions: extensions},
	}
}

func resultFromErrorDetails[T any](errDets openapi.ErrorDetail, code int) Result[T] {
	return Result[T]{
		code:     code,
		hasError: true,
		error:    errDets,
	}
}

func errorResult[T any](errRes any, code int) Result[T] {
	return Result[T]{
		code:     code,
		hasError: true,
		error:    errRes,
	}
}

func internalServerErrorResult[T any](detail string) Result[T] {
	return resultWithError[T]("Internal server error", detail, http.StatusInternalServerError, nil)
}
