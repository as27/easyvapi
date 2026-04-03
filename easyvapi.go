// Package easyvapi provides a Go client for the easyVerein API v2.0.
//
// # Basic Usage
//
// Create a client with your API token and call service methods:
//
//	client := easyvapi.New("your-api-token")
//	members, err := client.Members.ListAll(context.Background(), nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, m := range members {
//		fmt.Println(m.ID, m.MembershipNumber)
//	}
//
// # Features
//
// - Automatic pagination with lazy Iterator[T]
// - Automatic token refresh when signaled by API
// - Rate-limit handling with exponential backoff
// - Flexible query filtering for field selection
// - Comprehensive error handling
//
// # Rate Limiting
//
// The API enforces a limit of 100 requests per minute. The client automatically
// throttles when X-RateLimit-Remaining falls below 5. Use errors.As to detect
// RateLimitError and implement retry logic:
//
//	var rateLimitErr *easyvapi.RateLimitError
//	if errors.As(err, &rateLimitErr) {
//		time.Sleep(rateLimitErr.RetryAfter)
//		// retry
//	}
package easyvapi

import (
	"net/http"
	"time"

	"github.com/as27/easyvapi/model"
)

const defaultBaseURL = "https://easyverein.com/api/v2.0"

// Client is the main entry point for interacting with the easyVerein API.
// It provides access to all resource services and manages authentication,
// rate limiting, and token refresh automatically.
//
// Create a new Client with [New], which returns a fully initialized client
// with default configuration (30s timeout, standard base URL). Customize
// behavior using functional options like [WithHTTPClient] or [WithDebug].
type Client struct {
	// Members provides access to the /member endpoint for CRUD operations on members.
	Members *MemberService
	// ContactDetails provides access to the /contact-details endpoint.
	ContactDetails *ContactDetailsService
	// Invoices provides access to the /invoice endpoint for financial documents.
	Invoices *InvoiceService
	// InvoiceItems provides access to the /invoice-item endpoint for invoice line items.
	InvoiceItems *InvoiceItemService
	// Bookings provides access to the /booking endpoint for financial transactions.
	Bookings *BookingService
	// BookingProjects provides access to the /booking-project endpoint.
	BookingProjects *BookingProjectService
	// BillingAccounts provides access to the /billing-account endpoint.
	BillingAccounts *BillingAccountService
	// BankAccounts provides access to the /bank-account endpoint.
	BankAccounts *BankAccountService
	// AccountingPlans provides access to the /accounting-plan endpoint.
	AccountingPlans *AccountingPlanService
	// CustomTaxRates provides access to the /custom-tax-rate endpoint.
	CustomTaxRates *CustomTaxRateService
	// Cancellations provides access to the /cancellation endpoint.
	Cancellations *CancellationService
	// ApplicationForms provides access to the /application-form endpoint.
	ApplicationForms *ApplicationFormService
	// ApplicationFormElements provides access to the /application-form-element endpoint.
	ApplicationFormElements *ApplicationFormElementService
	// InventoryObjects provides access to the /inventory-object endpoint.
	InventoryObjects *InventoryObjectService
	// InventoryObjectGroups provides access to the /inventory-object-group endpoint.
	InventoryObjectGroups *InventoryObjectGroupService
	// Lendings provides access to the /lending endpoint.
	Lendings *LendingService
	// ArticleObjects provides access to the /article-object endpoint.
	ArticleObjects *ArticleObjectService
	// Locations provides access to the /location endpoint for event venues.
	Locations *LocationService
	// Calendars provides access to the /calendar endpoint.
	Calendars *CalendarService
	// Announcements provides access to the /announcement endpoint.
	Announcements *AnnouncementService
	// AnniversaryMailings provides access to the /anniversary-mailing endpoint.
	AnniversaryMailings *AnniversaryMailingService
	// CustomFields provides access to the /custom-field endpoint.
	CustomFields *CustomFieldService
	// CustomFieldCollections provides access to the /custom-field-collection endpoint.
	CustomFieldCollections *CustomFieldCollectionService
	// CustomFilters provides access to the /custom-filter endpoint.
	CustomFilters *CustomFilterService
	// DocumentTemplates provides access to the /document-template endpoint.
	DocumentTemplates *DocumentTemplateService
	// ContactDetailsGroups provides access to the /contact-details-group endpoint.
	ContactDetailsGroups *ContactDetailsGroupService
	// ContactDetailsLogs provides access to the /contact-details-log endpoint.
	ContactDetailsLogs *ContactDetailsLogService
	// FormerMemberData provides read-only access to the /former-member-data endpoint.
	FormerMemberData *FormerMemberDataService
	// ChairmanLevels provides access to the /chairman-level endpoint.
	ChairmanLevels *ChairmanLevelService
	// ChairmanNotes provides access to the /chairman-note endpoint.
	ChairmanNotes *ChairmanNoteService
	// Events provides access to the /event endpoint for calendar events.
	Events *EventService
	// MemberGroups provides access to the /member-group endpoint for member categories.
	MemberGroups *MemberGroupService
	// Organization provides access to the /organization endpoint (singleton).
	Organization *OrganizationService
	// FileSystemPaths provides access to the /file-system-path endpoint.
	FileSystemPaths *FileSystemPathService
	// Wastebasket provides access to the /wastebasket endpoint for deleted objects.
	Wastebasket *WastebasketService
	// ChatSettings provides access to the /chat-settings endpoint (singleton).
	ChatSettings *ChatSettingsService
	// Forums provides access to the /forum endpoint.
	Forums *ForumService
	// DosbSports provides access to the /dosb-sport endpoint.
	DosbSports *DosbSportService
	// LsbSports provides access to the /lsb-sport endpoint.
	LsbSports *LsbSportService
	// Apply provides access to the /apply endpoint for membership applications.
	Apply *ApplyService

	httpClient     *http.Client
	baseURL        string
	token          string
	debug          bool
	onTokenRefresh func(newToken string)
}

