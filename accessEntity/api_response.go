// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package accessEntity

import "github.com/apid/apidVerifyApiKey/common"

type ApiProductSuccessResponse struct {
	// api product
	ApiProduct *ApiProductDetails `json:"apiProduct,omitempty"`
	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType,omitempty"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue,omitempty"`
	// secondary identifier type
	SecondaryIdentifierType string `json:"secondaryIdentifierType,omitempty"`
	// secondary identifier value
	SecondaryIdentifierValue string `json:"secondaryIdentifierValue,omitempty"`
}

type AppCredentialSuccessResponse struct {
	// app credential
	AppCredential *AppCredentialDetails `json:"appCredential,omitempty"`
	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType,omitempty"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue,omitempty"`
}

type AppSuccessResponse struct {
	// app
	App *AppDetails `json:"app,omitempty"`
	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType,omitempty"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue,omitempty"`
	// secondary identifier type
	SecondaryIdentifierType string `json:"secondaryIdentifierType,omitempty"`
	// secondary identifier value
	SecondaryIdentifierValue string `json:"secondaryIdentifierValue,omitempty"`
}

type CompanyDevelopersSuccessResponse struct {
	// company developers
	CompanyDevelopers []*CompanyDeveloperDetails `json:"companyDevelopers"`
	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType,omitempty"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue,omitempty"`
}

type CompanySuccessResponse struct {
	// company
	Company *CompanyDetails `json:"company,omitempty"`
	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType,omitempty"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue,omitempty"`
}

type DeveloperSuccessResponse struct {
	// developer
	Developer *DeveloperDetails `json:"developer,omitempty"`
	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType,omitempty"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue,omitempty"`
}

type ApiProductDetails struct {
	// api proxies
	ApiProxies []string `json:"apiProxies,omitempty"`
	// api resources
	ApiResources []string `json:"apiResources,omitempty"`
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
	// quota interval
	QuotaInterval int64 `json:"quotaInterval,omitempty"`
	// quota limit
	QuotaLimit int64 `json:"quotaLimit,omitempty"`
	// quota time unit
	QuotaTimeUnit string `json:"quotaTimeUnit,omitempty"`
	// scopes
	Scopes []string `json:"scopes,omitempty"`
}

type AppDetails struct {

	// access type
	AccessType string `json:"accessType,omitempty"`
	// api products
	ApiProducts []string `json:"apiProducts"`
	// app credentials
	AppCredentials []*CredentialDetails `json:"appCredentials"`
	// app family
	AppFamily string `json:"appFamily,omitempty"`
	// app parent Id
	AppParentID string `json:"appParentId,omitempty"`
	// app parent status
	AppParentStatus string `json:"appParentStatus,omitempty"`
	// Developer or Company
	AppType string `json:"appType,omitempty"`
	// Attributes associated with the app.
	Attributes []common.Attribute `json:"attributes"`
	// callback Url
	CallbackUrl string `json:"callbackUrl,omitempty"`
	// ISO-8601
	CreatedAt string `json:"createdAt,omitempty"`
	// created by
	CreatedBy string `json:"createdBy,omitempty"`
	// display name
	DisplayName string `json:"displayName,omitempty"`
	// id
	Id string `json:"id,omitempty"`
	// key expires in
	KeyExpiresIn string `json:"keyExpiresIn,omitempty"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`
	// name
	Name string `json:"name,omitempty"`
	// scopes
	Scopes []string `json:"scopes"`
	// status
	Status string `json:"status,omitempty"`
}

