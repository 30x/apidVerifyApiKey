package accessEntity

import "github.com/apid/apidVerifyApiKey/common"

type ApiProductDetails struct {
	// api proxies
	APIProxies []string `json:"apiProxies,omitempty"`

	// api resources
	APIResources []string `json:"apiResources,omitempty"`

	// approval type
	ApprovalType string `json:"approvalType,omitempty"`

	// Attributes associated with the apiproduct.
	Attributes []common.Attribute `json:"attributes,omitempty"`

	// ISO-8601
	CreatedAt string `json:"createdAt,omitempty"`

	// created by
	CreatedBy string `json:"createdBy,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	// display name
	DisplayName string `json:"displayName,omitempty"`

	// environments
	Environments []string `json:"environments,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`

	// last modified by
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType,omitempty"`

	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue,omitempty"`

	// quota interval
	QuotaInterval int64 `json:"quotaInterval,omitempty"`

	// quota limit
	QuotaLimit int64 `json:"quotaLimit,omitempty"`

	// quota time unit
	QuotaTimeUnit string `json:"quotaTimeUnit,omitempty"`

	// scopes
	Scopes []string `json:"scopes,omitempty"`

	// secondary identifier type
	SecondaryIdentifierType string `json:"secondaryIdentifierType,omitempty"`

	// secondary identifier value
	SecondaryIdentifierValue string `json:"secondaryIdentifierValue,omitempty"`
}
