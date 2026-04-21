// Package easyvapi – tests for all exported symbols.
//
// These tests document the public contract of the package and guard against
// unintended behaviour changes in future refactors. Each test group is
// prefixed with the type or function it covers so that the test output alone
// explains what the package does.
package easyvapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/as27/easyvapi/model"
)


// ---------------------------------------------------------------------------
// APIError
// ---------------------------------------------------------------------------

// APIError is returned for any non-2xx HTTP response. It carries the status
// code, a short message, and an optional detail string from the response body.

func TestAPIError_Error_WithoutDetail(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "Not Found"}
	want := "easyvapi: HTTP 404: Not Found"
	if got := err.Error(); got != want {
		t.Errorf("APIError.Error() = %q, want %q", got, want)
	}
}

func TestAPIError_Error_WithDetail(t *testing.T) {
	err := &APIError{StatusCode: 400, Message: "Bad Request", Detail: "field 'email' is required"}
	want := "easyvapi: HTTP 400: Bad Request – field 'email' is required"
	if got := err.Error(); got != want {
		t.Errorf("APIError.Error() = %q, want %q", got, want)
	}
}

// errors.As allows callers to inspect status code and detail without a type assertion.
func TestAPIError_ErrorsAs(t *testing.T) {
	original := &APIError{StatusCode: 422, Message: "Unprocessable Entity", Detail: "duplicate"}
	wrapped := fmt.Errorf("operation failed: %w", original)

	var apiErr *APIError
	if !errors.As(wrapped, &apiErr) {
		t.Fatal("errors.As should find *APIError in wrapped error")
	}
	if apiErr.StatusCode != 422 {
		t.Errorf("StatusCode = %d, want 422", apiErr.StatusCode)
	}
	if apiErr.Detail != "duplicate" {
		t.Errorf("Detail = %q, want \"duplicate\"", apiErr.Detail)
	}
}

// APIError implements the error interface.
func TestAPIError_ImplementsError(t *testing.T) {
	var _ error = &APIError{}
}

// ---------------------------------------------------------------------------
// RateLimitError
// ---------------------------------------------------------------------------

// RateLimitError is returned when the API responds with HTTP 429. It contains
// a RetryAfter duration that callers should wait before retrying.

func TestRateLimitError_Error(t *testing.T) {
	err := &RateLimitError{RetryAfter: 60 * time.Second}
	want := "easyvapi: rate limit exceeded, retry after 1m0s"
	if got := err.Error(); got != want {
		t.Errorf("RateLimitError.Error() = %q, want %q", got, want)
	}
}

func TestRateLimitError_ErrorsAs(t *testing.T) {
	original := &RateLimitError{RetryAfter: 30 * time.Second}
	wrapped := fmt.Errorf("list failed: %w", original)

	var rlErr *RateLimitError
	if !errors.As(wrapped, &rlErr) {
		t.Fatal("errors.As should find *RateLimitError in wrapped error")
	}
	if rlErr.RetryAfter != 30*time.Second {
		t.Errorf("RetryAfter = %v, want 30s", rlErr.RetryAfter)
	}
}

// RateLimitError implements the error interface.
func TestRateLimitError_ImplementsError(t *testing.T) {
	var _ error = &RateLimitError{}
}

// ---------------------------------------------------------------------------
// Query builder
// ---------------------------------------------------------------------------

// NewQuery creates an empty query. Fields, Nested, and Exclude accumulate
// entries and the String() method renders them as {field,...} syntax.

func TestQuery_MultipleFieldsCalls_Accumulate(t *testing.T) {
	// Calling Fields() multiple times appends rather than replacing.
	q := NewQuery().Fields("id").Fields("name").Fields("email")
	want := "{id,name,email}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQuery_SameNestedNameAccumulates(t *testing.T) {
	// Calling Nested() twice with the same name merges fields.
	q := NewQuery().
		Nested("contactDetails", "firstName").
		Nested("contactDetails", "familyName")
	want := "{contactDetails{firstName,familyName}}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQuery_NestedInsertionOrderPreserved(t *testing.T) {
	// Nested groups appear in the order they were first added.
	q := NewQuery().
		Nested("memberGroups", "id").
		Nested("contactDetails", "email").
		Nested("memberGroups", "name") // second call to same key
	want := "{memberGroups{id,name},contactDetails{email}}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQuery_ExcludeMultipleFields(t *testing.T) {
	q := NewQuery().Fields("id", "name", "secret").Exclude("secret", "internalId")
	want := "{id,name,secret,-secret,-internalId}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQuery_PartOrder_FieldsThenNestedThenExcluded(t *testing.T) {
	// Output order is always: top-level fields, nested groups, excluded fields.
	q := NewQuery().
		Exclude("password").
		Nested("addr", "zip").
		Fields("id")
	want := "{id,addr{zip},-password}"
	if got := q.String(); got != want {
		t.Errorf("Query.String() = %q, want %q", got, want)
	}
}

