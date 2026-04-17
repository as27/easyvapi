package easyvapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/as27/easyvapi/model"
)

// MemberService manages all CRUD operations on the /member endpoint.
// Use this service to list, retrieve, create, update, and delete member records.
type MemberService struct {
	client *Client
}

// defaultMemberQuery requests all fields defined in model.Member so the API
// only returns what the client can decode. Callers can override this via
// ListOptions.Query to request a smaller subset of fields for better performance.
var defaultMemberQuery = NewQuery().
	Fields("id", "membershipNumber", "joinDate", "resignationDate",
		"resignationNoticeDate", "declarationOfApplication",
		"_paymentStartDate", "paymentAmount", "paymentIntervallMonths",
		"useBalanceForMembershipFee", "bulletinBoardNewPostNotification",
		"integrationDosbSport", "integrationDosbGender",
		"integrationLsbSport", "integrationLsbGender",
		"_isApplication", "_relatedMember", "relatedMembers", "org",
		"_deleteAfterDate", "_deletedBy",
		"declarationOfResignation", "declarationOfConsent", "sepaMandateFile",
		"customFields",
		"_applicationDate", "_applicationWasAcceptedAt",
		"_isChairman", "_chairmanPermissionGroup", "_profilePicture",
		"emailOrUserName", "signatureText", "_editableByRelatedMembers",
		"requirePasswordChange", "_isBlocked", "blockReason",
		"applicationKind", "wantsToCancelAt", "cancelReason",
		"showWarningsAndNotesToAdminsInProfile", "applicationForm",
		"_isMatrixSearchable", "matrixBlockReason", "blockedFromMatrix",
		"_matrixCommunicationPermission", "useMatrixGroupSettings").
	Nested("contactDetails", "id", "firstName", "familyName", "salutation",
		"street", "zip", "city", "country", "privateEmail", "primaryEmail",
		"privatePhone", "mobilePhone", "dateOfBirth").
	Nested("memberGroups", "id", "memberGroup{id,name,short}")

// MemberListOptions holds all filter and pagination options for Member list
// requests.
type MemberListOptions struct {
	ListOptions
	// Email filters members by their primary e-mail address.
	Email string
	// MembershipNumber filters by the member's numeric membership number.
	MembershipNumber string
	// IsBlocked filters by the blocked status of the member.
	IsBlocked *bool
	// IsApplication filters members that are still in the application process.
	IsApplication *bool
	// JoinDateGte filters members whose join date is on or after this date
	// (YYYY-MM-DD).
	JoinDateGte string
	// JoinDateLte filters members whose join date is on or before this date
	// (YYYY-MM-DD).
	JoinDateLte string
	// ResignationDateIsNull when set to true returns only active members (no
	// resignation date), false returns only resigned members.
	ResignationDateIsNull *bool
	// MemberGroups filters by the given member group IDs.
	MemberGroups []int
}

// memberListParams converts opts into URL query parameters.
func memberListParams(opts *MemberListOptions) url.Values {
	params := url.Values{}
	if opts == nil {
		applyListOptions(params, ListOptions{}, defaultMemberQuery)
		return params
	}
	applyListOptions(params, opts.ListOptions, defaultMemberQuery)
	if opts.Email != "" {
		params.Set("contactDetails__primaryEmail", opts.Email)
	}
	if opts.MembershipNumber != "" {
		params.Set("membershipNumber", opts.MembershipNumber)
	}
	if opts.IsBlocked != nil {
		params.Set("isBlocked", strconv.FormatBool(*opts.IsBlocked))
	}
	if opts.IsApplication != nil {
		params.Set("isApplication", strconv.FormatBool(*opts.IsApplication))
	}
	if opts.JoinDateGte != "" {
		params.Set("joinDate__gte", opts.JoinDateGte)
	}
	if opts.JoinDateLte != "" {
		params.Set("joinDate__lte", opts.JoinDateLte)
	}
	if opts.ResignationDateIsNull != nil {
		params.Set("resignationDate__isnull", strconv.FormatBool(*opts.ResignationDateIsNull))
	}
	for _, id := range opts.MemberGroups {
		params.Add("memberGroups", strconv.Itoa(id))
	}
	return params
}

