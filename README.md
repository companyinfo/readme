# go-readmeio

A Go client for the ReadMe API v2 that provides typed DTOs and convenient services for common Docs operations:
- Categories (create, list, get-by-title, update, delete)
- Guides (create, get, update, delete)
- References (create, get, update, delete)
- API Definitions (create, get, update, delete, validate)

This SDK wraps JSON:API-style payloads with idiomatic Go types and adds light client-side validation and structured error handling.

## Features

- Simple client setup with pluggable HTTP client and base URL
- Strongly-typed request/response models for categories, guides, references, and API definitions
- Consistent CRUD service interfaces
- Granular error details via APIError and APIFieldError
- Small, dependency-light surface built on resty

## Installation

```
go get go.companyinfo.dev/go-readmeio@latest
```

Requires Go 1.21+.

## Quickstart

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.companyinfo.dev/go-readmeio"
)

func main() {
	apiKey := os.Getenv("README_API_KEY")
	if apiKey == "" {
		log.Fatal("set README_API_KEY")
	}

	client, err := readme.New(
		apiKey,
		// Optional: override base URL or HTTP client.
		// readme.WithBaseURL("https://api.readme.com/v2"),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	branch := "v0.0" // Docs branch (a.k.a. version)
	section := readme.CategoryTypeReference

	// List categories
	cats, err := client.Categories.List(ctx, branch, section)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d %s categories\n", len(cats), section)

	// Create a guide (minimal)
	if len(cats) > 0 {
		g, err := client.Guides.Create(ctx, branch, readme.GuideParams{
			Title:    "Hello from readme",
			Category: &readme.ResourceRef{URI: cats[0].URI},
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created guide %q (slug=%s)\n", g.Title, g.Slug)
	}
}
```

## Client configuration

- WithBaseURL(url string): override API base URL (default https://api.readme.com/v2)
- WithHTTPClient(hc *resty.Client): supply a custom resty client (timeouts, retries, logging, etc.)

Authentication: pass your ReadMe API key to readme.New. The client sends Authorization: Bearer <key> on all requests.

## Services

- Categories: Create, List (by section), GetByTitle, Update, Delete
- Guides: Create, Get, Update, Delete
- References: Create, Get, Update, Delete
- APIDefinitions: Create, Get, Update, Delete, Validate

All operations require a branch (version) path parameter (e.g., v1.0, v0.0).

### Minimal examples

Create a category:
```go
created, err := client.Categories.Create(ctx, "v0.0", readme.CategoryParams{
  Title:   "Payments",
  Section: readme.CategoryTypeReference,
})
```

Create a guide:
```go
g, err := client.Guides.Create(ctx, "v0.0", readme.GuideParams{
  Title:    "Getting Started",
  Category: &readme.ResourceRef{URI: "<category-uri>"},
})
```

Create a reference:
```go
r, err := client.Reference.Create(ctx, "v0.0", readme.ReferenceParams{
  Title:    "List Pets",
  Type:     "basic",
  Category: &readme.ResourceRef{URI: "<category-uri>"},
  API:      &readme.ReferenceAPI{Method: "get", Path: "/pets", Source: "api"},
})
```

Upload an API definition:
```go
err := client.APIDefinitions.Create(ctx, "v0.0", readme.APIDefinitionParams{
  FileName: "petstore.json",
  Schema:   `{"openapi":"3.0.0", ...}`,
})
```

## Error handling

API errors are returned as APIError with optional field-level errors:

```go
if err != nil {
  var apiErr *readme.APIError
  if errors.As(err, &apiErr) {
    log.Printf("request failed: %s", apiErr.Error())
  } else {
    log.Fatal(err)
  }
}
```

Example message:
```
422 We encountered validation errors while processing your input.: The JSON you sent was the right format, but had data our endpoint couldn't process. [category.uri: The supplied category URI must be of a category within the guide section of your docs.]
```

## Example program

An example is included to exercise endpoints:

```
go run ./examples
```

Adjust constants in examples/main.go to try different flows.

## Development

- Build: `go build ./...`
- Test: `go test ./...`
- Lint: run your preferred linter (e.g., golangci-lint) if configured

Project layout highlights:
- readme Client with pluggable options
- Category, Guide, Reference services and DTOs
- Centralized API error types in error.go

## Security

- Keep your ReadMe API key secret; prefer environment variables over hardcoding
- Avoid committing credentials; consider using a secrets manager for CI/CD

## Links

- ReadMe API v2 docs: https://docs.readme.com/main/
- Guides API: https://docs.readme.com/main/reference/getguide
- References API: https://docs.readme.com/main/reference/getreference
- Categories API: https://docs.readme.com/main/reference/getcategories-1

## Status

Early-stage, subject to change. Feedback and contributions welcome.
