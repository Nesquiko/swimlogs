package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

func UpdateSession(s oapiGen.Session) (oapiGen.Session, error) {
	reqJson, err := json.Marshal(s)
	if err != nil {
		return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPut,
		BaseUrl+"/sessions/"+s.Id.String(),
		bytes.NewBuffer(reqJson),
	)
	if err != nil {
		return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", err)
	}
	req.Header.Add(ContentType, AppJsonType)

	res, err := Client.Do(req)
	if err != nil {
		return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", ErrServerUnreachable)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusNotFound:
		return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", ErrNotFound)
	case http.StatusConflict:
		return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", ErrEditConflict)
	case http.StatusInternalServerError:
		return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", ErrInternalServerError)
	}

	dec := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusBadRequest {
		dest := oapiGen.InvalidSession{}
		err = dec.Decode(&dest)
		if err != nil {
			return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", err)
		}

		return oapiGen.Session{}, &BadRequestSession{dest}
	}

	dest := oapiGen.Session{}
	err = dec.Decode(&dest)
	if err != nil {
		return oapiGen.Session{}, fmt.Errorf("UpdateSession: %w", err)
	}

	return dest, nil
}

func GetSessions() ([]oapiGen.Session, error) {
	res, err := http.Get(BaseUrl + "/sessions")
	if err != nil {
		return nil, fmt.Errorf("GetSessions: %w", ErrServerUnreachable)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("GetSessions: %w", ErrInternalServerError)
	}

	var dest oapiGen.GetAllSessions200JSONResponse
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return nil, fmt.Errorf("GetSessions: %w", err)
	}

	return dest.Sessions, nil
}

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