// List returns a lazy Iterator over all Member records matching opts.
// Pages are fetched on-demand as iteration progresses, making this
// memory-efficient for large result sets.
// Pass nil for opts to use default filtering and the standard query.
//
// Example:
//
//	iter := client.Members.List(ctx, nil)
//	for iter.Next() {
//		member := iter.Value()
//		fmt.Printf("%s: %s\n", member.MembershipNumber, member.ContactDetails.FirstName)
//	}
//	if err := iter.Err(); err != nil {
//		log.Fatal(err)
//	}
func (s *MemberService) List(ctx context.Context, opts *MemberListOptions) *Iterator[model.Member] {
	startURL := s.client.buildURL("/member", memberListParams(opts))
	return newIterator(startURL, func(pageURL string) ([]model.Member, *string, error) {
		return fetchPage[model.Member](s.client, ctx, pageURL)
	})
}

// ListAll fetches all Member records matching opts and returns them as a slice.
// This is a convenience wrapper around [MemberService.List] that collects all pages
// into memory. For very large result sets, consider using List with Iterator
// to process records one at a time.
// Pass nil for opts to use default filtering and the standard query.
//
// Example:
//
//	members, err := client.Members.ListAll(ctx, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Found %d members\n", len(members))
func (s *MemberService) ListAll(ctx context.Context, opts *MemberListOptions) ([]model.Member, error) {
	var all []model.Member
	iter := s.List(ctx, opts)
	for iter.Next() {
		all = append(all, iter.Value())
	}
	return all, iter.Err()
}

// Get retrieves a single Member by its ID. The query parameter optionally
// restricts which fields are returned. Pass nil to use the default query,
// or create a custom query with [NewQuery] to fetch only needed fields.
//
// Example:
//
//	member, err := client.Members.Get(ctx, 123456, nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%s\n", member.ContactDetails.FirstName)
//
// With a custom query:
//
//	query := easyvapi.NewQuery().
//		Fields("id", "membershipNumber").
//		Nested("contactDetails", "firstName", "familyName")
//	member, err := client.Members.Get(ctx, 123456, query)
func (s *MemberService) Get(ctx context.Context, id int, query *Query) (*model.Member, error) {
	params := url.Values{}
	if qs := query.String(); qs != "" {
		params.Set("query", qs)
	}
	resp, err := s.client.get(ctx, fmt.Sprintf("/member/%d", id), params)
	if err != nil {
		return nil, err
	}
	var m model.Member
	if err := s.client.decodeJSON(resp, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Create creates a new Member and returns the created record with all fields
// populated by the API. The ID field is assigned by the server.
//
// Example:
//
//	newMember := &model.MemberCreate{
//		JoinDate:      "2026-03-31",
//		PaymentAmount: 25.00,
//	}
//	created, err := client.Members.Create(ctx, *newMember)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Created member with ID: %d\n", created.ID)
func (s *MemberService) Create(ctx context.Context, m model.MemberCreate) (*model.Member, error) {
	resp, err := s.client.do(ctx, "POST", s.client.buildURL("/member", nil), m)
	if err != nil {
		return nil, err
	}
	var created model.Member
	if err := s.client.decodeJSON(resp, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

// Update applies a partial update (PATCH) to the Member with the given ID.
// Only non-zero fields in the MemberCreate struct are sent to the API.
// The updated record is returned with all current field values.
//
// Example:
//
//	updated, err := client.Members.Update(ctx, 123456, model.MemberCreate{
//		PaymentAmount: 30.00,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Updated payment amount to: %.2f\n", updated.PaymentAmount)
func (s *MemberService) Update(ctx context.Context, id int, m model.MemberCreate) (*model.Member, error) {
	resp, err := s.client.do(ctx, "PATCH", s.client.buildURL(fmt.Sprintf("/member/%d", id), nil), m)
	if err != nil {
		return nil, err
	}
	var updated model.Member
	if err := s.client.decodeJSON(resp, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

// Delete removes the Member with the given ID. The member record is
// immediately deleted and cannot be recovered. No confirmation is required.
//
// Example:
//
//	err := client.Members.Delete(ctx, 123456)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Member deleted")
func (s *MemberService) Delete(ctx context.Context, id int) error {
	resp, err := s.client.do(ctx, "DELETE", s.client.buildURL(fmt.Sprintf("/member/%d", id), nil), nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