type CredentialDetails struct {
	// api product references
	ApiProductReferences []string `json:"apiProductReferences"`
	// app Id
	AppID string `json:"appId,omitempty"`
	// app status
	AppStatus string `json:"appStatus,omitempty"`
	// Attributes associated with the client Id.
	Attributes []common.Attribute `json:"attributes"`
	// consumer key
	ConsumerKey string `json:"consumerKey,omitempty"`
	// consumer secret
	ConsumerSecret string `json:"consumerSecret,omitempty"`
	// expires at
	ExpiresAt string `json:"expiresAt,omitempty"`
	// issued at
	IssuedAt string `json:"issuedAt,omitempty"`
	// method type
	MethodType string `json:"methodType,omitempty"`
	// scopes
	Scopes []string `json:"scopes"`
	// status
	Status string `json:"status,omitempty"`
}

/*
type ApiProductReferenceDetails struct {
	// status of the api product
	Status         string `json:"status,omitempty"`
	// name of the api product
	ApiProduct string `json:"apiProduct,omitempty"`
}
*/
type AppCredentialDetails struct {
	// app Id
	AppID string `json:"appId,omitempty"`
	// app name
	AppName string `json:"appName,omitempty"`
	// Attributes associated with the app credential
	Attributes []common.Attribute `json:"attributes"`
	// consumer key
	ConsumerKey string `json:"consumerKey,omitempty"`
	// consumer key status
	ConsumerKeyStatus *ConsumerKeyStatusDetails `json:"consumerKeyStatus,omitempty"`
	// consumer secret
	ConsumerSecret string `json:"consumerSecret,omitempty"`
	// developer Id
	DeveloperID string `json:"developerId,omitempty"`
	// redirect uris
	RedirectUris []string `json:"redirectURIs"`
	// scopes
	Scopes []string `json:"scopes"`
	// status
	Status string `json:"status,omitempty"`
}

type ConsumerKeyStatusDetails struct {
	// app credential
	AppCredential *CredentialDetails `json:"appCredential,omitempty"`
	// app Id
	AppID string `json:"appId,omitempty"`
	// app name
	AppName string `json:"appName,omitempty"`
	// app status
	AppStatus string `json:"appStatus,omitempty"`
	// app type
	AppType string `json:"appType,omitempty"`
	// developer Id
	DeveloperID string `json:"developerId,omitempty"`
	// developer status
	DeveloperStatus string `json:"developerStatus,omitempty"`
	// is valid key
	IsValidKey bool `json:"isValidKey,omitempty"`
}

type CompanyDetails struct {

	// apps
	Apps []string `json:"apps"`
	// Attributes associated with the company.
	Attributes []common.Attribute `json:"attributes"`
	// ISO-8601
	CreatedAt string `json:"createdAt,omitempty"`
	// created by
	CreatedBy string `json:"createdBy,omitempty"`
	// display name
	DisplayName string `json:"displayName,omitempty"`
	// id
	ID string `json:"id,omitempty"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`
	// name
	Name string `json:"name,omitempty"`
	// status
	Status string `json:"status,omitempty"`
}

type CompanyDeveloperDetails struct {
	// company name
	CompanyName string `json:"companyName,omitempty"`
	// ISO-8601
	CreatedAt string `json:"createdAt,omitempty"`
	// created by
	CreatedBy string `json:"createdBy,omitempty"`
	// developer email
	DeveloperEmail string `json:"developerEmail,omitempty"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`
	// roles
	Roles []string `json:"roles"`
}

type DeveloperDetails struct {
	// apps
	Apps []string `json:"apps"`
	// Attributes associated with the developer.
	Attributes []common.Attribute `json:"attributes"`
	// companies
	Companies []string `json:"companies"`
	// ISO-8601
	CreatedAt string `json:"createdAt,omitempty"`
	// created by
	CreatedBy string `json:"createdBy,omitempty"`
	// email
	Email string `json:"email,omitempty"`
	// first name
	FirstName string `json:"firstName,omitempty"`
	// id
	ID string `json:"id,omitempty"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt,omitempty"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy,omitempty"`
	// last name
	LastName string `json:"lastName,omitempty"`
	// password
	Password string `json:"password,omitempty"`
	// status
	Status string `json:"status,omitempty"`
	// user name
	UserName string `json:"userName,omitempty"`
}
