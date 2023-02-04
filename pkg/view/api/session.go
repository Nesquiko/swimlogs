package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

func CreateSession(s oapiGen.Session) error {
	sessionJSON, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}

	res, err := Client.Post(BaseUrl+"/sessions", AppJsonType, bytes.NewBuffer(sessionJSON))
	if err != nil {
		return fmt.Errorf("CreateSession: %w", ErrServerUnreachable)
	}
	defer res.Body.Close()

	var dest oapiGen.CreateSession400JSONResponse
	switch res.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusInternalServerError:
		return fmt.Errorf("CreateSession: %w", ErrInternalServerError)
	case http.StatusBadRequest:
		dest = oapiGen.CreateSession400JSONResponse{}
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}

	return &BadRequestSession{oapiGen.InvalidSession(dest.InvalidSessionErrorResponseJSONResponse)}
}

type BadRequestSession struct {
	oapiGen.InvalidSession
}

func (r *BadRequestSession) Error() string {
	return fmt.Sprintf("invalid session, %v", r.InvalidSession)
}

func (r *BadRequestSession) Is(tgt error) bool {
	_, ok := tgt.(*BadRequestSession)
	return ok
}
