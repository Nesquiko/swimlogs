package server

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/Nesquiko/swimlogs/pkg/api"
	"github.com/Nesquiko/swimlogs/pkg/app"
)

func (s SwimLogsServer) CreateTraining(w http.ResponseWriter, r *http.Request) {
	request, decodeErr := Decode[api.NewTraining](w, r)
	if decodeErr != nil {
		encodeError(w, decodeErr)
		return
	}

	t, err := s.app.CreateTraining(r.Context(), request)
	if err != nil {
		if validationErr, ok := err.(*app.ValidationError); ok {
			apiErr := fromValidationError(validationErr)
			slog.Warn("invalid training", slog.String("error", validationErr.Error()))
			encodeError(w, apiErr)
			return
		}

		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "CreateTraining"),
		)
		encodeError(w, internalServerError())
		return
	}

	encode(w, http.StatusCreated, t)
}

func (s SwimLogsServer) TrainingById(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	t, err := s.app.TrainingById(r.Context(), id)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn("training not found",
			slog.String("error", err.Error()),
			slog.String("id", id.String()))
		encodeError(w, notFound("training", id.String()))
		return
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "TrainingById"),
		)
		encodeError(w, internalServerError())
		return
	}

	encode(w, http.StatusOK, t)
}

func (s SwimLogsServer) DeleteTraining(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	err := s.app.DeleteTraining(r.Context(), id)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn(
			"training not found",
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		encodeError(w, notFoundId("training", id))
		return
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "DeleteTraining"),
		)
		encodeError(w, internalServerError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s SwimLogsServer) DeleteSet(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	err := s.app.DeleteSet(r.Context(), id)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn(
			"set not found",
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		encodeError(w, notFoundId("set", id))
		return
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "DeleteSet"),
		)
		encodeError(w, internalServerError())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s SwimLogsServer) EditTrainingSession(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	request, decodeErr := Decode[api.EditSessionRequest](w, r)
	if decodeErr != nil {
		encodeError(w, decodeErr)
		return
	}

	if app.AllNilFields(request) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	t, err := s.app.EditTrainingSession(r.Context(), id, request)
	if err != nil {
		if validationErr, ok := err.(*app.ValidationError); ok {
			apiErr := fromValidationError(validationErr)
			slog.Warn("invalid training session", slog.String("error", validationErr.Error()))
			encodeError(w, apiErr)
			return
		}

		if errors.Is(err, app.ErrNotFound) {
			slog.Warn("training not found",
				slog.String("error", NotFoundCode),
				slog.String("id", id.String()))
			encodeError(w, notFoundId("training", id))
			return
		}

		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "EditTrainingSession"),
		)
		encodeError(w, internalServerError())
		return
	}

	encode(w, http.StatusOK, t)
}

func (s SwimLogsServer) SummariesPage(
	w http.ResponseWriter,
	r *http.Request,
	params api.SummariesPageParams,
) {
	summaries, total, err := s.app.TrainingSummariesPage(r.Context(), params)
	if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "SummariesPage"),
		)
		encodeError(w, internalServerError())
		return
	}

	pagination := api.Pagination{Page: params.Page, PageSize: len(summaries), Total: total}

	response := api.TrainingSummariesResponse{Summaries: summaries, Pagination: pagination}
	encode(w, http.StatusOK, response)
}

func (s SwimLogsServer) EditSet(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	request, decodeErr := Decode[api.EditSetRequest](w, r)
	if decodeErr != nil {
		encodeError(w, decodeErr)
		return
	}

	if app.AllNilFields(request) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	t, err := s.app.EditSet(r.Context(), id, request)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			slog.Warn(
				"set id not found",
				slog.String("error", err.Error()),
				slog.String("id", id.String()),
			)
			encodeError(w, notFound("set", id.String()))
			return
		}

		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "EditSet"),
		)
		encodeError(w, internalServerError())
		return
	}

	encode(w, http.StatusOK, t)
}

func (s SwimLogsServer) MoveSet(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	request, decodeErr := Decode[api.MoveSetRequest](w, r)
	if decodeErr != nil {
		encodeError(w, decodeErr)
		return
	}
	t, err := s.app.MoveSet(r.Context(), id, request.NewSetOrder)

	if errors.Is(err, app.ErrNotFound) {
		slog.Warn(
			"set not found",
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		encodeError(w, notFoundId("set", id))
		return
	} else if err != nil {
		slog.Error(
			UnexpectedError,
			slog.String("error", err.Error()),
			slog.String("where", "MoveSet"),
		)
		encodeError(w, internalServerError())
		return
	}

	encode(w, http.StatusOK, t)
}

func (s SwimLogsServer) ReplaceSetComponents(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	panic("unimplemented")
}
