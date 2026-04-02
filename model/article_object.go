package model

// ArticleObject represents an article/product in easyVerein (Artikel).
// Articles are used for shop items or event tickets.
type ArticleObject struct {
	// ID is the unique identifier of the article.
	ID int `json:"id"`
	// Name is the display name of the article.
	Name string `json:"name"`
	// Kind classifies the article type.
	Kind string `json:"kind"`
	// Price is the article price.
	Price flexFloat64 `json:"price"`
	// Description is an optional free-text description.
	Description string `json:"description"`
	// Unit is the unit of measure (e.g. "piece", "kg").
	Unit string `json:"unit"`
}

// ArticleObjectCreate holds the fields for creating or updating an article
// via POST / PATCH /article-object.
type ArticleObjectCreate struct {
	// Name is the display name (required for create).
	Name string `json:"name,omitempty"`
	// Kind classifies the article type.
	Kind string `json:"kind,omitempty"`
	// Price is the article price.
	Price float64 `json:"price,omitempty"`
	// Description is an optional free-text description.
	Description string `json:"description,omitempty"`
	// Unit is the unit of measure.
	Unit string `json:"unit,omitempty"`
}
