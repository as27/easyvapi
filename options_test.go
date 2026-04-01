package easyvapi

import (
	"testing"
)

func TestQueryString_Empty(t *testing.T) {
	q := NewQuery()
	if got := q.String(); got != "" {
		t.Errorf("empty Query.String() = %q, want \"\"", got)
	}
}

func TestQueryString_Fields(t *testing.T) {
	q := NewQuery().Fields("id", "name")
	want := "{id,name}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQueryString_Nested(t *testing.T) {
	q := NewQuery().
		Fields("id").
		Nested("contactDetails", "firstName", "familyName")
	want := "{id,contactDetails{firstName,familyName}}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQueryString_Exclude(t *testing.T) {
	q := NewQuery().Fields("id", "name").Exclude("password")
	want := "{id,name,-password}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQueryString_NilQuery(t *testing.T) {
	var q *Query
	if got := q.String(); got != "" {
		t.Errorf("nil Query.String() = %q, want \"\"", got)
	}
}
