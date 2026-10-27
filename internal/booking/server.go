// Package booking implements the generated Calendar API with process-local storage.
package booking

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/one-thousand-twoth/calendar/gen/api"
)

type handler struct {
	mu       sync.RWMutex
	bookings map[string]api.Booking
}

var _ api.Handler = (*handler)(nil)

// NewServer returns an isolated API server. Its data is lost on restart.
func NewServer() (*api.Server, error) {
	return api.NewServer(&handler{bookings: make(map[string]api.Booking)}, api.WithErrorHandler(handleError))
}

func (h *handler) GetHealth(context.Context) (*api.Health, error) {
	return &api.Health{Status: api.HealthStatusOk}, nil
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
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
