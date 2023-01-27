package api

import (
	"bytes"
	"encoding/json"
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

	if res.StatusCode != http.StatusOK {
		return err
	}

	return nil
}
