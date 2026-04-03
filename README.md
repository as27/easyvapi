# easyvapi – Go Client für easyVerein API v2.0

Ein moderner, typsicherer Go Client für die [easyVerein](https://www.easyverein.com/) REST API v2.0 mit automatischer Pagination, Rate-Limiting und Token-Refresh.

## Features

- ✅ **Vollständige CRUD-Operationen** für alle Ressourcen (Members, Contacts, Contact Groups, Contact Logs, Invoices, Invoice Items, Bookings, Booking Projects, Billing Accounts, Bank Accounts, Accounting Plans, Custom Tax Rates, Custom Fields, Custom Field Collections, Custom Filters, Document Templates, Chairman Levels, Chairman Notes, Former Member Data, Events, Member Groups, Locations, Calendars, Announcements, Anniversary Mailings, Application Forms, Application Form Elements, Inventory Objects, Inventory Object Groups, Lendings, Article Objects, Organization, File System Paths, Wastebasket, Chat Settings, Forums, DOSB Sports, LSB Sports)
- ✅ **Lazy Pagination** mit `Iterator[T]` – nur die benötigten Felder abrufen, Seiten bei Bedarf laden
- ✅ **Automatisches Token-Refresh** – Token wird automatisch erneuert, wenn die API es signalisiert
- ✅ **Intelligentes Rate-Limit-Handling** – automatisches Drosseln bei niedriger Rate-Limit-Verbrauch
- ✅ **Flexible Query-Filter** – einfache Builder-API zum Auswählen von Feldern und verschachtelten Objekten
- ✅ **Standard-Abfragen** – jeder Endpoint hat sinnvolle Defaults, die nur relevante Felder abrufen
- ✅ **Großzügige Fehlerausgaben** – aussagekräftige Error-Messages mit API-Details

## Installation

```bash
go get github.com/as27/easyvapi
```

## Schnelstart

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/as27/easyvapi"
)

