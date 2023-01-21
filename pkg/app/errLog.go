package app

import (
	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

func (app *swimLogsApp) RecordError(
	request oapiGen.RecordErrorRequestObject,
) (oapiGen.RecordErrorResponseObject, error) {
	errLog := request.Body
	app.logger.Warnf("error='%s' user description='%s'", errLog.ErrMsg, *errLog.UserDesc)

	return oapiGen.RecordError200Response{}, nil
}
