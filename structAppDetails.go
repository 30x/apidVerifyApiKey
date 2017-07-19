package apidVerifyApiKey

// Fields related to app
type AppDetails struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	AccessType string `json:"accessType,omitempty"`

	CallbackUrl string `json:"callbackUrl,omitempty"`

	DisplayName string `json:"displayName,omitempty"`

	Status string `json:"status,omitempty"`

	Apiproducts []string `json:"apiproducts,omitempty"`

	AppFamily string `json:"appFamily,omitempty"`

	CreatedAt int64 `json:"created_at,omitempty"`

	CreatedBy string `json:"created_by,omitempty"`

	LastmodifiedAt int64 `json:"lastmodified_at,omitempty"`

	LastmodifiedBy string `json:"lastmodified_by,omitempty"`

	Company string `json:"company,omitempty"`

	// Attributes associated with the app.
	Attributes []Attribute `json:"attributes,omitempty"`
}
