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

package apidVerifyApiKey

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type sucResponseDetail struct {
	Key             string `json:"key"`
	ExpiresAt       int64  `json:"expiresAt"`
	IssuedAt        string `json:"issuedAt"`
	Status          string `json:"status"`
	Type            string `json:"cType"`
	RedirectionURIs string `json:"redirectionURIs"`
	AppId           string `json:"appId"`
	AppName         string `json:"appName"`
}

type errResultDetail struct {
	ErrorCode string `json:"errorCode"`
	Reason    string `json:"reason"`
}

type kmsResponseSuccess struct {
	RspInfo sucResponseDetail `json:"result"`
	Type    string            `json:"type"`
}

type kmsResponseFail struct {
	ErrInfo errResultDetail `json:"result"`
	Type    string          `json:"type"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	kind  string `json:"kind"`
}

type ClientIdDetails struct {
	ClientId     string      `json:"clientId"`
	ClientSecret string      `json:"clientSecret"`
	RedirectUris []string    `json:"redirectUris"`
	Status       string      `json:"status"`
	Attributes   []Attribute `json:"attributes"`
}

type DeveloperDetails struct {
	Id             string      `json:"id"`
	UserName       string      `json:"userName"`
	FirstName      string      `json:"firstName"`
	LastName       string      `json:"lastName"`
	Email          string      `json:"email"`
	Status         string      `json:"status"`
	Apps           []string    `json:"apps"`
	CreatedAt      int         `json:"createdAt"`
	CreatedBy      string      `json:"createdBy"`
	LastModifiedAt int         `json:"lastModifiedAt"`
	LastModifiedBy string      `json:"lastModifiedBy"`
	Company        string      `json:"company"`
	Attributes     []Attribute `json:"attributes"`
}

type CompanyDetails struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	DisplayName    string      `json:"displayName"`
	Status         string      `json:"status"`
	Apps           []string    `json:"apps"`
	CreatedAt      int         `json:"createdAt"`
	CreatedBy      string      `json:"createdBy"`
	LastModifiedAt int         `json:"lastModifiedAt"`
	LastModifiedBy string      `json:"lastModifiedBy"`
	Attributes     []Attribute `json:"attributes"`
}

type AppDetails struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	AccessType     string      `json:"accessType"`
	CallbackUrl    string      `json:"callbackUrl"`
	DisplayName    string      `json:"displayName"`
	Status         string      `json:"status"`
	ApiProducts    []string    `json:"apiProducts"`
	AppFamily      string      `json:"appFamily"`
	CreatedAt      int         `json:"createdAt"`
	CreatedBy      string      `json:"createdBy"`
	LastModifiedAt int         `json:"lastModifiedAt"`
	LastModifiedBy string      `json:"lastModifiedBy"`
	Company        string      `json:"company"`
	Attributes     []Attribute `json:"attributes"`
}

type ApiProductDetails struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	DisplayName    string      `json:"displayName"`
	QuotaLimit     int         `json:"quotaLimit"`
	QuotaInterval  int         `json:"quotaInterval"`
	QuotaTimeUnit  int         `json:"quotaTimeUnit"`
	Status         string      `json:"status"`
	CreatedAt      int         `json:"createdAt"`
	CreatedBy      string      `json:"createdBy"`
	LastModifiedAt int         `json:"lastModifiedAt"`
	LastModifiedBy string      `json:"lastModifiedBy"`
	Company        string      `json:"company"`
	Environments   []string    `json:"environments"`
	ApiProxies     []string    `json:"apiProxies"`
	Attributes     []Attribute `json:"attributes"`
}

// handle client API
func handleRequest(w http.ResponseWriter, r *http.Request) {

	db := getDB()
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("initializing"))
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse form"))
		return
	}

	f := r.Form
	elems := []string{"action", "key", "uriPath", "organizationName", "environmentName", "apiProxyName"}
	for _, elem := range elems {
		if _, ok := f[elem]; !ok {
			log.Debug("Input params Incomplete: " + elem)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Missing element: %s", elem)))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := verifyAPIKey(f)
	if err != nil {
		log.Errorf("error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Debugf("handleVerifyAPIKey result %s", b)
	w.Write(b)
}

// returns []byte to be written to client
func verifyAPIKey(f url.Values) ([]byte, error) {

	action := f.Get("action")
	key := f.Get("key")
	path := f.Get("uriPath")
	organizationName := f.Get("organizationName")
	environmentName := f.Get("environmentName")
	apiProxyName := f.Get("apiProxyName")
	validateAgainstApiProxiesAndEnvs := f.Get("validateAgainstApiProxiesAndEnvs")
	if action != "verify" {

		reason := "action must be 'verify'"
		errorCode := "INCORRECT_USER_INPUT"
		return errorResponse(reason, errorCode)
	}

	db := getDB()

	// DANGER: This relies on an external TABLE - EDGEX_DATA_SCOPE is maintained by apidApigeeSync
	var env, tenantId string
	error := db.QueryRow("SELECT env, scope FROM EDGEX_DATA_SCOPE WHERE id = ?;", scopeuuid).Scan(&env, &tenantId)

	switch {
	case error == sql.ErrNoRows:
		log.Error("verifyAPIKey: sql.ErrNoRows")
		reason := "ENV Validation Failed"
		errorCode := "ENV_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)
	case error != nil:
		reason := error.Error()
		errorCode := "SEARCH_INTERNAL_ERROR"
		return errorResponse(reason, errorCode)
	}

	log.Debug("Found tenant_id='", tenantId, "' with env='", env, "' for scopeuuid='", scopeuuid, "'")

	sSql := `
		SELECT
			ap.api_resources, 
			ap.environments, 
			c.issued_at,
			c.status,
			a.callback_url,
			ad.email,
			ad.id,
			"developer" as ctype
		FROM
			KMS_APP_CREDENTIAL AS c
			INNER JOIN KMS_APP AS a ON c.app_id = a.id
			INNER JOIN KMS_DEVELOPER AS ad
				ON ad.id = a.developer_id
			INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
				ON mp.appcred_id = c.id 
			INNER JOIN KMS_API_PRODUCT as ap ON ap.id = mp.apiprdt_id
		WHERE (UPPER(ad.status) = 'ACTIVE' 
			AND mp.apiprdt_id = ap.id 
			AND mp.app_id = a.id
			AND mp.appcred_id = c.id 
			AND UPPER(mp.status) = 'APPROVED' 
			AND UPPER(a.status) = 'APPROVED'
			AND c.id = $1 
			AND c.tenant_id = $2)
		UNION ALL
		SELECT
			ap.api_resources,
			ap.environments,
			c.issued_at,
			c.status,
			a.callback_url,
			ad.name,
			ad.id,
			"company" as ctype
		FROM
			KMS_APP_CREDENTIAL AS c
			INNER JOIN KMS_APP AS a ON c.app_id = a.id
			INNER JOIN KMS_COMPANY AS ad
				ON ad.id = a.company_id
			INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
				ON mp.appcred_id = c.id
			INNER JOIN KMS_API_PRODUCT as ap ON ap.id = mp.apiprdt_id
		WHERE (UPPER(ad.status) = 'ACTIVE'
			AND mp.apiprdt_id = ap.id
			AND mp.app_id = a.id
			AND mp.appcred_id = c.id
			AND UPPER(mp.status) = 'APPROVED'
			AND UPPER(a.status) = 'APPROVED'
			AND c.id = $1
			AND c.tenant_id = $2)
	;`

	/* these fields need to be nullable types for scanning.  This is because when using json snapshots,
	   and therefore being responsible for inserts, we were able to default everything to be not null.  With
	   sqlite snapshots, we are not necessarily guaranteed that
	*/
	var status, redirectionURIs, appName, appId, resName, resEnv, issuedAt, cType sql.NullString
	err := db.QueryRow(sSql, key, tenantId).Scan(&resName, &resEnv, &issuedAt, &status,
		&redirectionURIs, &appName, &appId, &cType)
	switch {
	case err == sql.ErrNoRows:
		reason := "API Key verify failed for (" + key + ", " + scopeuuid + ", " + path + ")"
		errorCode := "REQ_ENTRY_NOT_FOUND"
		return errorResponse(reason, errorCode)

	case err != nil:
		reason := err.Error()
		errorCode := "SEARCH_INTERNAL_ERROR"
		return errorResponse(reason, errorCode)
	}

	/*
	 * Perform all validations related to the Query made with the data
	 * we just retrieved
	 */
	result := validatePath(resName.String, path)
	if result == false {
		reason := "Path Validation Failed (" + resName.String + " vs " + path + ")"
		errorCode := "PATH_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)

	}

	/* Verify if the ENV matches */
	result = validateEnv(resEnv.String, env)
	if result == false {
		reason := "ENV Validation Failed (" + resEnv.String + " vs " + env + ")"
		errorCode := "ENV_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)
	}

	var expiresAt int64 = -1
	resp := kmsResponseSuccess{
		Type: "APIKeyContext",
		RspInfo: sucResponseDetail{
			Key:             key,
			ExpiresAt:       expiresAt,
			IssuedAt:        issuedAt.String,
			Status:          status.String,
			RedirectionURIs: redirectionURIs.String,
			Type:            cType.String,
			AppId:           appId.String,
			AppName:         appName.String},
	}
	return json.Marshal(resp)
}

func errorResponse(reason, errorCode string) ([]byte, error) {
	if errorCode == "SEARCH_INTERNAL_ERROR" {
		log.Error(reason)
	} else {
		log.Debug(reason)
	}
	resp := kmsResponseFail{
		Type: "ErrorResult",
		ErrInfo: errResultDetail{
			Reason:    reason,
			ErrorCode: errorCode},
	}
	return json.Marshal(resp)
}