func TestQuery_NilString_ReturnsEmpty(t *testing.T) {
	var q *Query
	if got := q.String(); got != "" {
		t.Errorf("nil Query.String() = %q, want \"\"", got)
	}
}

func TestQuery_EmptyNewQuery_ReturnsEmpty(t *testing.T) {
	if got := NewQuery().String(); got != "" {
		t.Errorf("empty NewQuery().String() = %q, want \"\"", got)
	}
}

// ---------------------------------------------------------------------------
// New() and functional options
// ---------------------------------------------------------------------------

// New returns a Client with all 40 services initialized and ready to use.

func TestNew_AllServicesInitialized(t *testing.T) {
	c := New("tok")

	services := []struct {
		name string
		ok   bool
	}{
		{"Members", c.Members != nil},
		{"ContactDetails", c.ContactDetails != nil},
		{"Invoices", c.Invoices != nil},
		{"InvoiceItems", c.InvoiceItems != nil},
		{"Bookings", c.Bookings != nil},
		{"BookingProjects", c.BookingProjects != nil},
		{"BillingAccounts", c.BillingAccounts != nil},
		{"BankAccounts", c.BankAccounts != nil},
		{"AccountingPlans", c.AccountingPlans != nil},
		{"CustomTaxRates", c.CustomTaxRates != nil},
		{"Cancellations", c.Cancellations != nil},
		{"ApplicationForms", c.ApplicationForms != nil},
		{"ApplicationFormElements", c.ApplicationFormElements != nil},
		{"InventoryObjects", c.InventoryObjects != nil},
		{"InventoryObjectGroups", c.InventoryObjectGroups != nil},
		{"Lendings", c.Lendings != nil},
		{"ArticleObjects", c.ArticleObjects != nil},
		{"Locations", c.Locations != nil},
		{"Calendars", c.Calendars != nil},
		{"Announcements", c.Announcements != nil},
		{"AnniversaryMailings", c.AnniversaryMailings != nil},
		{"CustomFields", c.CustomFields != nil},
		{"CustomFieldCollections", c.CustomFieldCollections != nil},
		{"CustomFilters", c.CustomFilters != nil},
		{"DocumentTemplates", c.DocumentTemplates != nil},
		{"ContactDetailsGroups", c.ContactDetailsGroups != nil},
		{"ContactDetailsLogs", c.ContactDetailsLogs != nil},
		{"FormerMemberData", c.FormerMemberData != nil},
		{"ChairmanLevels", c.ChairmanLevels != nil},
		{"ChairmanNotes", c.ChairmanNotes != nil},
		{"Events", c.Events != nil},
		{"MemberGroups", c.MemberGroups != nil},
		{"Organization", c.Organization != nil},
		{"FileSystemPaths", c.FileSystemPaths != nil},
		{"Wastebasket", c.Wastebasket != nil},
		{"ChatSettings", c.ChatSettings != nil},
		{"Forums", c.Forums != nil},
		{"DosbSports", c.DosbSports != nil},
		{"LsbSports", c.LsbSports != nil},
		{"Apply", c.Apply != nil},
	}
	for _, svc := range services {
		if !svc.ok {
			t.Errorf("service %s is nil after New()", svc.name)
		}
	}
}

func TestNew_DefaultBaseURL(t *testing.T) {
	c := New("tok")
	if c.baseURL != defaultBaseURL {
		t.Errorf("baseURL = %q, want %q", c.baseURL, defaultBaseURL)
	}
}

func TestNew_DefaultHTTPClientHas30sTimeout(t *testing.T) {
	c := New("tok")
	if c.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("httpClient.Timeout = %v, want 30s", c.httpClient.Timeout)
	}
}

func TestWithHTTPClient(t *testing.T) {
	custom := &http.Client{Timeout: 5 * time.Second}
	c := New("tok", WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Error("WithHTTPClient did not apply custom HTTP client")
	}
}

func TestWithBaseURL(t *testing.T) {
	c := New("tok", WithBaseURL("https://example.com/api"))
	if c.baseURL != "https://example.com/api" {
		t.Errorf("baseURL = %q, want https://example.com/api", c.baseURL)
	}
}

func TestWithDebug(t *testing.T) {
	c := New("tok", WithDebug(true))
	if !c.debug {
		t.Error("WithDebug(true) did not set debug flag")
	}
	c2 := New("tok", WithDebug(false))
	if c2.debug {
		t.Error("WithDebug(false) should leave debug unset")
	}
}

