package apidVerifyApiKey

// Response object for the verification of apikey. Verification of apikey response contains details such as developer-id,developer-email-id, other fields and attributes ; app-id,app-name, other fields and attributes;  apiproduct-name, fields and attributes ;
type VerifyApiKeySuccessResponse struct {
	Self string `json:"self,omitempty"`

	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`

	// Environment Identifier/Name
	Environment string `json:"environment,omitempty"`

	ClientId ClientIdDetails `json:"clientId,omitempty"`

	Developer DeveloperDetails `json:"developer,omitempty"`

	Company CompanyDetails `json:"company,omitempty"`

	App AppDetails `json:"app,omitempty"`

	ApiProduct ApiProductDetails `json:"apiProduct,omitempty"`

	// Identifier of the authorization code. This will be unique for each request.
	Identifier string `json:"identifier,omitempty"`

	Kind string `json:"kind,omitempty"`
}
