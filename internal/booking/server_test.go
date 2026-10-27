package booking

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/one-thousand-twoth/calendar/gen/api"
)

func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	server, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func request(server http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	server.ServeHTTP(rec, req)
	return rec
}

const validBody = `{"title":"Call","startsAt":"2026-11-01T15:00:00+08:00","endsAt":"2026-11-01T15:30:00+08:00"}`

func TestHealth(t *testing.T) {
	rec := request(newTestServer(t), http.MethodGet, "/health", "")
	if rec.Code != 200 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 1 || body["status"] != "ok" {
		t.Fatalf("body: %v", body)
	}
}

func TestCreateAndRead(t *testing.T) {
	server := newTestServer(t)
	rec := request(server, http.MethodPost, "/api/bookings", validBody)
	if rec.Code != 201 {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	var booking map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &booking); err != nil {
		t.Fatal(err)
	}
	id, err := uuid.Parse(booking["id"])
	if err != nil || id.Version() != 4 {
		t.Fatalf("invalid UUID: %v", booking)
	}
	if len(booking) != 4 || booking["title"] != "Call" || booking["startsAt"] != "2026-11-01T07:00:00Z" || booking["endsAt"] != "2026-11-01T07:30:00Z" {
		t.Fatalf("unexpected booking: %v", booking)
	}
	got := request(server, http.MethodGet, "/api/bookings/"+booking["id"], "")
	if got.Code != 200 || got.Body.String() != rec.Body.String() {
		t.Fatalf("read %d: %s", got.Code, got.Body)
	}
	checkError(t, request(newTestServer(t), http.MethodGet, "/api/bookings/"+booking["id"], ""), 404, "not_found")
}

func checkError(t *testing.T, rec *httptest.ResponseRecorder, status int, code string) {
	t.Helper()
	if rec.Code != status {
		t.Fatalf("want %d, got %d: %s", status, rec.Code, rec.Body)
	}
	if !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("content type: %v", rec.Header())
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if len(body) != 2 || body["code"] != code || body["message"] == "" {
		t.Fatalf("error body: %v", body)
	}
}

func TestInvalidRequests(t *testing.T) {
	server := newTestServer(t)
	cases := map[string]string{
		"empty title":           strings.Replace(validBody, `"Call"`, `""`, 1),
		"whitespace title":      strings.Replace(validBody, `"Call"`, `" \t\n\u2003"`, 1),
		"missing title":         strings.Replace(validBody, `"title":"Call",`, "", 1),
		"missing date":          `{"title":"Call","startsAt":"2026-11-01T15:00:00Z"}`,
		"invalid date":          strings.Replace(validBody, "2026-11-01T15:00:00+08:00", "yesterday", 1),
		"invalid calendar date": strings.Replace(validBody, "2026-11-01T15:00:00+08:00", "2026-02-30T15:00:00Z", 1),
		"no timezone":           strings.Replace(validBody, "2026-11-01T15:00:00+08:00", "2026-11-01T15:00:00", 1),
		"equal dates":           strings.Replace(validBody, "15:30:00", "15:00:00", 1),
		"reversed dates":        strings.Replace(validBody, "15:30:00", "14:00:00", 1),
		"malformed JSON":        `{"title":`,
		"trailing JSON":         validBody + `{}`,
		"null body":             `null`,
		"wrong field type":      strings.Replace(validBody, `"Call"`, `123`, 1),
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			checkError(t, request(server, http.MethodPost, "/api/bookings", body), 400, "bad_request")
		})
	}
	checkError(t, request(server, http.MethodGet, "/api/bookings/unknown", ""), 404, "not_found")
}

func TestConcurrentBookings(t *testing.T) {
	server := newTestServer(t)
	const count = 32
	ids := make(chan string, count)
	var wg sync.WaitGroup
	for range count {
		wg.Go(func() {
			rec := request(server, http.MethodPost, "/api/bookings", validBody)
			if rec.Code != 201 {
				t.Errorf("create %d: %s", rec.Code, rec.Body)
				return
			}
			var booking api.Booking
			if err := json.Unmarshal(rec.Body.Bytes(), &booking); err != nil {
				t.Error(err)
				return
			}
			ids <- booking.ID
			got := request(server, http.MethodGet, "/api/bookings/"+booking.ID, "")
			if got.Code != 200 || got.Body.String() != rec.Body.String() {
				t.Errorf("read %d: %s", got.Code, got.Body)
			}
		})
	}
	wg.Wait()
	close(ids)
	seen := make(map[string]bool)
	for id := range ids {
		if seen[id] {
			t.Errorf("duplicate id %s", id)
		}
		seen[id] = true
	}
	if len(seen) != count {
		t.Fatalf("stored %d, want %d", len(seen), count)
	}
}