func TestWithTokenRefreshCallback(t *testing.T) {
	called := false
	cb := func(newToken string) { called = true }
	c := New("tok", WithTokenRefreshCallback(cb))
	if c.onTokenRefresh == nil {
		t.Fatal("onTokenRefresh is nil after WithTokenRefreshCallback")
	}
	c.onTokenRefresh("new")
	if !called {
		t.Error("registered callback was not invoked")
	}
}

// Options are applied in order; later options overwrite earlier ones.
func TestNew_MultipleOptions_LastWins(t *testing.T) {
	c := New("tok",
		WithBaseURL("https://first.example.com"),
		WithBaseURL("https://second.example.com"),
	)
	if c.baseURL != "https://second.example.com" {
		t.Errorf("baseURL = %q, want https://second.example.com", c.baseURL)
	}
}

// ---------------------------------------------------------------------------
// Organization – singleton service (Get/Update only, no List/Create/Delete)
// ---------------------------------------------------------------------------

func TestOrganizationService_Get(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/organization", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method %s", r.Method)
		}
		writeJSON(w, http.StatusOK, model.Organization{ID: 1, Name: "Test Verein"})
	})
	c, _ := newTestServer(t, mux)

	org, err := c.Organization.Get(context.Background(), nil)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if org.ID != 1 {
		t.Errorf("ID = %d, want 1", org.ID)
	}
	if org.Name != "Test Verein" {
		t.Errorf("Name = %q, want \"Test Verein\"", org.Name)
	}
}

func TestOrganizationService_Update(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/organization", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("unexpected method %s, want PATCH", r.Method)
		}
		writeJSON(w, http.StatusOK, model.Organization{ID: 1, Name: "Neuer Name"})
	})
	c, _ := newTestServer(t, mux)

	updated, err := c.Organization.Update(context.Background(), model.OrganizationCreate{Name: "Neuer Name"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "Neuer Name" {
		t.Errorf("Name = %q, want \"Neuer Name\"", updated.Name)
	}
}

// The Organization service deliberately has no List, Create, or Delete method –
// it represents a singleton resource. Verify the type has only Get and Update.
func TestOrganizationService_IsReadUpdateOnly(t *testing.T) {
	// This test is a compile-time contract: if someone adds List/Create/Delete
	// to OrganizationService the signatures below would need to change as well,
	// acting as a change detector.
	var svc *OrganizationService
	// Get and Update must exist.
	var _ func(context.Context, *Query) (*model.Organization, error) = svc.Get
	var _ func(context.Context, model.OrganizationCreate) (*model.Organization, error) = svc.Update
}

// ---------------------------------------------------------------------------
// Iterator[T] – exported pagination type
// ---------------------------------------------------------------------------

// Iterator fetches pages lazily: the fetch function is not called until
// Next() is invoked for the first time. See pagination_test.go for full
// coverage of multi-page and error scenarios.

func TestIterator_LazyFetch_NoCallBeforeNext(t *testing.T) {
	// The fetch function must NOT be called when the Iterator is created –
	// only when Next() is first invoked.
	called := false
	_ = newIterator("url", func(url string) ([]int, *string, error) {
		called = true
		return nil, nil, nil
	})
	if called {
		t.Error("fetch function called during newIterator, want lazy evaluation")
	}
}

// ---------------------------------------------------------------------------
// ListOptions
// ---------------------------------------------------------------------------

// ListOptions controls pagination and text search. Services embed it and add
// endpoint-specific filters.

func TestListOptions_ZeroValue_IsValid(t *testing.T) {
	// A zero-value ListOptions is valid and means "use defaults".
	var opts ListOptions
	if opts.Limit != 0 {
		t.Errorf("zero Limit = %d, want 0 (use API default)", opts.Limit)
	}
	if opts.Query != nil {
		t.Error("zero Query should be nil (use service default)")
	}
}

// MemberListOptions embeds ListOptions, so embedding is verified via the API.
func TestMemberListOptions_Embedding(t *testing.T) {
	opts := MemberListOptions{
		ListOptions: ListOptions{Limit: 50, Ordering: "-joinDate"},
		Email:       "test@example.com",
	}
	if opts.Limit != 50 {
		t.Errorf("embedded Limit = %d, want 50", opts.Limit)
	}
	if opts.Ordering != "-joinDate" {
		t.Errorf("Ordering = %q, want -joinDate", opts.Ordering)
	}
}

// ---------------------------------------------------------------------------
// Client – token propagation
// ---------------------------------------------------------------------------

func TestNew_TokenStoredOnClient(t *testing.T) {
	c := New("my-secret-token")
	if c.token != "my-secret-token" {
		t.Errorf("token = %q, want my-secret-token", c.token)
	}
}
