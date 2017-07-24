package apidVerifyApiKey

// Fields related to developer
type DeveloperDetails struct {
	Id string `json:"id,omitempty"`

	UserName string `json:"userName,omitempty"`

	FirstName string `json:"firstName,omitempty"`

	LastName string `json:"lastName,omitempty"`

	Email string `json:"email,omitempty"`

	Status string `json:"status,omitempty"`

	Apps []string `json:"apps,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`

	CreatedBy string `json:"created_by,omitempty"`

	LastmodifiedAt string `json:"lastmodified_at,omitempty"`

	LastmodifiedBy string `json:"lastmodified_by,omitempty"`

	Company string `json:"company,omitempty"`

	// Attributes associated with the developer.
	Attributes []Attribute `json:"attributes,omitempty"`
}
