package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Nesquiko/swimlogs/pkg/api"
)

func encode[T any](w http.ResponseWriter, status int, response T) {
	encodeWithContentType(w, status, response, ApplicationJSON)
}

func encodeError(w http.ResponseWriter, err *ApiError) {
	encodeWithContentType(w, err.Status, err.ErrorDetail, ApplicationProblemJSON)
}

func encodeWithContentType[T any](
	w http.ResponseWriter,
	code int,
	response T,
	contentType string,
) {
	w.Header().Set(ContentType, contentType)
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("where", "encode"),
			slog.String("error", err.Error()),
		)
		http.Error(w, EncodingError, http.StatusInternalServerError)
	}
}

func Decode[T any](w http.ResponseWriter, r *http.Request) (T, *ApiError) {
	dst, err := decode[T](w, r)
	var decErr *decodeErr
	if errors.As(err, &decErr) {
		return dst, decodeErrToApiErrorWithCode(decErr, decErr.code)
	}
	if err != nil {
		return dst, decodeErrToApiError(err)
	}
	return dst, nil
}

type decodeErr struct {
	err  error
	code string
}

func (e *decodeErr) Error() string {
	return e.err.Error()
}

func decode[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var dst T
	r.Body = http.MaxBytesReader(w, r.Body, int64(MaxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		var invalidUnmarshalErr *json.InvalidUnmarshalError

		invalidFieldPrefix := "json: unknown field "
		largeBodyErrorStr := "http: request body too large"

		switch {
		case errors.As(err, &syntaxErr):
			return dst, fmt.Errorf(
				"body contains badly-formed JSON (at character %d)",
				syntaxErr.Offset,
			)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return dst, errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeErr):
			if unmarshalTypeErr.Field != "" {
				return dst, fmt.Errorf(
					"body contains incorrect JSON type for field %q",
					unmarshalTypeErr.Field,
				)
			}
			return dst, fmt.Errorf(
				"body contains incorrect JSON type (at character %d)",
				unmarshalTypeErr.Offset,
			)

		case errors.Is(err, io.EOF):
			return dst, errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), invalidFieldPrefix):
			fieldName := strings.TrimPrefix(err.Error(), invalidFieldPrefix)
			return dst, &decodeErr{
				fmt.Errorf("body contains unknown key %s", fieldName),
				ValidationErrorCode,
			}

		case err.Error() == largeBodyErrorStr:
			return dst, fmt.Errorf("body must not be larger than %d bytes", MaxBytes)

		case errors.As(err, &invalidUnmarshalErr):
			return dst, fmt.Errorf("invalid unmarshal target")

		default:
			return dst, err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return dst, errors.New("body must only contain a single JSON value")
	}

	return dst, nil
}

const (
	DecodingErrorCode  = "undecodable.request"
	DecodingErrorTitle = "Request couldn't be decoded"
)

func decodeErrToApiError(err error) *ApiError {
	return decodeErrToApiErrorWithCode(err, DecodingErrorCode)
}

func decodeErrToApiErrorWithCode(err error, code string) *ApiError {
	return &ApiError{
		ErrorDetail: api.ErrorDetail{
			Code:   code,
			Detail: err.Error(),
			Status: http.StatusBadRequest,
			Title:  DecodingErrorTitle,
		},
	}
}
