package readme

import (
	"errors"
	"strings"
)

// validateSection checks that a CategoryType is one of the accepted
// values ("guides" or "reference"), comparing case-insensitively after
// trimming surrounding whitespace. It returns an error describing the
// allowed values when the input is empty or unknown.
func validateSection(s CategoryType) error {
	switch canonSection(s) {
	case "guide", "reference":
		return nil
	default:
		return errors.New("section must be one of [guides reference]")
	}
}

// canonSection returns the canonical lowercase form of a CategoryType,
// trimmed of surrounding whitespace. It produces the exact spelling
// expected by the ReadMe API in the URL.
func canonSection(s CategoryType) string {
	return strings.ToLower(strings.TrimSpace(string(s)))
}

// CategoryType identifies the section a category belongs to in ReadMe API v2.
//
// In v2 the section is supplied as a path parameter (or in the request body
// when creating a category). The accepted values are "guides" and "reference".
type CategoryType string

const (
	// CategoryTypeReference represents the API Reference section.
	CategoryTypeReference CategoryType = "reference"
	// CategoryTypeGuides represents the Guides (knowledge base) section.
	CategoryTypeGuides CategoryType = "guides"
)

// Category models a category returned by the ReadMe API v2.
//
// See: https://docs.readme.com/main/reference/getcategories-1
type Category struct {
	// Title is the display title of the category and also serves as its unique
	// resource identifier within a (branch, section) pair in API v2.
	Title string `json:"title"`
	// Section is "guides" or "reference".
	Section CategoryType `json:"section,omitempty"`
	// Links contain the `project` URI for the category.
	Links CategoryLinks `json:"links"`
	// URI to this category.
	URI string `json:"uri"`
}

// CategoryLinks models the `links` block on a category.
type CategoryLinks struct {
	Project string `json:"project"`
}

// CategoryParams contains the fields used when creating or updating a category
// via the ReadMe API v2.
//
// Note: in v2 the version/branch is supplied as a path parameter and is no
// longer part of the request body. Validation tags (consumed by
// go-playground/validator) describe the requirements for the **Create**
// endpoint: `title` is required and `section` must be one of "guides" or
// "reference".
type CategoryParams struct {
	// Title is the title of the category.
	Title string `json:"title" validate:"required"`
	// Section is "guides" or "reference". Required on Create; optional on Update.
	Section CategoryType `json:"section,omitempty" validate:"required,oneof=guide reference"`
}
