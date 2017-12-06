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
package verifyApiKey

import (
	"errors"
	"github.com/apid/apidVerifyApiKey/common"
)

type ClientIdDetails struct {
	ClientId     string   `json:"clientId,omitempty"`
	ClientSecret string   `json:"clientSecret,omitempty"`
	RedirectURIs []string `json:"redirectURIs,omitempty"`
	Status       string   `json:"status,omitempty"`
	// Attributes associated with the client Id.
	Attributes []common.Attribute `json:"attributes,omitempty"`
}

type ApiProductDetails struct {
	Id             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	DisplayName    string   `json:"displayName,omitempty"`
	QuotaLimit     string   `json:"quota.limit,omitempty"`
	QuotaInterval  int64    `json:"quota.interval,omitempty"`
	QuotaTimeunit  string   `json:"quota.timeunit,omitempty"`
	Status         string   `json:"status,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	CreatedBy      string   `json:"created_by,omitempty"`
	LastmodifiedAt string   `json:"lastmodified_at,omitempty"`
	LastmodifiedBy string   `json:"lastmodified_by,omitempty"`
	Company        string   `json:"company,omitempty"`
	Environments   []string `json:"environments,omitempty"`
	Apiproxies     []string `json:"apiproxies,omitempty"`
	// Attributes associated with the apiproduct.
	Attributes []common.Attribute `json:"attributes,omitempty"`
	Resources  []string           `json:"-"`
}

type AppDetails struct {
	Id             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	AccessType     string   `json:"accessType,omitempty"`
	CallbackUrl    string   `json:"callbackUrl,omitempty"`
	DisplayName    string   `json:"displayName,omitempty"`
	Status         string   `json:"status,omitempty"`
	Apiproducts    []string `json:"apiproducts,omitempty"`
	AppFamily      string   `json:"appFamily,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	CreatedBy      string   `json:"created_by,omitempty"`
	LastmodifiedAt string   `json:"lastmodified_at,omitempty"`
	LastmodifiedBy string   `json:"lastmodified_by,omitempty"`
	Company        string   `json:"company,omitempty"`
	// Attributes associated with the app.
	Attributes []common.Attribute `json:"attributes,omitempty"`
}

type CompanyDetails struct {
	Id             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	DisplayName    string   `json:"displayName,omitempty"`
	Status         string   `json:"status,omitempty"`
	Apps           []string `json:"apps,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	CreatedBy      string   `json:"created_by,omitempty"`
	LastmodifiedAt string   `json:"lastmodified_at,omitempty"`
	LastmodifiedBy string   `json:"lastmodified_by,omitempty"`
	// Attributes associated with the company.
	Attributes []common.Attribute `json:"attributes,omitempty"`
}

type DeveloperDetails struct {
	Id             string   `json:"id,omitempty"`
	UserName       string   `json:"userName,omitempty"`
	FirstName      string   `json:"firstName,omitempty"`
	LastName       string   `json:"lastName,omitempty"`
	Email          string   `json:"email,omitempty"`
	Status         string   `json:"status,omitempty"`
	Apps           []string `json:"apps,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	CreatedBy      string   `json:"created_by,omitempty"`
	LastmodifiedAt string   `json:"lastmodified_at,omitempty"`
	LastmodifiedBy string   `json:"lastmodified_by,omitempty"`
	Company        string   `json:"company,omitempty"`
	// Attributes associated with the developer.
	Attributes []common.Attribute `json:"attributes,omitempty"`
}

type VerifyApiKeyRequest struct {
	Action           string `json:"action"`
	Key              string `json:"key"`
	UriPath          string `json:"uriPath"`
	OrganizationName string `json:"organizationName"`
	EnvironmentName  string `json:"environmentName"`
	ApiProxyName     string `json:"apiProxyName"`
	// when this flag is false, authentication of key and authorization for uripath is done and authorization for apiproxies and environments is skipped. Default is true.
	ValidateAgainstApiProxiesAndEnvs bool `json:"validateAgainstApiProxiesAndEnvs,omitempty"`
}

func (v *VerifyApiKeyRequest) validate() (bool, error) {
	var validationMsg string

	if v.Action == "" {
		validationMsg += " action"
	}

	if v.Key == "" {
		validationMsg += " key"
	}
	if v.OrganizationName == "" {
		validationMsg += " organizationName"
	}
	if v.UriPath == "" {
		validationMsg += " uriPath"
	}
	if v.ValidateAgainstApiProxiesAndEnvs {
		if v.ApiProxyName == "" {
			validationMsg += " apiProxyName"
		}
		if v.EnvironmentName == "" {
			validationMsg += " environmentName"
		}
	}

	if validationMsg != "" {
		validationMsg = "Missing mandatory fields in the request :" + validationMsg
		return false, errors.New(validationMsg)
	}
	return true, nil
}

type VerifyApiKeySuccessResponse struct {
	Self string `json:"self,omitempty"`
	// Organization Identifier/Name
	Organization string `json:"organization,omitempty"`
	// Environment Identifier/Name
	Environment string            `json:"environment,omitempty"`
	ClientId    ClientIdDetails   `json:"clientId,omitempty"`
	Developer   DeveloperDetails  `json:"developer,omitempty"`
	Company     CompanyDetails    `json:"company,omitempty"`
	App         AppDetails        `json:"app,omitempty"`
	ApiProduct  ApiProductDetails `json:"apiProduct,omitempty"`
	// Identifier of the authorization code. This will be unique for each request.
	Identifier string `json:"identifier,omitempty"`
	Kind       string `json:"kind,omitempty"`
}

type VerifyApiKeyRequestResponseDataWrapper struct {
	verifyApiKeyRequest         VerifyApiKeyRequest
	verifyApiKeySuccessResponse VerifyApiKeySuccessResponse
	tempDeveloperDetails        DeveloperDetails
	apiProducts                 []ApiProductDetails
	ctype                       string
	tenant_id                   string
}
