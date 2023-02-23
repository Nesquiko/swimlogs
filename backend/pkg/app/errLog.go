package app

import (
	"github.com/Nesquiko/swimlogs/oapi-generator/oapiGen"
)

func (app *swimLogsApp) RecordError(
	request oapiGen.RecordErrorRequestObject,
) (oapiGen.RecordErrorResponseObject, error) {
	errLog := request.Body
	userDesc := ""

	if errLog.UserDesc != nil {
		userDesc = *errLog.UserDesc
	}

	app.logger.Warnf("error='%s' user description='%s'", errLog.ErrMsg, userDesc)

	return oapiGen.RecordError200Response{}, nil
}
