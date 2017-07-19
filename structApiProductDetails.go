package apidVerifyApiKey

// Fields related to app
type ApiProductDetails struct {
	Id string `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	DisplayName string `json:"displayName,omitempty"`

	QuotaLimit int64 `json:"quota.limit,omitempty"`

	QuotaInterval int64 `json:"quota.interval,omitempty"`

	QuotaTimeunit int64 `json:"quota.timeunit,omitempty"`

	Status string `json:"status,omitempty"`

	CreatedAt int64 `json:"created_at,omitempty"`

	CreatedBy string `json:"created_by,omitempty"`

	LastmodifiedAt int64 `json:"lastmodified_at,omitempty"`

	LastmodifiedBy string `json:"lastmodified_by,omitempty"`

	Company string `json:"company,omitempty"`

	Environments []string `json:"environments,omitempty"`

	Apiproxies []string `json:"apiproxies,omitempty"`

	// Attributes associated with the apiproduct.
	Attributes []Attribute `json:"attributes,omitempty"`
}
