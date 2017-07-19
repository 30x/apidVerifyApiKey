package apidVerifyApiKey

// Fields related to company
type CompanyDetails struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	DisplayName string `json:"displayName,omitempty"`

	Status string `json:"status,omitempty"`

	Apps []string `json:"apps,omitempty"`

	CreatedAt int64 `json:"created_at,omitempty"`

	CreatedBy string `json:"created_by,omitempty"`

	LastmodifiedAt int64 `json:"lastmodified_at,omitempty"`

	LastmodifiedBy string `json:"lastmodified_by,omitempty"`

	// Attributes associated with the company.
	Attributes []Attribute `json:"attributes,omitempty"`
}
