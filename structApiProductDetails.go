package apidVerifyApiKey

// Fields related to app
type ApiProductDetails struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	DisplayName string `json:"displayName,omitempty"`

	QuotaLimit string `json:"quota.limit,omitempty"`

	QuotaInterval string `json:"quota.interval,omitempty"`

	QuotaTimeunit string `json:"quota.timeunit,omitempty"`

	Status string `json:"status,omitempty"`

	CreatedAt string `json:"created_at,omitempty"`

	CreatedBy string `json:"created_by,omitempty"`

	LastmodifiedAt string `json:"lastmodified_at,omitempty"`

	LastmodifiedBy string `json:"lastmodified_by,omitempty"`

	Company string `json:"company,omitempty"`

	Environments []string `json:"environments,omitempty"`

	Apiproxies []string `json:"apiproxies,omitempty"`

	// Attributes associated with the apiproduct.
	Attributes []Attribute `json:"attributes,omitempty"`
}
