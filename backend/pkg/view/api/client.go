package api

import (
	"errors"
	"net/http"
	"time"
)

var BaseUrl = "http://localhost:42069"

var ErrInternalServerError = errors.New("internal server error")
var ErrServerUnreachable = errors.New("server isn't reachable")
var ErrNotFound = errors.New("resource not found")
var ErrEditConflict = errors.New("edit conflict")

var Client = http.Client{
	Timeout: 5 * time.Second,
}
