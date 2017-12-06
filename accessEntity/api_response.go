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
	ApiProduct *ApiProductDetails `json:"apiProduct"`
	// Organization Identifier/Name
	Organization string `json:"organization"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue"`
	// secondary identifier type
	SecondaryIdentifierType string `json:"secondaryIdentifierType"`
	// secondary identifier value
	SecondaryIdentifierValue string `json:"secondaryIdentifierValue"`
}

type AppCredentialSuccessResponse struct {
	// app credential
	AppCredential *AppCredentialDetails `json:"appCredential"`
	// Organization Identifier/Name
	Organization string `json:"organization"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue"`
}

type AppSuccessResponse struct {
	// app
	App *AppDetails `json:"app"`
	// Organization Identifier/Name
	Organization string `json:"organization"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue"`
	// secondary identifier type
	SecondaryIdentifierType string `json:"secondaryIdentifierType"`
	// secondary identifier value
	SecondaryIdentifierValue string `json:"secondaryIdentifierValue"`
}

type CompanyDevelopersSuccessResponse struct {
	// company developers
	CompanyDevelopers []*CompanyDeveloperDetails `json:"companyDevelopers"`
	// Organization Identifier/Name
	Organization string `json:"organization"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue"`
}

type CompanySuccessResponse struct {
	// company
	Company *CompanyDetails `json:"company"`
	// Organization Identifier/Name
	Organization string `json:"organization"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue"`
}

type DeveloperSuccessResponse struct {
	// developer
	Developer *DeveloperDetails `json:"developer"`
	// Organization Identifier/Name
	Organization string `json:"organization"`
	// primary identifier type
	PrimaryIdentifierType string `json:"primaryIdentifierType"`
	// primary identifier value
	PrimaryIdentifierValue string `json:"primaryIdentifierValue"`
}

type ApiProductDetails struct {
	// api proxies
	ApiProxies []string `json:"apiProxies"`
	// api resources
	ApiResources []string `json:"apiResources"`
	// approval type
	ApprovalType string `json:"approvalType"`
	// Attributes associated with the apiproduct.
	Attributes []common.Attribute `json:"attributes"`
	// ISO-8601
	CreatedAt string `json:"createdAt"`
	// created by
	CreatedBy string `json:"createdBy"`
	// description
	Description string `json:"description"`
	// display name
	DisplayName string `json:"displayName"`
	// environments
	Environments []string `json:"environments"`
	// id
	ID string `json:"id"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy"`
	// name
	Name string `json:"name"`
	// quota interval
	QuotaInterval int64 `json:"quotaInterval"`
	// quota limit
	QuotaLimit int64 `json:"quotaLimit"`
	// quota time unit
	QuotaTimeUnit string `json:"quotaTimeUnit"`
	// scopes
	Scopes []string `json:"scopes"`
}

type AppDetails struct {

	// access type
	AccessType string `json:"accessType"`
	// api products
	ApiProducts []string `json:"apiProducts"`
	// app credentials
	AppCredentials []*CredentialDetails `json:"appCredentials"`
	// app family
	AppFamily string `json:"appFamily"`
	// app parent, developer's Id or company's name
	AppParentID string `json:"appParentId"`
	// app parent status
	AppParentStatus string `json:"appParentStatus"`
	// Developer or Company
	AppType string `json:"appType"`
	// Attributes associated with the app.
	Attributes []common.Attribute `json:"attributes"`
	// callback Url
	CallbackUrl string `json:"callbackUrl"`
	// ISO-8601
	CreatedAt string `json:"createdAt"`
	// created by
	CreatedBy string `json:"createdBy"`
	// display name
	DisplayName string `json:"displayName"`
	// id
	Id string `json:"id"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy"`
	// name
	Name string `json:"name"`
	// status
	Status string `json:"status"`
}

type CredentialDetails struct {
	// api product references
	ApiProductReferences []string `json:"apiProductReferences"`
	// app Id
	AppID string `json:"appId"`
	// app status
	AppStatus string `json:"appStatus"`
	// Attributes associated with the client Id.
	Attributes []common.Attribute `json:"attributes"`
	// consumer key
	ConsumerKey string `json:"consumerKey"`
	// consumer secret
	ConsumerSecret string `json:"consumerSecret"`
	// expires at
	ExpiresAt string `json:"expiresAt"`
	// issued at
	IssuedAt string `json:"issuedAt"`
	// method type
	MethodType string `json:"methodType"`
	// scopes
	Scopes []string `json:"scopes"`
	// status
	Status string `json:"status"`
}

/*
type ApiProductReferenceDetails struct {
	// status of the api product
	Status         string `json:"status"`
	// name of the api product
	ApiProduct string `json:"apiProduct"`
}
*/
type AppCredentialDetails struct {
	// app Id
	AppID string `json:"appId"`
	// app name
	AppName string `json:"appName"`
	// Attributes associated with the app credential
	Attributes []common.Attribute `json:"attributes"`
	// consumer key
	ConsumerKey string `json:"consumerKey"`
	// consumer key status
	ConsumerKeyStatus *ConsumerKeyStatusDetails `json:"consumerKeyStatus"`
	// consumer secret
	ConsumerSecret string `json:"consumerSecret"`
	// developer Id
	DeveloperID string `json:"developerId"`
	// redirect uris
	RedirectUris []string `json:"redirectURIs"`
	// scopes
	Scopes []string `json:"scopes"`
	// status
	Status string `json:"status"`
}

type ConsumerKeyStatusDetails struct {
	// app credential
	AppCredential *CredentialDetails `json:"appCredential"`
	// app Id
	AppID string `json:"appId"`
	// app name
	AppName string `json:"appName"`
	// app status
	AppStatus string `json:"appStatus"`
	// app type
	AppType string `json:"appType"`
	// developer Id
	DeveloperID string `json:"developerId"`
	// developer status
	DeveloperStatus string `json:"developerStatus"`
	// is valid key
	IsValidKey bool `json:"isValidKey"`
}

type CompanyDetails struct {

	// apps
	Apps []string `json:"apps"`
	// Attributes associated with the company.
	Attributes []common.Attribute `json:"attributes"`
	// ISO-8601
	CreatedAt string `json:"createdAt"`
	// created by
	CreatedBy string `json:"createdBy"`
	// display name
	DisplayName string `json:"displayName"`
	// id
	ID string `json:"id"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy"`
	// name
	Name string `json:"name"`
	// status
	Status string `json:"status"`
}

type CompanyDeveloperDetails struct {
	// company name
	CompanyName string `json:"companyName"`
	// ISO-8601
	CreatedAt string `json:"createdAt"`
	// created by
	CreatedBy string `json:"createdBy"`
	// developer email
	DeveloperEmail string `json:"developerEmail"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy"`
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
	CreatedAt string `json:"createdAt"`
	// created by
	CreatedBy string `json:"createdBy"`
	// email
	Email string `json:"email"`
	// first name
	FirstName string `json:"firstName"`
	// id
	ID string `json:"id"`
	// ISO-8601
	LastModifiedAt string `json:"lastModifiedAt"`
	// last modified by
	LastModifiedBy string `json:"lastModifiedBy"`
	// last name
	LastName string `json:"lastName"`
	// password
	Password string `json:"password"`
	// status
	Status string `json:"status"`
	// user name
	UserName string `json:"userName"`
}
