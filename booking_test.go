package easyvapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/as27/easyvapi/model"
)

// ---------------------------------------------------------------------------
// BookingService – Finanzbuchungen verwalten (/booking)
//
// Buchungen repräsentieren einzelne Finanztransaktionen im Verein.
// Jede Buchung hat mindestens: Betrag (Amount), Datum (Date) und ein
// Buchungskonto (BillingAccount, referenziert per ID).
//
// Besondere Endpunkte:
//   - BulkCreate: Mehrere Buchungen in einem HTTP-Request anlegen
//   - BulkUpdate: Mehrere Buchungen in einem HTTP-Request ändern (PATCH)
//
// Filteroptionen (BookingListOptions):
//   - BankAccount   – nach Bankkonto filtern
//   - BillingAccount – nach Buchungskonto filtern
//   - DateGt / DateLt – Datumsbereich (größer/kleiner als, YYYY-MM-DD)
// ---------------------------------------------------------------------------

// TestBookingService_Get prüft, dass Get die korrekte URL /booking/{id}
// aufruft und das Ergebnis deserialisiert.
func TestBookingService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking/10", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		writeJSON(w, http.StatusOK, model.Booking{
			ID:     10,
			Amount: 99.50,
			Date:   "2026-01-15",
		})
	})
	c, _ := newTestServer(t, mux)

	b, err := c.Bookings.Get(context.Background(), 10, nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if b.ID != 10 {
		t.Errorf("ID = %d, want 10", b.ID)
	}
	if b.Date != "2026-01-15" {
		t.Errorf("Date = %q, want \"2026-01-15\"", b.Date)
	}
}

// TestBookingService_ListAll prüft, dass ListAll alle Buchungen einer
// paginierten Antwort korrekt sammelt.
func TestBookingService_ListAll(t *testing.T) {
	bookings := []model.Booking{
		{ID: 1, Amount: 50.00, Date: "2026-01-01"},
		{ID: 2, Amount: 75.00, Date: "2026-01-02"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, pagedOf(bookings))
	})
	c, _ := newTestServer(t, mux)

	got, err := c.Bookings.ListAll(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].ID != 1 || got[1].ID != 2 {
		t.Errorf("IDs = %d, %d – want 1, 2", got[0].ID, got[1].ID)
	}
}

// TestBookingService_ListOptions_DateFilter prüft, dass DateGt und DateLt
// korrekt als URL-Parameter date__gt / date__lt übergeben werden.
// Diese Filter sind wichtig, um Buchungen für einen bestimmten Monat abzufragen.
func TestBookingService_ListOptions_DateFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("date__gt"); got != "2026-01-01" {
			t.Errorf("date__gt = %q, want \"2026-01-01\"", got)
		}
		if got := q.Get("date__lt"); got != "2026-01-31" {
			t.Errorf("date__lt = %q, want \"2026-01-31\"", got)
		}
		writeJSON(w, http.StatusOK, pagedOf([]model.Booking{}))
	})
	c, _ := newTestServer(t, mux)

	_, err := c.Bookings.ListAll(context.Background(), &BookingListOptions{
		DateGt: "2026-01-01",
		DateLt: "2026-01-31",
	})
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
}

// TestBookingService_ListOptions_BillingAccountFilter prüft, dass
// BillingAccount als URL-Parameter billingAccount übergeben wird.
func TestBookingService_ListOptions_BillingAccountFilter(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("billingAccount"); got != "42" {
			t.Errorf("billingAccount = %q, want \"42\"", got)
		}
		writeJSON(w, http.StatusOK, pagedOf([]model.Booking{}))
	})
	c, _ := newTestServer(t, mux)

	_, _ = c.Bookings.ListAll(context.Background(), &BookingListOptions{
		BillingAccount: 42,
	})
}

