package easyvapi

import (
	"errors"
	"testing"
)

// ---------------------------------------------------------------------------
// Iterator[T] – lazy, page-by-page Iteration über paginierte API-Antworten
//
// Der Iterator holt Seiten erst dann, wenn der aktuelle Puffer erschöpft ist
// ("lazy fetching"). So kann man auch sehr große Ergebnismengen durchlaufen,
// ohne alles auf einmal in den Speicher zu laden.
//
// Grundlegendes Nutzungsmuster:
//
//	iter := client.Members.List(ctx, nil)
//	for iter.Next() {
//	    item := iter.Value()
//	}
//	if err := iter.Err(); err != nil { … }
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Hilfsfunktion: buildIterator
//
// Erzeugt einen Iterator, der seine Seiten aus einem In-Memory-Slice von Pages
// bezieht. So können wir den Iterator isoliert testen, ohne HTTP-Server.
//
// pages ist ein Slice von Slices – jedes innere Slice ist eine "API-Seite".
// ---------------------------------------------------------------------------
func buildIterator(pages [][]int) *Iterator[int] {
	pageIdx := 0
	return newIterator("start", func(_ string) ([]int, *string, error) {
		if pageIdx >= len(pages) {
			return nil, nil, nil
		}
		items := pages[pageIdx]
		pageIdx++
		var next *string
		if pageIdx < len(pages) {
			url := "page-" + string(rune('0'+pageIdx))
			next = &url
		}
		return items, next, nil
	})
}

// TestIterator_EmptyResult prüft, dass ein leerer Iterator sofort terminiert
// und iter.Err() nil ist.
func TestIterator_EmptyResult(t *testing.T) {
	// Szenario: Die API gibt eine leere Seite zurück.
	iter := buildIterator([][]int{{}})

	if iter.Next() {
		t.Error("Next() = true on empty result, want false")
	}
	if err := iter.Err(); err != nil {
		t.Errorf("Err() = %v, want nil", err)
	}
}

// TestIterator_SinglePage prüft die normale Iteration über eine einzelne Seite.
// Value() muss die Elemente in der richtigen Reihenfolge liefern.
func TestIterator_SinglePage(t *testing.T) {
	// Szenario: Genau eine Seite mit drei Elementen.
	iter := buildIterator([][]int{{10, 20, 30}})

	want := []int{10, 20, 30}
	for i, expected := range want {
		if !iter.Next() {
			t.Fatalf("Next() = false at index %d, want true", i)
		}
		if got := iter.Value(); got != expected {
			t.Errorf("Value() = %d, want %d", got, expected)
		}
	}
	// Nach dem letzten Element muss Next() false zurückgeben.
	if iter.Next() {
		t.Error("Next() = true after last element, want false")
	}
	if err := iter.Err(); err != nil {
		t.Errorf("Err() = %v, want nil", err)
	}
}

// TestIterator_MultiPage prüft, dass der Iterator beim Seitenübergang
// automatisch die nächste Seite lädt ("lazy pagination").
//
// Das ist der Kernmechanismus: Sobald der aktuelle Puffer erschöpft ist
// und nextURL != nil, wird fetchFunc erneut aufgerufen.
func TestIterator_MultiPage(t *testing.T) {
	// Szenario: Drei Seiten à zwei Elemente = sechs Elemente insgesamt.
	iter := buildIterator([][]int{
		{1, 2},
		{3, 4},
		{5, 6},
	})

	var got []int
	for iter.Next() {
		got = append(got, iter.Value())
	}
	if err := iter.Err(); err != nil {
		t.Fatalf("Err() = %v, want nil", err)
	}
	want := []int{1, 2, 3, 4, 5, 6}
	if len(got) != len(want) {
		t.Fatalf("collected %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("item[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

// TestIterator_FetchError prüft, dass ein Fehler beim Laden einer Seite
// korrekt in iter.Err() landet und Next() danach false zurückgibt.
//
// Das verhindert, dass ein HTTP-Fehler mitten in der Iteration stillschweigend
// ignoriert wird.
func TestIterator_FetchError(t *testing.T) {
	// Szenario: Erste Seite liefert Daten, zweite Seite schlägt fehl.
	fetchErr := errors.New("network error")
	callCount := 0
	iter := newIterator("start", func(_ string) ([]int, *string, error) {
		callCount++
		if callCount == 1 {
			// Erste Seite: zwei Elemente, nächste Seite signalisieren.
			next := "page-2"
			return []int{1, 2}, &next, nil
		}
		// Zweite Seite: Fehler
		return nil, nil, fetchErr
	})

	// Erstes und zweites Element sind noch erreichbar.
	if !iter.Next() {
		t.Fatal("Next() = false, want true for first element")
	}
	if !iter.Next() {
		t.Fatal("Next() = false, want true for second element")
	}
	// Beim dritten Aufruf muss der Fetch fehlschlagen → Next() = false.
	if iter.Next() {
		t.Error("Next() = true after fetch error, want false")
	}
	if !errors.Is(iter.Err(), fetchErr) {
		t.Errorf("Err() = %v, want %v", iter.Err(), fetchErr)
	}
}

// TestIterator_NextAfterError prüft, dass weitere Next()-Aufrufe nach einem
// Fehler sicher sind und weiterhin false zurückgeben.
func TestIterator_NextAfterError(t *testing.T) {
	fetchErr := errors.New("oops")
	iter := newIterator("start", func(_ string) ([]int, *string, error) {
		return nil, nil, fetchErr
	})

	if iter.Next() {
		t.Error("first Next() = true, want false")
	}
	// Erneuter Aufruf darf nicht paniken oder einen anderen Fehler liefern.
	if iter.Next() {
		t.Error("second Next() = true after error, want false")
	}
	if !errors.Is(iter.Err(), fetchErr) {
		t.Errorf("Err() = %v, want %v", iter.Err(), fetchErr)
	}
}
