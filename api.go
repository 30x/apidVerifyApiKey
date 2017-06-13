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
	elems := []string{"action", "key", "uriPath", "scopeuuid"}
	for _, elem := range elems {
		if f.Get(elem) == "" {
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

	key := f.Get("key")
	scopeuuid := f.Get("scopeuuid")
	path := f.Get("uriPath")
	action := f.Get("action")

	if key == "" || scopeuuid == "" || path == "" || action != "verify" {
		log.Debug("Input params Invalid/Incomplete")
		reason := "Input Params Incomplete or Invalid"
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
