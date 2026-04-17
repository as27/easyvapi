package easyvapi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/as27/easyvapi/model"
)

// ---------------------------------------------------------------------------
// EventService – Veranstaltungen verwalten (/event)
//
// Events repräsentieren Vereinsveranstaltungen mit Start-/Enddatum.
// Sie können öffentlich (IsPublic) oder intern sein und einem Kalender
// zugeordnet werden.
//
// Besondere Endpunkte:
//   - Copy:               Veranstaltung duplizieren (GET /event/{id}/copy)
//   - GenerateInvoices:   Rechnungen für alle Teilnehmer erzeugen
//   - InviteGroups:       Mitgliedergruppe zur Veranstaltung einladen
//   - Participations:     CRUD für Teilnahme-Einträge pro Veranstaltung
//
// Filteroptionen (EventListOptions):
//   - StartGte / StartLte – Zeitfenster (ISO 8601)
//   - Calendar             – Kalender-ID
//   - IsPublic             – Öffentlichkeitsstatus (Pointer, damit nil = kein Filter)
// ---------------------------------------------------------------------------

// TestEventService_Get prüft, dass Get die URL /event/{id} aufruft und
// das Event-Objekt korrekt deserialisiert.
func TestEventService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/event/7", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		writeJSON(w, http.StatusOK, model.Event{
			ID:       7,
			Name:     "Jahreshauptversammlung",
			Start:    "2026-03-15T18:00:00Z",
			IsPublic: true,
		})
	})
	c, _ := newTestServer(t, mux)

	e, err := c.Events.Get(context.Background(), 7, nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.ID != 7 {
		t.Errorf("ID = %d, want 7", e.ID)
	}
	if e.Name != "Jahreshauptversammlung" {
		t.Errorf("Name = %q, want \"Jahreshauptversammlung\"", e.Name)
	}
	if !e.IsPublic {
		t.Error("IsPublic = false, want true")
	}
}

// TestEventService_ListAll prüft das Sammeln aller Events über ListAll.
func TestEventService_ListAll(t *testing.T) {
	events := []model.Event{
		{ID: 1, Name: "Sommerfest"},
		{ID: 2, Name: "Weihnachtsfeier"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, pagedOf(events))
	})
	c, _ := newTestServer(t, mux)

	got, err := c.Events.ListAll(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}

// TestEventService_ListOptions_DateRange prüft, dass StartGte und StartLte
// korrekt als URL-Parameter start__gte / start__lte übergeben werden.
// Diese Filter sind wichtig, um Events in einem bestimmten Zeitfenster
// abzufragen (z. B. "alle Events im März 2026").
func TestEventService_ListOptions_DateRange(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("start__gte"); got != "2026-03-01T00:00:00Z" {
			t.Errorf("start__gte = %q, want \"2026-03-01T00:00:00Z\"", got)
		}
		if got := q.Get("start__lte"); got != "2026-03-31T23:59:59Z" {
			t.Errorf("start__lte = %q, want \"2026-03-31T23:59:59Z\"", got)
		}
		writeJSON(w, http.StatusOK, pagedOf([]model.Event{}))
	})
	c, _ := newTestServer(t, mux)

	_, _ = c.Events.ListAll(context.Background(), &EventListOptions{
		StartGte: "2026-03-01T00:00:00Z",
		StartLte: "2026-03-31T23:59:59Z",
	})
}

// TestEventService_ListOptions_IsPublicFilter prüft, dass IsPublic korrekt
// als URL-Parameter isPublic übergeben wird.
//
// IsPublic ist ein *bool-Pointer, damit der Wert "kein Filter" (nil) von
// "explizit false" unterschieden werden kann. Ein nil-Wert sendet keinen
// isPublic-Parameter.
func TestEventService_ListOptions_IsPublicFilter(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("isPublic"); got != "true" {
				t.Errorf("isPublic = %q, want \"true\"", got)
			}
			writeJSON(w, http.StatusOK, pagedOf([]model.Event{}))
		})
		c, _ := newTestServer(t, mux)
		b := true
		_, _ = c.Events.ListAll(context.Background(), &EventListOptions{IsPublic: &b})
	})

	t.Run("false", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("isPublic"); got != "false" {
				t.Errorf("isPublic = %q, want \"false\"", got)
			}
			writeJSON(w, http.StatusOK, pagedOf([]model.Event{}))
		})
		c, _ := newTestServer(t, mux)
		b := false
		_, _ = c.Events.ListAll(context.Background(), &EventListOptions{IsPublic: &b})
	})

	t.Run("nil_no_param", func(t *testing.T) {
		// Ist IsPublic nil, darf kein isPublic-Parameter in der URL erscheinen.
		mux := http.NewServeMux()
		mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Has("isPublic") {
				t.Error("isPublic should not be set when IsPublic is nil")
			}
			writeJSON(w, http.StatusOK, pagedOf([]model.Event{}))
		})
		c, _ := newTestServer(t, mux)
		_, _ = c.Events.ListAll(context.Background(), &EventListOptions{IsPublic: nil})
	})
}