func main() {
	// Client erstellen
	client := easyvapi.New("dein-api-token")

	// Alle Mitglieder abrufen (mit Standard-Query)
	members, err := client.Members.ListAll(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, m := range members {
		fmt.Printf("ID: %d, Mitgliedsnummer: %s\n", m.ID, m.MembershipNumber)
	}
}
```

## Verwendungsbeispiele

### 1. Mitglieder filtern und abrufen

```go
import "github.com/as27/easyvapi"

// Alle aktiven Mitglieder abrufen (ohne Resignationsdatum)
opts := &easyvapi.MemberListOptions{
	ListOptions: easyvapi.ListOptions{Limit: 100},
}
isActive := true
opts.ResignationDateIsNull = &isActive

members, err := client.Members.ListAll(ctx, opts)
if err != nil {
	log.Fatal(err)
}
```

### 2. Mit Iterator – Speicher-effizient arbeiten

```go
// Iterator für Lazy-Loading verwenden
iter := client.Members.List(ctx, nil)
for iter.Next() {
	m := iter.Value()
	fmt.Printf("Verarbeite Mitglied: %s\n", m.MembershipNumber)
	// Nur eine Seite wird im Speicher gehalten
}
if err := iter.Err(); err != nil {
	log.Fatal(err)
}
```

### 3. Nur bestimmte Felder abrufen (Query-Filter)

```go
// Nur ID und E-Mail abrufen – viel schneller!
query := easyvapi.NewQuery().
	Fields("id").
	Nested("contactDetails", "privateEmail")

opts := &easyvapi.MemberListOptions{
	ListOptions: easyvapi.ListOptions{Query: query},
}

members, err := client.Members.ListAll(ctx, opts)
if err != nil {
	log.Fatal(err)
}
```

### 4. Kontaktdaten mit Filter abrufen

```go
// Alle Kontakte aus Deutschland filtern
opts := &easyvapi.ContactDetailsListOptions{
	Country: "DE",
}

contacts, err := client.ContactDetails.ListAll(ctx, opts)
if err != nil {
	log.Fatal(err)
}
```

### 5. Buchungen in einem Zeitraum abrufen

```go
import "github.com/as27/easyvapi"

// Buchungen von Januar bis März 2026
opts := &easyvapi.BookingListOptions{
	DateGt: "2026-01-01",
	DateLt: "2026-03-31",
}

bookings, err := client.Bookings.ListAll(ctx, opts)
if err != nil {
	log.Fatal(err)
}

for _, b := range bookings {
	fmt.Printf("Betrag: %.2f EUR, Datum: %s\n", b.Amount, b.Date)
}
```

### 6. Neue Buchung erstellen

```go
booking := &model.BookingCreate{
	Amount:         150.50,
	BillingAccount: 12345,
	Date:           "2026-03-31",
	Description:    "Mitgliedsbeitrag",
	Receiver:       "Mitglied Max Mustermann",
}

created, err := client.Bookings.Create(ctx, *booking)
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Buchung erstellt mit ID: %d\n", created.ID)
```

### 7. Veranstaltungen mit Suchtext filtern

```go
// Alle Events mit "Workshop" im Namen
opts := &easyvapi.EventListOptions{
	ListOptions: easyvapi.ListOptions{
		Search: "Workshop",
	},
}

events, err := client.Events.ListAll(ctx, opts)
if err != nil {
	log.Fatal(err)
}
```

### 8. Mitglied aktualisieren

```go
// Mitglied aktualisieren (PATCH)
updated, err := client.Members.Update(ctx, 123456, model.MemberCreate{
	PaymentAmount: 25.00,
})
if err != nil {
	log.Fatal(err)
}

fmt.Printf("Mitglied aktualisiert: %s\n", updated.MembershipNumber)
```

### 9. Rate-Limit-Fehler behandeln

```go
import "errors"

data, err := client.Members.ListAll(ctx, nil)
if err != nil {
	var rateLimitErr *easyvapi.RateLimitError
	if errors.As(err, &rateLimitErr) {
		fmt.Printf("Rate-Limit erreicht. Bitte %v warten.\n", rateLimitErr.RetryAfter)
		time.Sleep(rateLimitErr.RetryAfter)
		// Erneut versuchen...
	} else {
		log.Fatal(err)
	}
}
```

### 10. Token-Refresh-Callback

```go
// Token nach Refresh speichern
client := easyvapi.New(
	"initial-token",
	easyvapi.WithTokenRefreshCallback(func(newToken string) {
		fmt.Println("Token erneuert! Speichere neuen Token...")
		// Token in Datei/DB speichern
		os.WriteFile("token.txt", []byte(newToken), 0600)
	}),
)
```

## Query-Builder

Mit dem `Query`-Builder kannst du präzise wählen, welche Felder die API zurückgeben soll:

```go
query := easyvapi.NewQuery().
	Fields("id", "membershipNumber", "joinDate").
	Nested("contactDetails", "firstName", "familyName", "privateEmail").
	Exclude("password")

// Erzeugt: {id,membershipNumber,joinDate,contactDetails{firstName,familyName,privateEmail},-password}
```

## Konfiguration

### Custom HTTP-Client

```go
import "net/http"
import "time"

httpClient := &http.Client{
	Timeout: 60 * time.Second,
}

client := easyvapi.New(
	"token",
	easyvapi.WithHTTPClient(httpClient),
)
```

### Debug-Modus

Im Debug-Modus lehnt der JSON-Decoder unbekannte Felder ab – nützlich, um API-Änderungen frühzeitig zu erkennen:

```go
client := easyvapi.New(
	"token",
	easyvapi.WithDebug(true),
)
```

### Custom Base-URL

```go
client := easyvapi.New(
	"token",
	easyvapi.WithBaseURL("https://custom.easyverein.com/api/v2.0"),
)
```

## Standard-Queries pro Endpoint

Jeder Endpoint setzt automatisch eine sinnvolle Standard-Query, wenn keine eigene Query angegeben wird:

| Endpoint | Standard-Felder |
|---|---|
| **Members** | id, membershipNumber, joinDate, resignationDate, paymentAmount, paymentIntervallMonths, _isBlocked, _isApplication, _relatedMember, contactDetails{...}, memberGroups{...} |
| **ContactDetails** | id, firstName, familyName, salutation, street, zip, city, country, privateEmail, primaryEmail, privatePhone, mobilePhone, dateOfBirth |
| **Invoices** | id, invNumber, date, receiver, totalPrice, kind, isDraft, isTemplate |
| **InvoiceItems** | id, title, quantity, unitPrice, taxRate, taxName, description, billingAccount, gross |
| **Bookings** | id, amount, date, receiver, billingId |
| **BookingProjects** | *Keine Query-Unterstützung (API-Limitierung)* |
| **BillingAccounts** | *Keine Query-Unterstützung (API-Limitierung)* |
| **BankAccounts** | *Keine Query-Unterstützung (API-Limitierung)* |
| **AccountingPlans** | id, name, description |
| **CustomTaxRates** | id, taxName, customTaxRate |
| **ContactDetailsGroups** | *Keine Query-Unterstützung (API-Limitierung)* |
| **ContactDetailsLogs** | id, contactDetails, title, message, date |
| **FormerMemberData** | *Keine Query-Unterstützung (API-Limitierung); read-only* |
| **ChairmanLevels** | *Keine Query-Unterstützung (API-Limitierung)* |
| **ChairmanNotes** | id, text, date, _deleteAfterDate |
| **CustomFields** | *Keine Query-Unterstützung (API-Limitierung)* |
| **CustomFieldCollections** | id, name, orderSequence, position |
| **CustomFilters** | *Keine Query-Unterstützung (API-Limitierung)* |
| **DocumentTemplates** | *Keine Query-Unterstützung (API-Limitierung)* |
| **Events** | id, name, start, end, allDay, isPublic, canceled, locationName |
| **MemberGroups** | *Keine Query-Unterstützung (API-Limitierung)* |
| **Organization** | *Singleton: nur Get und Update* |
| **FileSystemPaths** | *Keine Query-Unterstützung (API-Limitierung)* |
| **Wastebasket** | *Keine Query-Unterstützung (API-Limitierung); kein Create/Update/Delete* |
| **ChatSettings** | *Singleton: nur Get und Update* |
| **Forums** | *Keine Query-Unterstützung (API-Limitierung)* |
| **DosbSports** | *Keine Query-Unterstützung (API-Limitierung)* |
| **LsbSports** | *Keine Query-Unterstützung (API-Limitierung)* |

## Fehlerbehandlung

Das Paket wirft zwei spezialisierte Error-Typen:

### APIError

Nicht-2xx HTTP-Responses der API:

```go
var apiErr *easyvapi.APIError
if errors.As(err, &apiErr) {
	fmt.Printf("API-Fehler %d: %s\n", apiErr.StatusCode, apiErr.Message)
	if apiErr.Detail != "" {
		fmt.Printf("Details: %s\n", apiErr.Detail)
	}
}
```

### RateLimitError

Tritt auf bei HTTP 429 (zu viele Requests):

```go
var rateLimitErr *easyvapi.RateLimitError
if errors.As(err, &rateLimitErr) {
	fmt.Printf("Rate-Limit. Retry nach: %v\n", rateLimitErr.RetryAfter)
	time.Sleep(rateLimitErr.RetryAfter)
}
```

## Best Practices

### 1. Context mit Timeout verwenden

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

members, err := client.Members.ListAll(ctx, nil)
```

### 2. Query-Filter für bessere Performance

```go
// ❌ Schlecht: Alle Felder abrufen
members, _ := client.Members.ListAll(ctx, nil)

// ✅ Besser: Nur benötigte Felder
query := easyvapi.NewQuery().Fields("id", "membershipNumber")
opts := &easyvapi.MemberListOptions{
	ListOptions: easyvapi.ListOptions{Query: query},
}
members, _ := client.Members.ListAll(ctx, opts)
```

### 3. Iterator für große Datenmengen

```go
// ❌ Speicher-intensiv: Alle auf einmal laden
members, _ := client.Members.ListAll(ctx, nil)

// ✅ Speicher-effizient: Lazy-Iteration
iter := client.Members.List(ctx, nil)
for iter.Next() {
	m := iter.Value()
	// Verarbeite einzeln
}
```

### 4. Rate-Limit beachten

Die API hat ein Limit von 100 Requests/Minute. Das Paket drosselt automatisch, aber:

```go
// Pausen zwischen verschiedenen Endpoints einbauen
for i, endpoint := range endpoints {
	if i > 0 {
		time.Sleep(5 * time.Second)
	}
	// Abruf...
}
```

## Lizenz

MIT

## Support

Für Fragen oder Bugs: [GitHub Issues](https://github.com/as27/easyvapi/issues)
