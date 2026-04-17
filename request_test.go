package easyvapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/as27/easyvapi/model"
)

// ---------------------------------------------------------------------------
// buildURL – vollständige Request-URL konstruieren
//
// buildURL fügt Pfad und URL-Parameter an die Basis-URL des Clients an.
// Sonderfälle: führende "/" am Pfad ist optional; Basis-URL darf trailing "/"
// haben.
// ---------------------------------------------------------------------------

func TestBuildURL_Simple(t *testing.T) {
	c := New("token", WithBaseURL("https://api.example.com/v1"))

	got := c.buildURL("/member", url.Values{})
	want := "https://api.example.com/v1/member"
	if got != want {
		t.Errorf("buildURL = %q, want %q", got, want)
	}
}

// TestBuildURL_WithParams prüft, dass URL-Parameter korrekt angehängt werden.
// Die Reihenfolge der Parameter kann variieren (url.Values.Encode() sortiert
// alphabetisch), daher parsen wir die URL zum Vergleich.
func TestBuildURL_WithParams(t *testing.T) {
	c := New("token", WithBaseURL("https://api.example.com/v1"))
	params := url.Values{"limit": {"100"}, "ordering": {"name"}}

	rawURL := c.buildURL("/member", params)

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse URL: %v", err)
	}
	q := parsed.Query()
	if q.Get("limit") != "100" {
		t.Errorf("limit = %q, want \"100\"", q.Get("limit"))
	}
	if q.Get("ordering") != "name" {
		t.Errorf("ordering = %q, want \"name\"", q.Get("ordering"))
	}
}

// TestBuildURL_TrailingSlashInBase prüft, dass eine trailing "/" in der
// Basis-URL nicht zu einem doppelten Slash in der fertigen URL führt.
func TestBuildURL_TrailingSlashInBase(t *testing.T) {
	c := New("token", WithBaseURL("https://api.example.com/v1/"))

	got := c.buildURL("/member", url.Values{})
	want := "https://api.example.com/v1/member"
	if got != want {
		t.Errorf("buildURL = %q, want %q", got, want)
	}
}

// TestBuildURL_PathWithoutLeadingSlash prüft, dass fehlende "/" am Pfad
// automatisch ergänzt wird.
func TestBuildURL_PathWithoutLeadingSlash(t *testing.T) {
	c := New("token", WithBaseURL("https://api.example.com/v1"))

	got := c.buildURL("member", url.Values{}) // kein führendes /
	want := "https://api.example.com/v1/member"
	if got != want {
		t.Errorf("buildURL = %q, want %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// HTTP-Fehlerbehandlung
//
// Der Client wandelt alle nicht-2xx-Antworten in *APIError um. Dabei liest er
// den Response-Body aus und versucht ihn als JSON zu parsen:
//   - { "detail": "..." }  → APIError.Detail
//   - { "message": "..." } → APIError.Message
// Ist kein JSON-Feld vorhanden, wird der Rohtext als Detail verwendet.
// ---------------------------------------------------------------------------

// TestAPIError_DetailField prüft, dass das "detail"-Feld aus dem JSON-Body
// korrekt in APIError.Detail landet.
func TestAPIError_DetailField(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/member/1", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"detail": "joinDate is required",
		})
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Members.Get(context.Background(), 1, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Detail != "joinDate is required" {
		t.Errorf("Detail = %q, want \"joinDate is required\"", apiErr.Detail)
	}
}

// TestAPIError_RawBodyFallback prüft, dass bei nicht-JSON-Antworten der
// Rohtext als Detail übernommen wird. Das passiert z. B. bei Proxy-Fehlern
// oder unerwarteten Server-Responses.
func TestAPIError_RawBodyFallback(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/member/2", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Server Error"))
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Members.Get(context.Background(), 2, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
	if apiErr.Detail == "" {
		t.Error("Detail should contain raw body text, got empty string")
	}
}

// TestAPIError_MessageFallback prüft, dass APIError.Message mit dem
// HTTP-Statustext befüllt wird, wenn die API kein eigenes message-Feld sendet.
// So ist der Fehler immer lesbar, auch ohne spezifische API-Fehlermeldung.
func TestAPIError_MessageFallback(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/member/3", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusForbidden, map[string]string{
			"detail": "insufficient permissions",
		})
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Members.Get(context.Background(), 3, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	// Message soll mit http.StatusText befüllt sein ("Forbidden")
	if apiErr.Message == "" {
		t.Error("Message should not be empty for 403 response")
	}
}

// ---------------------------------------------------------------------------
// Rate-Limit-Handling (429 Too Many Requests)
//
// Bei HTTP 429 gibt der Client einen *RateLimitError zurück, der die
// empfohlene Wartezeit enthält. Die Wartezeit stammt aus dem Retry-After-Header
// (Sekunden als Integer) oder fällt auf 60 Sekunden zurück.
// ---------------------------------------------------------------------------