// TestEventService_Create prüft POST /event.
func TestEventService_Create(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/event", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body model.EventCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		writeJSON(w, http.StatusCreated, model.Event{
			ID:    20,
			Name:  body.Name,
			Start: body.Start,
			End:   body.End,
		})
	})
	c, _ := newTestServer(t, mux)

	created, err := c.Events.Create(context.Background(), model.EventCreate{
		Name:  "Kursus",
		Start: "2026-05-10T09:00:00Z",
		End:   "2026-05-10T12:00:00Z",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != 20 {
		t.Errorf("ID = %d, want 20", created.ID)
	}
	if created.Name != "Kursus" {
		t.Errorf("Name = %q, want \"Kursus\"", created.Name)
	}
}

// TestEventService_Update prüft PATCH /event/{id}.
func TestEventService_Update(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/event/20", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		var body model.EventCreate
		_ = json.NewDecoder(r.Body).Decode(&body)
		writeJSON(w, http.StatusOK, model.Event{ID: 20, Name: body.Name})
	})
	c, _ := newTestServer(t, mux)

	updated, err := c.Events.Update(context.Background(), 20, model.EventCreate{Name: "Kursus (aktualisiert)"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "Kursus (aktualisiert)" {
		t.Errorf("Name = %q, want \"Kursus (aktualisiert)\"", updated.Name)
	}
}

// TestEventService_Delete prüft DELETE /event/{id}.
func TestEventService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/event/20", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	c, _ := newTestServer(t, mux)

	if err := c.Events.Delete(context.Background(), 20); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

// TestEventService_Copy prüft den Copy-Endpunkt GET /event/{id}/copy.
// Dieser Endpunkt erstellt eine Kopie des Events und gibt das neue Event zurück.
func TestEventService_Copy(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/event/5/copy", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		writeJSON(w, http.StatusOK, model.Event{ID: 99, Name: "Kopie von Event"})
	})
	c, _ := newTestServer(t, mux)

	copied, err := c.Events.Copy(context.Background(), 5)
	if err != nil {
		t.Fatalf("Copy: %v", err)
	}
	if copied.ID != 99 {
		t.Errorf("ID = %d, want 99", copied.ID)
	}
}

// TestEventService_Participations_CRUD prüft alle Teilnahme-Operationen
// für eine Veranstaltung.
//
// Teilnahmen sind unter /event/{id}/participation verschachtelt.
// Das ermöglicht es, alle Teilnehmer einer Veranstaltung zu verwalten.
func TestEventService_Participations_CRUD(t *testing.T) {
	const eventID = 12

	mux := http.NewServeMux()

	// List
	mux.HandleFunc("/event/12/participation", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, pagedOf([]model.Participation{{ID: 1}, {ID: 2}}))
		case http.MethodPost:
			var body model.ParticipationCreate
			_ = json.NewDecoder(r.Body).Decode(&body)
			writeJSON(w, http.StatusCreated, model.Participation{ID: 10})
		default:
			t.Errorf("unexpected method %s on /event/12/participation", r.Method)
		}
	})

	// Get / Update / Delete einzelner Teilnahme
	mux.HandleFunc("/event/12/participation/1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, model.Participation{ID: 1})
		case http.MethodPatch:
			writeJSON(w, http.StatusOK, model.Participation{ID: 1})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected method %s on /event/12/participation/1", r.Method)
		}
	})

	c, _ := newTestServer(t, mux)
	ctx := context.Background()

	// ListAllParticipations
	parts, err := c.Events.ListAllParticipations(ctx, eventID, nil)
	if err != nil {
		t.Fatalf("ListAllParticipations: %v", err)
	}
	if len(parts) != 2 {
		t.Errorf("len(parts) = %d, want 2", len(parts))
	}

	// GetParticipation
	p, err := c.Events.GetParticipation(ctx, eventID, 1)
	if err != nil {
		t.Fatalf("GetParticipation: %v", err)
	}
	if p.ID != 1 {
		t.Errorf("participation ID = %d, want 1", p.ID)
	}

	// CreateParticipation
	created, err := c.Events.CreateParticipation(ctx, eventID, model.ParticipationCreate{})
	if err != nil {
		t.Fatalf("CreateParticipation: %v", err)
	}
	if created.ID != 10 {
		t.Errorf("created ID = %d, want 10", created.ID)
	}

	// UpdateParticipation
	_, err = c.Events.UpdateParticipation(ctx, eventID, 1, model.ParticipationCreate{})
	if err != nil {
		t.Fatalf("UpdateParticipation: %v", err)
	}

	// DeleteParticipation
	if err := c.Events.DeleteParticipation(ctx, eventID, 1); err != nil {
		t.Fatalf("DeleteParticipation: %v", err)
	}
}

// TestEventListParams_CalendarFilter prüft, dass Calendar-ID korrekt als
// URL-Parameter "calendar" übergeben wird.
func TestEventListParams_CalendarFilter(t *testing.T) {
	params := eventListParams(&EventListOptions{Calendar: 3})

	if got := params.Get("calendar"); got != "3" {
		t.Errorf("calendar = %q, want \"3\"", got)
	}
}