// TestBookingService_Create prüft, dass Create einen POST-Request an /booking
// sendet und die Antwort korrekt zurückgibt.
//
// Pflichtfelder einer Buchung: Amount, BillingAccount, Date.
func TestBookingService_Create(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var body model.BookingCreate
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":     55,
			"amount": body.Amount,
			"date":   body.Date,
		})
	})
	c, _ := newTestServer(t, mux)

	created, err := c.Bookings.Create(context.Background(), model.BookingCreate{
		Amount:         120.00,
		BillingAccount: 7,
		Date:           "2026-03-31",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != 55 {
		t.Errorf("ID = %d, want 55", created.ID)
	}
	if created.Date != "2026-03-31" {
		t.Errorf("Date = %q, want \"2026-03-31\"", created.Date)
	}
}

// TestBookingService_Update prüft, dass Update einen PATCH-Request an
// /booking/{id} schickt und das aktualisierte Objekt zurückgibt.
//
// PATCH statt PUT: Es werden nur die gesendeten Felder geändert, alle
// anderen Felder bleiben unverändert.
func TestBookingService_Update(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking/5", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		var body model.BookingCreate
		_ = json.NewDecoder(r.Body).Decode(&body)
		writeJSON(w, http.StatusOK, map[string]any{
			"id":     5,
			"amount": body.Amount,
			"date":   body.Date,
		})
	})
	c, _ := newTestServer(t, mux)

	updated, err := c.Bookings.Update(context.Background(), 5, model.BookingCreate{
		Amount: 200.00,
		Date:   "2026-04-01",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Date != "2026-04-01" {
		t.Errorf("Date = %q, want \"2026-04-01\"", updated.Date)
	}
}

// TestBookingService_Delete prüft, dass Delete einen DELETE-Request an
// /booking/{id} sendet und bei HTTP 204 keinen Fehler zurückgibt.
func TestBookingService_Delete(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking/3", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s, want DELETE", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	c, _ := newTestServer(t, mux)

	if err := c.Bookings.Delete(context.Background(), 3); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

// TestBookingService_BulkCreate prüft den Bulk-Create-Endpunkt POST /booking/bulk.
//
// BulkCreate ist effizienter als einzelne Create-Aufrufe: alle Buchungen
// werden in einem einzigen HTTP-Request an den Server geschickt.
// Rückgabe: Slice mit den angelegten Buchungen inklusive server-seitig
// vergebener IDs.
func TestBookingService_BulkCreate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking/bulk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		var bodies []model.BookingCreate
		if err := json.NewDecoder(r.Body).Decode(&bodies); err != nil {
			t.Errorf("decode body: %v", err)
		}
		// Server vergibt IDs.
		result := make([]map[string]any, len(bodies))
		for i, b := range bodies {
			result[i] = map[string]any{"id": i + 100, "amount": b.Amount, "date": b.Date}
		}
		writeJSON(w, http.StatusCreated, result)
	})
	c, _ := newTestServer(t, mux)

	bookings := []model.BookingCreate{
		{Amount: 10.00, BillingAccount: 1, Date: "2026-01-01"},
		{Amount: 20.00, BillingAccount: 1, Date: "2026-01-02"},
	}
	created, err := c.Bookings.BulkCreate(context.Background(), bookings)
	if err != nil {
		t.Fatalf("BulkCreate: %v", err)
	}
	if len(created) != 2 {
		t.Fatalf("len = %d, want 2", len(created))
	}
	if created[0].ID != 100 || created[1].ID != 101 {
		t.Errorf("IDs = %d, %d – want 100, 101", created[0].ID, created[1].ID)
	}
}

// TestBookingService_BulkUpdate prüft den Bulk-Update-Endpunkt PATCH /booking/bulk.
//
// BulkUpdate aktualisiert mehrere bestehende Buchungen in einem Request.
// Die Buchungsobjekte müssen IDs enthalten, damit der Server sie zuordnen kann.
func TestBookingService_BulkUpdate(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking/bulk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("method = %s, want PATCH", r.Method)
		}
		var bodies []model.Booking
		_ = json.NewDecoder(r.Body).Decode(&bodies)
		writeJSON(w, http.StatusOK, bodies) // Server gibt dieselben Objekte zurück
	})
	c, _ := newTestServer(t, mux)

	toUpdate := []model.Booking{
		{ID: 10, Amount: 15.00, Date: "2026-02-01"},
		{ID: 11, Amount: 25.00, Date: "2026-02-02"},
	}
	updated, err := c.Bookings.BulkUpdate(context.Background(), toUpdate)
	if err != nil {
		t.Fatalf("BulkUpdate: %v", err)
	}
	if len(updated) != 2 {
		t.Fatalf("len = %d, want 2", len(updated))
	}
	if updated[0].ID != 10 || updated[1].ID != 11 {
		t.Errorf("IDs = %d, %d – want 10, 11", updated[0].ID, updated[1].ID)
	}
}

// TestBookingListParams_NoOptions prüft, dass bookingListParams mit nil-Optionen
// die Default-Parameter (limit=100, default query) setzt und keine
// optionalen Filter-Parameter enthält.
func TestBookingListParams_NoOptions(t *testing.T) {
	params := bookingListParams(nil)

	if params.Get("limit") != "100" {
		t.Errorf("limit = %q, want \"100\"", params.Get("limit"))
	}
	// Kein BankAccount- oder BillingAccount-Filter gesetzt
	for _, key := range []string{"bankAccount", "billingAccount", "date__gt", "date__lt"} {
		if params.Has(key) {
			t.Errorf("unexpected param %q in default params", key)
		}
	}
}

// TestBookingListParams_AllFilters prüft, dass alle Filteroptionen korrekt
// in URL-Parameter übersetzt werden.
func TestBookingListParams_AllFilters(t *testing.T) {
	opts := &BookingListOptions{
		BankAccount:    5,
		BillingAccount: 7,
		DateGt:         "2026-01-01",
		DateLt:         "2026-12-31",
	}
	params := bookingListParams(opts)

	checks := []struct{ key, want string }{
		{"bankAccount", "5"},
		{"billingAccount", "7"},
		{"date__gt", "2026-01-01"},
		{"date__lt", "2026-12-31"},
	}
	for _, c := range checks {
		if got := params.Get(c.key); got != c.want {
			t.Errorf("%s = %q, want %q", c.key, got, c.want)
		}
	}
}

// Stelle sicher, dass url.Values importiert wird (wird in den Params-Tests genutzt).
var _ = url.Values{}
