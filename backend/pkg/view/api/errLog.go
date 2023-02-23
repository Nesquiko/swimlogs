package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

const ContentType = "Content-Type"
const AppJsonType = "application/json"

func SendErrorLog(errLog oapiGen.ErrorLog) {
	errLogJson, err := json.Marshal(errLog)
	if err != nil {
		log.Print(err)
		return
	}

	res, err := Client.Post(BaseUrl+"/logs", AppJsonType, bytes.NewBuffer(errLogJson))
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}
}
