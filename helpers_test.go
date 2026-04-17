package easyvapi

import (
	"net/url"
	"testing"
)

// ---------------------------------------------------------------------------
// applyListOptions – URL-Parameter für paginierte API-Anfragen aufbauen
//
// Diese Funktion befüllt ein url.Values-Objekt mit den allgemeinen Parametern,
// die jeder List-Endpunkt benötigt:
//
//   - limit    → Seitengröße (Standard 100)
//   - query    → Feldauswahl im Format {id,name,...} (optional)
//   - ordering → Sortierfeld, Prefix "-" für absteigend (optional)
//   - search   → Volltextsuche (optional)
//
// Services rufen applyListOptions in ihrer listParams-Funktion auf.
// Ist opts.Query == nil, wird defaultQuery verwendet.
// Ist auch defaultQuery == nil (z. B. bei /member-group), wird kein
// query-Parameter gesetzt.
// ---------------------------------------------------------------------------

// TestApplyListOptions_Defaults prüft, dass ohne explizite Optionen der
// Standard-Limit von 100 gesetzt wird und kein query-Parameter erscheint,
// wenn auch defaultQuery nil ist.
func TestApplyListOptions_Defaults(t *testing.T) {
	params := url.Values{}
	applyListOptions(params, ListOptions{}, nil)

	if got := params.Get("limit"); got != "100" {
		t.Errorf("limit = %q, want \"100\"", got)
	}
	if params.Has("query") {
		t.Errorf("query should not be set when both opts.Query and defaultQuery are nil")
	}
}

// TestApplyListOptions_CustomLimit prüft, dass ein explizit gesetztes Limit
// korrekt übernommen wird (hier 25 statt des Standards 100).
func TestApplyListOptions_CustomLimit(t *testing.T) {
	params := url.Values{}
	applyListOptions(params, ListOptions{Limit: 25}, nil)

	if got := params.Get("limit"); got != "25" {
		t.Errorf("limit = %q, want \"25\"", got)
	}
}

// TestApplyListOptions_DefaultQuery prüft, dass defaultQuery verwendet wird,
// wenn opts.Query nil ist. Das ist der Normalfall: Jeder Service definiert
// eine defaultQuery mit allen Modell-Feldern.
func TestApplyListOptions_DefaultQuery(t *testing.T) {
	defaultQ := NewQuery().Fields("id", "name")
	params := url.Values{}
	applyListOptions(params, ListOptions{}, defaultQ)

	want := "{id,name}"
	if got := params.Get("query"); got != want {
		t.Errorf("query = %q, want %q", got, want)
	}
}

// TestApplyListOptions_CustomQueryOverridesDefault prüft, dass eine explizite
// opts.Query die defaultQuery überschreibt. Das ermöglicht Clients, die
// Antwortgröße weiter zu reduzieren ("Ich brauche nur id und email").
func TestApplyListOptions_CustomQueryOverridesDefault(t *testing.T) {
	defaultQ := NewQuery().Fields("id", "name", "email")
	customQ := NewQuery().Fields("id") // nur IDs abrufen
	params := url.Values{}
	applyListOptions(params, ListOptions{Query: customQ}, defaultQ)

	want := "{id}"
	if got := params.Get("query"); got != want {
		t.Errorf("query = %q, want %q", got, want)
	}
}

// TestApplyListOptions_Ordering prüft, dass Ordering korrekt als URL-Parameter
// übergeben wird. Prefix "-" bedeutet absteigende Sortierung (API-Konvention).
func TestApplyListOptions_Ordering(t *testing.T) {
	tests := []struct {
		ordering string
		want     string
	}{
		{"name", "name"},       // aufsteigend
		{"-joinDate", "-joinDate"}, // absteigend
	}
	for _, tc := range tests {
		params := url.Values{}
		applyListOptions(params, ListOptions{Ordering: tc.ordering}, nil)
		if got := params.Get("ordering"); got != tc.want {
			t.Errorf("ordering(%q) = %q, want %q", tc.ordering, got, tc.want)
		}
	}
}

// TestApplyListOptions_Search prüft, dass der search-Parameter gesetzt wird.
// Die API führt eine case-insensitive Suche über alle durchsuchbaren Felder durch.
func TestApplyListOptions_Search(t *testing.T) {
	params := url.Values{}
	applyListOptions(params, ListOptions{Search: "Mustermann"}, nil)

	if got := params.Get("search"); got != "Mustermann" {
		t.Errorf("search = %q, want \"Mustermann\"", got)
	}
}

// TestApplyListOptions_NoOrderingOrSearch prüft, dass optionale Parameter
// (Ordering, Search) nicht als leere Strings gesetzt werden.
// Leere Parameter würden unnötigen Traffic erzeugen und könnten die API
// verwirren.
func TestApplyListOptions_NoOrderingOrSearch(t *testing.T) {
	params := url.Values{}
	applyListOptions(params, ListOptions{}, nil) // kein Ordering, kein Search

	if params.Has("ordering") {
		t.Error("ordering should not be set when empty")
	}
	if params.Has("search") {
		t.Error("search should not be set when empty")
	}
}
