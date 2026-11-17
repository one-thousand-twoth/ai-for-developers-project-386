// Package booking adapts the generated Calendar API to calendar domain services.
package booking

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/one-thousand-twoth/calendar/gen/api"
	"github.com/one-thousand-twoth/calendar/internal/meetingtype"
)

type handler struct {
	meetingTypes *meetingtype.Service
	mu           sync.RWMutex
	bookings     map[string]api.Booking
}

var _ api.Handler = (*handler)(nil)

// NewServer returns the API router. Only legacy demo bookings are process-local.
func NewServer(types *meetingtype.Service) (*api.Server, error) {
	return api.NewServer(&handler{meetingTypes: types, bookings: make(map[string]api.Booking)}, api.WithErrorHandler(handleError))
}

func (h *handler) GetRoot(ctx context.Context) (*api.Health, error) {
	return h.GetHealth(ctx)
}

func (h *handler) GetHealth(context.Context) (*api.Health, error) {
	return &api.Health{Status: api.HealthStatusOk}, nil
}

func (h *handler) GetMeetingTypes(ctx context.Context) (api.GetMeetingTypesRes, error) {
	types, err := h.meetingTypes.List(ctx)
	if errors.Is(err, meetingtype.ErrUnavailable) {
		return &api.GetMeetingTypesServiceUnavailable{
			Code: api.GetMeetingTypesServiceUnavailableCodeTemporarilyUnavailable, Message: "storage temporarily unavailable; retry later",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	result := make(api.GetMeetingTypesOKApplicationJSON, 0, len(types))
	for _, value := range types {
		result = append(result, toAPI(value))
	}
	return &result, nil
}

func (h *handler) CreateMeetingType(ctx context.Context, req *api.MeetingTypeInput) (api.CreateMeetingTypeRes, error) {
	input := meetingtype.MeetingType{Name: req.Name, DurationMinutes: int(req.DurationMinutes)}
	if description, ok := req.Description.Get(); ok {
		input.Description = &description
	}
	created, err := h.meetingTypes.Create(ctx, input)
	if errors.Is(err, meetingtype.ErrInvalid) {
		return &api.ApiError{Code: "bad_request", Message: err.Error()}, nil
	}
	if errors.Is(err, meetingtype.ErrDuplicate) {
		return &api.CreateMeetingTypeConflict{
			Code: api.CreateMeetingTypeConflictCodeDuplicateMeetingType, Message: "meeting type name already exists",
		}, nil
	}
	if errors.Is(err, meetingtype.ErrUnavailable) {
		return &api.CreateMeetingTypeServiceUnavailable{
			Code: api.CreateMeetingTypeServiceUnavailableCodeTemporarilyUnavailable, Message: "storage temporarily unavailable; retry later",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	result := toAPI(created)
	return &result, nil
}

func toAPI(value meetingtype.MeetingType) api.MeetingType {
	result := api.MeetingType{ID: value.ID, Name: value.Name, DurationMinutes: int32(value.DurationMinutes)}
	if value.Description != nil {
		result.Description = api.NewOptString(*value.Description)
	}
	return result
}

func (h *handler) CreateBooking(_ context.Context, req *api.BookingInput) (api.CreateBookingRes, error) {
	if strings.TrimSpace(req.Title) == "" {
		return &api.ApiError{Code: "bad_request", Message: "title must contain a non-whitespace character"}, nil
	}
	if !req.EndsAt.After(req.StartsAt) {
		return &api.ApiError{Code: "bad_request", Message: "endsAt must be strictly after startsAt"}, nil
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	booking := api.Booking{ID: id.String(), Title: req.Title, StartsAt: req.StartsAt.UTC(), EndsAt: req.EndsAt.UTC()}
	h.mu.Lock()
	h.bookings[booking.ID] = booking
	h.mu.Unlock()
	return &booking, nil
}

func (h *handler) GetBooking(_ context.Context, params api.GetBookingParams) (api.GetBookingRes, error) {
	h.mu.RLock()
	booking, ok := h.bookings[params.ID]
	h.mu.RUnlock()
	if !ok {
		return &api.GetBookingNotFound{Code: "not_found", Message: "booking not found"}, nil
	}
	return &booking, nil
}

func handleError(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	status := ogenerrors.ErrorCode(err)
	body := api.ApiError{Code: "internal_error", Message: "internal server error"}
	if status == http.StatusBadRequest {
		body = api.ApiError{Code: "bad_request", Message: "invalid request: " + err.Error()}
	} else if status >= http.StatusInternalServerError {
		log.Printf("request failed: %v", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
