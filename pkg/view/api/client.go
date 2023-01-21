package api

import (
	"net/http"
	"time"
)

var BaseUrl = "http://localhost:42069"

var Client = http.Client{
	Timeout: 5 * time.Second,
}