// TestRateLimitError_DefaultRetryAfter prüft das Fallback auf 60 Sekunden,
// wenn kein Retry-After-Header gesetzt ist.
func TestRateLimitError_DefaultRetryAfter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/member/1", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Members.Get(context.Background(), 1, nil)
	var rlErr *RateLimitError
	if !errors.As(err, &rlErr) {
		t.Fatalf("error type = %T, want *RateLimitError", err)
	}
	if rlErr.RetryAfter != 60*time.Second {
		t.Errorf("RetryAfter = %v, want 60s", rlErr.RetryAfter)
	}
}

// TestRateLimitError_RetryAfterHeader prüft, dass der Retry-After-Header
// korrekt in eine time.Duration umgewandelt wird.
func TestRateLimitError_RetryAfterHeader(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/member/1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "30")
		w.WriteHeader(http.StatusTooManyRequests)
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Members.Get(context.Background(), 1, nil)
	var rlErr *RateLimitError
	if !errors.As(err, &rlErr) {
		t.Fatalf("error type = %T, want *RateLimitError", err)
	}
	if rlErr.RetryAfter != 30*time.Second {
		t.Errorf("RetryAfter = %v, want 30s", rlErr.RetryAfter)
	}
}

// ---------------------------------------------------------------------------
// Token-Refresh
//
// Wenn der Server den Header "tokenRefreshNeeded: true" setzt, ruft der Client
// automatisch GET /refresh-token auf, speichert das neue Token und wiederholt
// den ursprünglichen Request – transparent für den Aufrufer.
// Die Callback-Funktion (WithTokenRefreshCallback) wird mit dem neuen Token
// aufgerufen, damit der Aufrufer das Token persistieren kann.
// ---------------------------------------------------------------------------

// TestTokenRefresh_AutoRetry prüft, dass nach einem tokenRefreshNeeded-Header
// der Request automatisch wiederholt wird und das Ergebnis korrekt zurückkommt.
func TestTokenRefresh_AutoRetry(t *testing.T) {
	firstCall := true
	mux := http.NewServeMux()

	// /refresh-token liefert ein neues Token.
	mux.HandleFunc("/refresh-token", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"token": "new-token"})
	})

	// /member/1 signalisiert beim ersten Aufruf Refresh, beim zweiten liefert
	// es die eigentliche Antwort.
	mux.HandleFunc("/member/1", func(w http.ResponseWriter, r *http.Request) {
		if firstCall {
			firstCall = false
			w.Header().Set("tokenRefreshNeeded", "true")
			// Wir müssen trotzdem einen gültigen Body liefern, da der Client
			// die Antwort verwirft und den Request nach dem Refresh wiederholt.
			writeJSON(w, http.StatusOK, model.Member{ID: 1})
			return
		}
		writeJSON(w, http.StatusOK, model.Member{ID: 1, MembershipNumber: "M-Refreshed"})
	})

	var capturedToken string
	c, _ := newTestServer(t, mux)
	c.onTokenRefresh = func(newToken string) {
		capturedToken = newToken
	}

	m, err := c.Members.Get(context.Background(), 1, nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if m.MembershipNumber != "M-Refreshed" {
		t.Errorf("MembershipNumber = %q, want \"M-Refreshed\"", m.MembershipNumber)
	}
	if capturedToken != "new-token" {
		t.Errorf("capturedToken = %q, want \"new-token\"", capturedToken)
	}
	if c.token != "new-token" {
		t.Errorf("client.token = %q, want \"new-token\"", c.token)
	}
}

// ---------------------------------------------------------------------------
// Content-Type-Header
//
// Bei POST/PATCH-Anfragen mit einem Body muss der Client
// Content-Type: application/json setzen. Bei GET/DELETE (kein Body)
// darf kein Content-Type-Header gesendet werden.
// ---------------------------------------------------------------------------

// TestContentTypeHeader_SetOnCreate prüft, dass Content-Type korrekt
// gesetzt wird, wenn ein Request-Body vorhanden ist (POST).
func TestContentTypeHeader_SetOnCreate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/member", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Content-Type = %q, want \"application/json\"", ct)
		}
		writeJSON(w, http.StatusCreated, model.Member{ID: 1})
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Members.Create(context.Background(), model.MemberCreate{JoinDate: "2024-01-01"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
}

// TestRequestBody_SentCorrectly prüft, dass der Request-Body korrekt als JSON
// serialisiert und vom Server empfangen wird. Das stellt sicher, dass keine
// Felder beim Marshalling verloren gehen.
func TestRequestBody_SentCorrectly(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/member", func(w http.ResponseWriter, r *http.Request) {
		var body model.MemberCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		if body.JoinDate != "2025-06-01" {
			t.Errorf("body.JoinDate = %q, want \"2025-06-01\"", body.JoinDate)
		}
		if body.PaymentAmount != 42.50 {
			t.Errorf("body.PaymentAmount = %v, want 42.50", body.PaymentAmount)
		}
		writeJSON(w, http.StatusCreated, model.Member{ID: 1})
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Members.Create(context.Background(), model.MemberCreate{
		JoinDate:      "2025-06-01",
		PaymentAmount: 42.50,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
}
