package easyvapi

import (
	"strings"
)

// ListOptions holds generic pagination and filtering options shared by all
// List endpoints. Each service extends this with endpoint-specific filters.
// If Query is nil, the endpoint uses its default query which requests all
// fields defined in the corresponding model struct.
type ListOptions struct {
	// Limit is the maximum number of results per page. Default is 100.
	// Set to 0 to use the endpoint default (100).
	Limit int
	// Query restricts which fields are included in the response using the
	// easyVerein query syntax. If nil, the endpoint uses its default query.
	// Use [NewQuery] to build a custom query.
	Query *Query
	// Ordering is the field name to order results by. Prefix with "-" for descending order.
	// Example: "name" for ascending, "-joinDate" for descending.
	Ordering string
	// Search is a full-text search string applied across searchable fields.
	// The API performs case-insensitive substring matching.
	Search string
}

// Query builds an easyVerein field-selection query string of the form
// {field1,nested{f2},-excluded}. This allows clients to request only
// the fields they need, reducing response size and improving performance.
//
// Example:
//
//	q := easyvapi.NewQuery().
//		Fields("id", "name").
//		Nested("contactDetails", "email", "phone").
//		Exclude("password")
//	// Renders to: {id,name,contactDetails{email,phone},-password}
type Query struct {
	fields   []string
	nested   map[string][]string
	excluded []string
	// preserve insertion order for nested keys
	nestedOrder []string
}

// NewQuery creates an empty Query builder. Use the builder methods to add fields.
//
//	query := easyvapi.NewQuery().
//		Fields("id", "email").
//		Nested("contactDetails", "firstName", "familyName")
func NewQuery() *Query {
	return &Query{
		nested: make(map[string][]string),
	}
}

// Fields adds top-level fields to include in the response.
// Multiple calls accumulate fields.
//
//	query := easyvapi.NewQuery().
//		Fields("id", "name").
//		Fields("email")  // now includes id, name, email
func (q *Query) Fields(fields ...string) *Query {
	q.fields = append(q.fields, fields...)
	return q
}

// Nested adds a nested object selector. The name parameter is the field name
// (e.g., "contactDetails"), and the fields parameter specifies which nested
// fields to include (e.g., "id", "firstName", "email").
// Multiple calls to Nested with the same name accumulate fields.
//
//	query := easyvapi.NewQuery().
//		Nested("contactDetails", "id", "firstName", "familyName").
//		Nested("memberGroups", "id", "name")
func (q *Query) Nested(name string, fields ...string) *Query {
	if _, exists := q.nested[name]; !exists {
		q.nestedOrder = append(q.nestedOrder, name)
	}
	q.nested[name] = append(q.nested[name], fields...)
	return q
}

// Exclude marks fields to be excluded from the response. Excluded fields
// are prefixed with "-" in the query string. This is useful when you want
// most fields but need to omit sensitive or large fields.
//
//	query := easyvapi.NewQuery().
//		Fields("id", "name", "password").
//		Exclude("password")  // request all but exclude password
func (q *Query) Exclude(fields ...string) *Query {
	q.excluded = append(q.excluded, fields...)
	return q
}

// String renders the Query into the easyVerein query syntax.
// Returns an empty string if the query is nil or has no fields.
// Example output: {id,contactDetails{firstName},-password}
//
//	q := easyvapi.NewQuery().
//		Fields("id").
//		Nested("contactDetails", "email").
//		Exclude("password")
//	fmt.Println(q.String())  // Output: {id,contactDetails{email},-password}
func (q *Query) String() string {
	if q == nil {
		return ""
	}
	var parts []string
	parts = append(parts, q.fields...)
	for _, name := range q.nestedOrder {
		fs := q.nested[name]
		parts = append(parts, name+"{"+strings.Join(fs, ",")+"}")
	}
	for _, f := range q.excluded {
		parts = append(parts, "-"+f)
	}
	if len(parts) == 0 {
		return ""
	}
	return "{" + strings.Join(parts, ",") + "}"
}
