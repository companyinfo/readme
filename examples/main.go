// Example flow demonstrating how to use the readme client against the
// ReadMe API v2. It walks through:
//
//  1. Constructing a Client with an API key (and optionally a custom base URL).
//  2. Listing categories for a branch + section.
//  3. Looking up a single category by title.
//  4. Creating, fetching, updating and deleting a Guide page.
//  5. Creating, fetching, updating and deleting a Reference page
//     (including the Reference-specific `api` block: method + path).
//
// Configuration is read from environment variables so the example stays
// runnable without code changes:
//
//	README_API_KEY   (required) — your ReadMe API key
//	README_BASE_URL  (optional) — defaults to https://api.readme.com/v2
//	README_BRANCH    (optional) — defaults to "v0.0"
//	README_RUN_WRITE (optional) — set to "1" to run the create/update/delete
//	                              flows. Defaults to read-only.
//
// Run with:
//
//	README_API_KEY=rdme_xxx README_RUN_WRITE=1 go run ./examples
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.companyinfo.dev/go-readmeio"
)

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	apiKey := "***********************"

	baseURL := envOr("README_BASE_URL", "https://api.readme.com/v2")
	branch := envOr("README_BRANCH", "v0.0")

	// set this to true if you want to run the examples that actually create pages on ReadMe
	runWrite := false

	client, err := readme.New(apiKey, readme.WithBaseURL(baseURL))
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// ---------------------------------------------------------------------
	// 1. Categories — read-only listing for both sections.
	// ---------------------------------------------------------------------
	guideCats, err := client.Categories.List(ctx, branch, readme.CategoryTypeGuides)
	if err != nil {
		log.Fatalf("categories.List(guides): %v", err)
	}
	fmt.Printf("Fetched %d guide categories on branch %q\n", len(guideCats), branch)
	for i, c := range guideCats {
		fmt.Printf("  %d. %s (uri=%s)\n", i+1, c.Title, c.URI)
	}

	refCats, err := client.Categories.List(ctx, branch, readme.CategoryTypeReference)
	if err != nil {
		log.Fatalf("categories.List(reference): %v", err)
	}
	fmt.Printf("Fetched %d reference categories on branch %q\n", len(refCats), branch)

	// 1b. Look up a single category by title (uses the first guide category
	//     returned above, if any).
	if len(guideCats) > 0 {
		first := guideCats[0]
		cat, err := client.Categories.GetByTitle(ctx, branch, "guides", first.Title)
		if err != nil {
			log.Fatalf("categories.GetByTitle(%q): %v", first.Title, err)
		}
		fmt.Printf("Resolved category by title: %+v\n", cat)
	}

	if !runWrite {
		fmt.Println("README_RUN_WRITE != 1 — skipping create/update/delete flow.")
		return
	}

	// ---------------------------------------------------------------------
	// 2. Guide CRUD flow.
	// ---------------------------------------------------------------------
	if len(guideCats) == 0 {
		log.Fatal("no guide categories on this branch — cannot create a guide")
	}
	guideCatURI := guideCats[0].URI

	guide, err := client.Guides.Create(ctx, branch, readme.GuideParams{
		Title:    "Hello from readme",
		Category: &readme.ResourceRef{URI: guideCatURI},
		Content: &readme.GuideContent{
			Body:    "# Hello\n\nCreated by the readme example.",
			Excerpt: "Example guide created via readme.",
		},
	})
	if err != nil {
		log.Fatalf("guides.Create: %v", err)
	}
	fmt.Printf("Created guide: title=%q slug=%q\n", guide.Title, guide.Slug)

	got, err := client.Guides.Get(ctx, branch, guide.Slug)
	if err != nil {
		log.Fatalf("guides.Get: %v", err)
	}
	fmt.Printf("Fetched guide: title=%q\n", got.Title)

	updated, err := client.Guides.Update(ctx, branch, guide.Slug, readme.GuideParams{
		Title: "Hello from readme (updated)",
	})
	if err != nil {
		log.Fatalf("guides.Update: %v", err)
	}
	fmt.Printf("Updated guide title to: %q\n", updated.Title)

	if err := client.Guides.Delete(ctx, branch, guide.Slug); err != nil {
		log.Fatalf("guides.Delete: %v", err)
	}
	fmt.Printf("Deleted guide %q\n", guide.Slug)

	// ---------------------------------------------------------------------
	// 3. Reference CRUD flow (note the Reference-specific `api` block).
	// ---------------------------------------------------------------------
	if len(refCats) == 0 {
		log.Println("no reference categories on this branch — skipping reference flow")
		return
	}
	refCatURI := refCats[0].URI

	ref, err := client.Reference.Create(ctx, branch, readme.ReferenceParams{
		Title:    "Get pets",
		Category: &readme.ResourceRef{URI: refCatURI},
		API: &readme.ReferenceAPI{
			Method: "get",
			Path:   "/pets",
		},
		Content: &readme.GuideContent{
			Excerpt: "Returns all pets from the system.",
		},
	})
	if err != nil {
		log.Fatalf("reference.Create: %v", err)
	}
	fmt.Printf("Created reference: title=%q slug=%q api=%+v\n", ref.Title, ref.Slug, ref.API)

	gotRef, err := client.Reference.Get(ctx, branch, ref.Slug)
	if err != nil {
		log.Fatalf("reference.Get: %v", err)
	}
	fmt.Printf("Fetched reference: title=%q\n", gotRef.Title)

	updatedRef, err := client.Reference.Update(ctx, branch, ref.Slug, readme.ReferenceParams{
		API: &readme.ReferenceAPI{Method: "post", Path: "/pets"},
	})
	if err != nil {
		log.Fatalf("reference.Update: %v", err)
	}
	fmt.Printf("Updated reference api block to: %+v\n", updatedRef.API)

	if err := client.Reference.Delete(ctx, branch, ref.Slug); err != nil {
		log.Fatalf("reference.Delete: %v", err)
	}
	fmt.Printf("Deleted reference %q\n", ref.Slug)

	// ---------------------------------------------------------------------
	// 4. API Definition flow.
	// ---------------------------------------------------------------------
	spec := `{"openapi":"3.0.0","info":{"title":"Example API","version":"1.0.0"},"paths":{}}`
	filename := "example.json"

	fmt.Println("Validating API definition...")
	warns, err := client.APIDefinitions.Validate(ctx, branch, readme.APIDefinitionParams{
		Schema:   spec,
		FileName: filename,
	})
	if err != nil {
		log.Fatalf("apiDefinitions.Validate: %v", err)
	}
	fmt.Printf("Validation result: %s\n", warns)

	fmt.Println("Creating API definition...")
	err = client.APIDefinitions.Create(ctx, branch, readme.APIDefinitionParams{
		Schema:   spec,
		FileName: filename,
	})
	if err != nil {
		log.Fatalf("apiDefinitions.Create: %v", err)
	}
	fmt.Printf("Created API definition: %s\n", filename)

	apiDef, err := client.APIDefinitions.Get(ctx, branch, filename)
	if err != nil {
		log.Fatalf("apiDefinitions.Get: %v", err)
	}
	fmt.Printf("Fetched API definition: title=%q id=%q\n", apiDef.Title, apiDef.ID)

	fmt.Println("Updating API definition...")
	err = client.APIDefinitions.Update(ctx, branch, filename, readme.APIDefinitionParams{
		Schema:   `{"openapi":"3.0.0","info":{"title":"Example API (updated)","version":"1.0.0"},"paths":{}}`,
		FileName: filename,
	})
	if err != nil {
		log.Fatalf("apiDefinitions.Update: %v", err)
	}
	fmt.Println("Updated API definition.")

	if err := client.APIDefinitions.Delete(ctx, branch, filename); err != nil {
		log.Fatalf("apiDefinitions.Delete: %v", err)
	}
	fmt.Printf("Deleted API definition %q\n", filename)
}