// Option is a functional option for configuring a Client.
// Options are applied in order during [New].
type Option func(*Client)

// WithHTTPClient sets a custom http.Client for the API client.
// The default client has a 30-second timeout. Use this to customize
// the timeout, TLS configuration, or other HTTP behavior.
//
//	client := easyvapi.New(token, easyvapi.WithHTTPClient(&http.Client{
//		Timeout: 60 * time.Second,
//	}))
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// WithBaseURL overrides the default API base URL.
// Default is "https://easyverein.com/api/v2.0".
// Use this for testing or custom deployments.
//
//	client := easyvapi.New(token, easyvapi.WithBaseURL("https://test.example.com/api/v2.0"))
func WithBaseURL(u string) Option {
	return func(c *Client) {
		c.baseURL = u
	}
}

// WithDebug enables debug mode. When enabled, JSON decoding rejects
// unknown fields, which helps detect API schema changes early.
// Useful during development to catch incompatibilities before
// they reach production.
//
//	client := easyvapi.New(token, easyvapi.WithDebug(true))
func WithDebug(debug bool) Option {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithTokenRefreshCallback registers a callback function that is called
// whenever the API signals a token refresh via the tokenRefreshNeeded header.
// The callback receives the new token and should persist it securely
// (e.g., write to a config file or encrypted storage).
//
//	client := easyvapi.New(token, easyvapi.WithTokenRefreshCallback(func(newToken string) {
//		os.WriteFile("token.txt", []byte(newToken), 0600)
//	}))
func WithTokenRefreshCallback(fn func(newToken string)) Option {
	return func(c *Client) {
		c.onTokenRefresh = fn
	}
}

// New creates a new Client authenticated with the given token.
// The default client uses a 30-second HTTP timeout and the standard
// easyVerein API base URL. Customize behavior with functional options.
//
//	client := easyvapi.New("my-api-token")
//	client2 := easyvapi.New("token", easyvapi.WithDebug(true), easyvapi.WithHTTPClient(customClient))
func New(token string, opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		baseURL:    defaultBaseURL,
		token:      token,
	}
	for _, o := range opts {
		o(c)
	}
	c.Members = &MemberService{client: c}
	c.ContactDetails = &ContactDetailsService{client: c}
	c.Invoices = &InvoiceService{client: c}
	c.InvoiceItems = &InvoiceItemService{client: c}
	c.Bookings = &BookingService{client: c}
	c.BookingProjects = &BookingProjectService{client: c}
	c.BillingAccounts = &BillingAccountService{client: c}
	c.BankAccounts = &BankAccountService{client: c}
	c.AccountingPlans = &AccountingPlanService{client: c}
	c.CustomTaxRates = &CustomTaxRateService{client: c}
	c.Cancellations = &CancellationService{client: c}
	c.ApplicationForms = &ApplicationFormService{client: c}
	c.ApplicationFormElements = &ApplicationFormElementService{client: c}
	c.InventoryObjects = &InventoryObjectService{client: c}
	c.InventoryObjectGroups = &InventoryObjectGroupService{client: c}
	c.Lendings = &LendingService{client: c}
	c.ArticleObjects = &ArticleObjectService{client: c}
	c.Locations = &LocationService{client: c}
	c.Calendars = &CalendarService{client: c}
	c.Announcements = &AnnouncementService{client: c}
	c.AnniversaryMailings = &AnniversaryMailingService{client: c}
	c.CustomFields = &CustomFieldService{client: c}
	c.CustomFieldCollections = &CustomFieldCollectionService{client: c}
	c.CustomFilters = &CustomFilterService{client: c}
	c.DocumentTemplates = &DocumentTemplateService{client: c}
	c.ContactDetailsGroups = &ContactDetailsGroupService{client: c}
	c.ContactDetailsLogs = &ContactDetailsLogService{client: c}
	c.FormerMemberData = &FormerMemberDataService{client: c}
	c.ChairmanLevels = &ChairmanLevelService{client: c}
	c.ChairmanNotes = &ChairmanNoteService{client: c}
	c.Events = &EventService{client: c}
	c.MemberGroups = &MemberGroupService{client: c}
	c.Organization = &OrganizationService{client: c}
	c.FileSystemPaths = &FileSystemPathService{client: c}
	c.Wastebasket = &WastebasketService{client: c}
	c.ChatSettings = &ChatSettingsService{client: c}
	c.Forums = &ForumService{client: c}
	c.DosbSports = &DosbSportService{client: c}
	c.LsbSports = &LsbSportService{client: c}
	c.Apply = &ApplyService{client: c}
	return c
}

// ensure model is imported (used by service files in the same package).
var _ = model.Member{}
