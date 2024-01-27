package it

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/Nesquiko/swimlogs/apidef"
	"github.com/Nesquiko/swimlogs/pkg/server"
)

func TestCreateTraining(t *testing.T) {
	req, err := json.Marshal(apidef.CreateTrainingRequest{
		DurationMin: 60,
		Sets: []apidef.NewTrainingSet{
			{
				DistanceMeters: 100,
				Repeat:         1,
				SetOrder:       0,
				StartType:      apidef.None,
				TotalDistance:  100,
			},
		},
		Start:         time.Now(),
		TotalDistance: 100,
	})
	require.NoError(t, err)

	url := TH.ts.URL + "/trainings"
	res, err := http.Post(url, server.ApplicationJSON, bytes.NewBuffer(req))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode)

	response, err := io.ReadAll(res.Body)
	res.Body.Close()
	require.NoError(t, err)
	log.Info().Bytes("response", response).Msg("response from server")
}
