# easyvapi – easyVerein API v2.0 Go Client

**Quick Reference for Claude Code**

## Essentials

- **Init**: `client := easyvapi.New(token)` + services: `Members`, `ContactDetails`, `Invoices`, `Bookings`, `Events`, `MemberGroups`
- **List**: `iter := client.Members.List(ctx, nil)` (lazy Iterator) or `all := client.Members.ListAll(ctx, nil)` (all at once)
- **Get**: `client.Members.Get(ctx, id, nil)`
- **Create/Update/Delete**: `client.Members.Create/Update/Delete(ctx, ...)`

## How It Works

```go
// Standard flow: List → Iterate → Process
iter := client.Members.List(ctx, nil)
for iter.Next() {
    m := iter.Value()  // one member
}
if err := iter.Err(); err != nil { ... }

// Filtering
opts := &easyvapi.MemberListOptions{
    ListOptions: easyvapi.ListOptions{
        Limit: 100,                                    // default
        Query: easyvapi.NewQuery().Fields("id"),      // custom fields
    },
}
members, err := client.Members.ListAll(ctx, opts)
```

## Key Features

| Feature | Details |
|---|---|
| **Pagination** | Lazy Iterator[T] - pages fetch on-demand. ListAll() collects all. Default limit: 100/page. |
| **Query** | `NewQuery().Fields("id","name").Nested("contactDetails","email").Exclude("password")` → `{id,name,contactDetails{email},-password}` |
| **Default Queries** | Each service (e.g., `defaultMemberQuery`) auto-selects model fields. Nil opts → use default. Exception: `member_group` has no query support. |
| **Rate-Limit** | Auto-throttles if X-RateLimit-Remaining < 5 (sleeps 10s). HTTP 429 → `*RateLimitError{RetryAfter time.Duration}` |
| **Token Refresh** | Auto-refreshes on `tokenRefreshNeeded` header, retries once. Callback via `WithTokenRefreshCallback()`. |
| **HTTP** | 30s timeout, `Accept: application/json` always set, proper URL encoding. |
| **Errors** | `*APIError{StatusCode, Message, Detail}` (non-2xx), `*RateLimitError{RetryAfter}` (429). Use `errors.As()`. |

## Services (40 total)

Each has: `List(ctx, opts) *Iterator[T]`, `ListAll(ctx, opts) ([]T, error)`, `Get(ctx, id, query)`, `Create`, `Update`, `Delete`

