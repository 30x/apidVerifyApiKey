package apidVerifyApiKey

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

type sucResponseDetail struct {
	Key             string `json:"key"`
	ExpiresAt       int64  `json:"expiresAt"`
	IssuedAt        int64  `json:"issuedAt"`
	Status          string `json:"status"`
	RedirectionURIs string `json:"redirectionURIs"`
	DeveloperAppId  string `json:"developerId"`
	DeveloperAppNam string `json:"developerAppName"`
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
	if r.Method != "POST" {
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to parse form"))
	}

	f := r.Form
	elems := []string{"action", "key", "uriPath", "organization", "environment"}
	for _, elem := range elems {
		if f.Get(elem) == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Missing element: %s", elem)))
		}
	}

	org := f.Get("organization")
	key := f.Get("key")
	path := f.Get("uriPath")
	env := f.Get("environment")
	action := f.Get("action")

	b, err := verifyAPIKey(key, path, env, org, action)
	if err != nil {
		log.Errorf("error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Debugf("handleVerifyAPIKey result %s", b)
	w.Write(b)
}

// todo: The following was basically just copied from old APID - needs review.

// returns []byte to be written to client
func verifyAPIKey(key, path, env, org, action string) ([]byte, error) {
	var (
		sSql                               string
		status, redirectionURIs            string
		developerAppName, developerId      string
		resName, resEnv, reason, errorCode string
		issuedAt, expiresAt                int64
	)

	if key == "" || org == "" || path == "" || env == "" || action != "verify" {
		log.Error("Input params Invalid/Incomplete")
		reason = "Input Params Incomplete or Invalid"
		errorCode = "INCORRECT_USER_INPUT"
		return errorResponse(reason, errorCode)
	}

	db, err := data.DB()
	if err != nil {
		log.Errorf("Unable to access DB")
		reason = err.Error()
		errorCode = "SEARCH_INTERNAL_ERROR"
		return errorResponse(reason, errorCode)
	}

	sSql = "SELECT ap.api_resources, ap.environments, c.issued_at, c.app_status, a.callback_url, d.username, d.id FROM APP_CREDENTIAL AS c INNER JOIN APP AS a ON c.app_id = a.id INNER JOIN DEVELOPER AS d ON a.developer_id = d.id INNER JOIN APP_CREDENTIAL_APIPRODUCT_MAPPER as mp ON mp.appcred_id = c.id INNER JOIN API_PRODUCT as ap ON ap.id = mp.apiprdt_id WHERE (UPPER(d.status) = 'ACTIVE' AND mp.apiprdt_id = ap.id AND mp.app_id = a.id AND mp.appcred_id = c.id AND UPPER(mp.status) = 'APPROVED' AND UPPER(a.status) = 'APPROVED' AND UPPER(c.status) = 'APPROVED' AND c.id = '" + key + "' AND c._apid_scope = '" + org + "');"

	err = db.QueryRow(sSql).Scan(&resName, &resEnv, &issuedAt, &status,
		&redirectionURIs, &developerAppName, &developerId)
	expiresAt = -1
	switch {
	case err == sql.ErrNoRows:
		reason = "API Key verify failed for (" + key + ", " + org + ", " + path + ", " + env + ")"
		errorCode = "REQ_ENTRY_NOT_FOUND"
		return errorResponse(reason, errorCode)

	case err != nil:
		reason = err.Error()
		errorCode = "SEARCH_INTERNAL_ERROR"
		return errorResponse(reason, errorCode)
	}

	/*
	 * Perform all validations related to the Query made with the data
	 * we just retrieved
	 */
	result := validatePath(resName, path)
	if result == false {
		reason = "Path Validation Failed (" + resName + " vs " + path + ")"
		errorCode = "PATH_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)

	}

	/* Verify if the ENV matches */
	result = validateEnv(resEnv, env)
	if result == false {
		reason = "ENV Validation Failed (" + resEnv + " vs " + env + ")"
		errorCode = "ENV_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)
	}

	resp := kmsResponseSuccess{
		Type: "APIKeyContext",
		RspInfo: sucResponseDetail{
			Key:             key,
			ExpiresAt:       expiresAt,
			IssuedAt:        issuedAt,
			Status:          status,
			RedirectionURIs: redirectionURIs,
			DeveloperAppId:  developerId,
			DeveloperAppNam: developerAppName},
	}
	return json.Marshal(resp)
}

func errorResponse(reason, errorCode string) ([]byte, error) {

	log.Error(reason)
	resp := kmsResponseFail{
		Type: "ErrorResult",
		ErrInfo: errResultDetail{
			Reason:    reason,
			ErrorCode: errorCode},
	}
	return json.Marshal(resp)
}
