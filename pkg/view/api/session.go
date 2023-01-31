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
		return err
	}

	res, err := Client.Post(BaseUrl+"/sessions", AppJsonType, bytes.NewBuffer(sessionJSON))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}

	var dest oapiGen.CreateSession400JSONResponse
	switch res.StatusCode {
	case http.StatusInternalServerError:
		return ErrInternalServerError

	case http.StatusBadRequest:
		dest = oapiGen.CreateSession400JSONResponse{}
	}

	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&dest)
	if err != nil {
		return err
	}

	return &BadRequestSession{Errs: dest.AdditionalProperties}
}

type BadRequestSession struct {
	Errs map[string]string
}

func (r *BadRequestSession) Error() string {
	return fmt.Sprintf("invalid session, %v", r.Errs)
}

func (r *BadRequestSession) Is(tgt error) bool {
	_, ok := tgt.(*BadRequestSession)
	return ok
}