| Service | Model | Filters | Notes |
|---|---|---|---|
| Members | Member | Email, MembershipNumber, IsBlocked, Dates, Groups | Related nested objects: contactDetails, memberGroups |
| ContactDetails | ContactDetails | Country, FirstName, FamilyName, IsCompany | Flat struct |
| Invoices | Invoice | Kind, IsDraft, IsTemplate, DateRange | Financial docs |
| InvoiceItems | InvoiceItem | RelatedInvoice, BillingAccount, Title | Invoice line items |
| Bookings | Booking | DateRange, BillingAccount | + BulkCreate/BulkUpdate |
| BookingProjects | BookingProject | Name | **No query parameter support** |
| BillingAccounts | BillingAccount | Name, AccountKind | **No query parameter support** |
| BankAccounts | BankAccount | Name | **No query parameter support** |
| AccountingPlans | AccountingPlan | — | Kontenplan |
| CustomTaxRates | CustomTaxRate | TaxName | Steuersätze |
| Cancellations | — | — | Submit only (POST /cancellation) |
| ApplicationForms | ApplicationForm | Title, Public, Language, FormularKind | **No query parameter support** |
| ApplicationFormElements | ApplicationFormElement | ApplicationForm, Kind, Required | **No query parameter support**; + BulkCreate/BulkUpdate |
| InventoryObjects | InventoryObject | Name, Identifier, LendingAvailable | **No query parameter support** |
| InventoryObjectGroups | InventoryObjectGroup | Name | **No query parameter support** |
| Lendings | Lending | State, InventoryObject, LendingPerson | **No query parameter support**; + BulkCreate/BulkUpdate |
| ArticleObjects | ArticleObject | Name, Kind | **No query parameter support** |
| Locations | Location | Name, Country, Zip | **No query parameter support** |
| Calendars | Calendar | Name | **No query parameter support** |
| Announcements | Announcement | Platform, ShowBanner | **No query parameter support** |
| AnniversaryMailings | AnniversaryMailing | — | **No query parameter support** |
| CustomFields | CustomField | Label, FieldKind, FieldCollection, ShowInMemberArea | **No query parameter support** |
| CustomFieldCollections | CustomFieldCollection | — | Gruppen eigener Felder |
| CustomFilters | CustomFilter | Name, Model | **No query parameter support** |
| DocumentTemplates | DocumentTemplate | Title, DocumentKind | **No query parameter support** |
| ContactDetailsGroups | ContactDetailsGroup | Name | **No query parameter support** |
| ContactDetailsLogs | ContactDetailsLog | ContactDetails (ID) | Log entries per contact |
| FormerMemberData | FormerMemberData | — | Read-only; **No query parameter support** |
| ChairmanLevels | ChairmanLevel | — | **No query parameter support** |
| ChairmanNotes | ChairmanNote | DateGte, DateLte | Interne Vorstandsnotizen |
| Events | Event | DateRange, Calendar, IsPublic | + Copy, GenerateInvoices, InviteGroups, ListParticipations, CreateParticipation, UpdateParticipation, DeleteParticipation |
| MemberGroups | MemberGroup | Name | **No query parameter support** |
| Organization | Organization | — | Singleton: Get, Update only |
| FileSystemPaths | FileSystemPath | Name | **No query parameter support** |
| Wastebasket | WastebasketItem | Model | List, Restore; **No query parameter support** |
| ChatSettings | ChatSettings | — | Singleton: Get, Update only |
| Forums | Forum | Name | **No query parameter support** |
| DosbSports | DosbSport | Name | **No query parameter support** |
| LsbSports | LsbSport | Name | **No query parameter support** |
| Apply | — | — | Submit only (POST /apply) |

## Code Patterns

**Create:**
```go
booking := &model.BookingCreate{Amount: 100.0, BillingAccount: 123, Date: "2026-03-31"}
created, err := client.Bookings.Create(ctx, *booking)
```

**Update (PATCH):**
```go
updated, err := client.Members.Update(ctx, id, model.MemberCreate{PaymentAmount: 25.0})
```

**Custom Query (speed up):**
```go
q := easyvapi.NewQuery().Fields("id", "membershipNumber").Nested("contactDetails", "email")
members, err := client.Members.ListAll(ctx, &easyvapi.MemberListOptions{
    ListOptions: easyvapi.ListOptions{Query: q},
})
```

## File Map

- `easyvapi.go` - Client, New(), Options (WithHTTPClient, WithDebug, WithTokenRefreshCallback, WithBaseURL)
- `options.go` - ListOptions, Query builder
- `pagination.go` - Iterator[T]
- `errors.go` - APIError, RateLimitError
- `request.go` - HTTP internals (do/doOnce, rate-limit, token-refresh)
- `helpers.go` - applyListOptions, fetchPage
- `member.go`, `contact_details.go`, `contact_details_group.go`, `contact_details_log.go`, `invoice.go`, `invoice_item.go`, `booking.go`, `booking_project.go`, `billing_account.go`, `bank_account.go`, `accounting_plan.go`, `custom_tax_rate.go`, `cancellation.go`, `former_member_data.go`, `chairman_level.go`, `chairman_note.go`, `custom_field.go`, `custom_field_collection.go`, `custom_filter.go`, `document_template.go`, `event.go`, `member_group.go` - Services + defaults
- `model/` - Data structs

## Adding New Endpoint

1. Create `model/newtype.go` with data struct
2. Create `newtype.go` service: `type NewTypeService struct{client *Client}`, `defaultNewTypeQuery` var, `List/ListAll/Get/Create/Update/Delete`
3. Update `easyvapi.go`: add field to Client, initialize in New()
4. In service params func: `applyListOptions(params, opts.ListOptions, defaultNewTypeQuery)` (or `nil` if no query)
