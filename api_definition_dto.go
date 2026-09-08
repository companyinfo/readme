package readme

import "time"

// APIDefinition models the response from the ReadMe API v2 for API definitions.
//
// See: https://docs.readme.com/main/reference/getapis
type APIDefinition struct {
	ID        string              `json:"id"`
	Version   string              `json:"version"` // branch/version slug
	Title     string              `json:"title"`
	UpdatedAt string              `json:"updated_at,omitzero"`
	CreatedAt time.Time           `json:"created_at"`
	Filename  string              `json:"filename"`
	LegacyId  string              `json:"legacy_id"`
	Source    APIDefinitionSource `json:"source"`
	StatusUrl string              `json:"status_url"`
	Type      string              `json:"type"`
	Upload    APIDefinitionUpload `json:"upload"`
	Uri       string              `json:"uri"`
}

type APIDefinitionSource struct {
	Current  string `json:"current"`
	Original string `json:"original"`
	SyncUrl  string `json:"sync_url"`
}
type APIDefinitionUpload struct {
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	Warnings string `json:"warnings"`
}

// APIDefinitionParams is the request body for creating/updating/validating an API definition.
//
// See:
//   - https://docs.readme.com/main/reference/createapi
//   - https://docs.readme.com/main/reference/updateapi
//   - https://docs.readme.com/main/reference/validateapi
type APIDefinitionParams struct {
	// Schema is the OpenAPI or Swagger specification as a string (JSON or YAML).
	FileName     string `json:"filename,omitzero"`
	Schema       string `json:"schema,omitzero"`
	UploadSource string `json:"upload_source,omitzero"`
	Url          string `json:"url,omitzero"`
}
