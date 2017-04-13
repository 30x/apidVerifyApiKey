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
	IssuedAt        int64  `json:"issuedAt"`
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

	// DANGER: This relies on an external TABLE - DATA_SCOPE is maintained by apidApigeeSync
	var env, tenantId string
	error := db.QueryRow("SELECT env, scope FROM DATA_SCOPE WHERE id = ?;", scopeuuid).Scan(&env, &tenantId)

	switch {
	case error == sql.ErrNoRows:
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
			APP_CREDENTIAL AS c 
			INNER JOIN APP AS a ON c.app_id = a.id
			INNER JOIN DEVELOPER AS ad 
				ON ad.id = a.developer_id
			INNER JOIN APP_CREDENTIAL_APIPRODUCT_MAPPER as mp 
				ON mp.appcred_id = c.id 
			INNER JOIN API_PRODUCT as ap ON ap.id = mp.apiprdt_id
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
			APP_CREDENTIAL AS c
			INNER JOIN APP AS a ON c.app_id = a.id
			INNER JOIN COMPANY AS ad
				ON ad.id = a.company_id
			INNER JOIN APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
				ON mp.appcred_id = c.id
			INNER JOIN API_PRODUCT as ap ON ap.id = mp.apiprdt_id
		WHERE (UPPER(ad.status) = 'ACTIVE'
			AND mp.apiprdt_id = ap.id
			AND mp.app_id = a.id
			AND mp.appcred_id = c.id
			AND UPPER(mp.status) = 'APPROVED'
			AND UPPER(a.status) = 'APPROVED'
			AND c.id = $1
			AND c.tenant_id = $2)
	;`

	var status, redirectionURIs, appName, appId, resName, resEnv, cType string
	var issuedAt int64
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
	 * Perform all validations related to the Query made with the dataService
	 * we just retrieved
	 */
	result := validatePath(resName, path)
	if result == false {
		reason := "Path Validation Failed (" + resName + " vs " + path + ")"
		errorCode := "PATH_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)

	}

	/* Verify if the ENV matches */
	result = validateEnv(resEnv, env)
	if result == false {
		reason := "ENV Validation Failed (" + resEnv + " vs " + env + ")"
		errorCode := "ENV_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)
	}

	var expiresAt int64 = -1
	resp := kmsResponseSuccess{
		Type: "APIKeyContext",
		RspInfo: sucResponseDetail{
			Key:             key,
			ExpiresAt:       expiresAt,
			IssuedAt:        issuedAt,
			Status:          status,
			RedirectionURIs: redirectionURIs,
			Type:            cType,
			AppId:           appId,
			AppName:         appName},
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
