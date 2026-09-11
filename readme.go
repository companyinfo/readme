// Package readme provides utilities for generating and updating README files.
package readme

import (
	"context"
	"errors"

	"github.com/go-resty/resty/v2"
)

const (
	// defaultBaseURL is the hostname for all ReadMe API v2 requests.
	// See: https://api.readme.com/v2
	defaultBaseURL = "https://api.readme.com/v2"
)

// Client is a minimal HTTP client for the ReadMe API v2. It uses Bearer
// authentication with your API key and a configurable base URL and HTTP client.
type Client struct {
	// HTTPClient is the underlying HTTP client used for requests.
	HTTPClient *resty.Client
	// BaseURL is the API base URL. Defaults to https://api.readme.com/v2.
	BaseURL string
	// ApiKey is the ReadMe API key used for Bearer authentication.
	ApiKey string
	// Categories is the service for accessing categories.
	Categories CategoryService
	// Guides is the service for accessing guides.
	Guides GuideService
	// Reference is the service for accessing reference (API reference) pages.
	Reference ReferenceService
	// APIDefinitions is the service for accessing API definitions.
	APIDefinitions APIDefinitionService
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.BaseURL = baseURL
		}
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *resty.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.HTTPClient = hc
		}
	}
}

// New creates a new ReadMe API v2 client using Bearer authentication.
// The provided ApiKey must be a valid ReadMe API key.
// By default, the client uses a resty.Client and is configured to talk to https://api.readme.com/v2.
// You can customize these with options.
func New(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("ApiKey is required")
	}

	c := &Client{
		HTTPClient: resty.New(),
		BaseURL:    defaultBaseURL,
		ApiKey:     apiKey,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(c)
		}
	}

	// Apply final configuration to the underlying Resty client.
	c.HTTPClient.
		SetBaseURL(c.BaseURL).
		SetAuthToken(apiKey).
		SetHeader("Accept", "application/json")

	// Initialize services
	c.Categories = NewCategoryClient(c)
	c.Guides = NewGuideClient(c)
	c.Reference = NewReferenceClient(c)
	c.APIDefinitions = NewAPIDefinitionClient(c)

	return c, nil
}

// NewRequest creates a new Resty request with the provided context applied.
// Services should use this helper to ensure consistent request setup.
func (c *Client) NewRequest(ctx context.Context) *resty.Request {
	return c.HTTPClient.R().SetContext(ctx)
}
