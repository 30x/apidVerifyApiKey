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
	"github.com/30x/apid-core"
	"net/http"
	"io/ioutil"
)

// handle client API
func handleRequestv2(w http.ResponseWriter, r *http.Request) {

	db := getDB()
	if db == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("initializing"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println(string(body))
	var verifyApiKeyReq VerifyApiKeyRequest
	err = json.Unmarshal(body, &verifyApiKeyReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Info(verifyApiKeyReq)

	if verifyApiKeyReq.Action == "" || verifyApiKeyReq.ApiProxyName == "" || verifyApiKeyReq.EnvironmentName == "" || verifyApiKeyReq.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Missing element: %s", verifyApiKeyReq)))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := verifyAPIKeyv2(verifyApiKeyReq)
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
func verifyAPIKeyv2(verifyApiKeyReq VerifyApiKeyRequest) ([]byte, error) {

	key := verifyApiKeyReq.Key
	organizationName := verifyApiKeyReq.OrganizationName
	environmentName := verifyApiKeyReq.EnvironmentName
	path := verifyApiKeyReq.UriPath
	action := verifyApiKeyReq.Action

	if key == "" || organizationName == "" || environmentName == "" || path == "" || action != "verify" {
		log.Debug("Input params Invalid/Incomplete")
		reason := "Input Params Incomplete or Invalid"
		errorCode := "INCORRECT_USER_INPUT"
		return errorResponse(reason, errorCode)
	}

	db := getDB()

	sSql := `
		SELECT
			ap.id,
			ap.api_resources,
			ap.environments,
			c.issued_at,
			c.status,
			a.callback_url,
			a.id,
			ad.email,
			ad.id,
			"developer" as ctype,
			c.consumer_secret,
			a.callback_url,
			c.tenant_id
		FROM
			KMS_APP_CREDENTIAL AS c
			INNER JOIN KMS_APP AS a
				ON c.app_id = a.id
			INNER JOIN KMS_DEVELOPER AS ad
				ON ad.id = a.developer_id
			INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
				ON mp.appcred_id = c.id
			INNER JOIN KMS_API_PRODUCT as ap
				ON ap.id = mp.apiprdt_id
			INNER JOIN KMS_ORGANIZATION AS o
				ON o.tenant_id = c.tenant_id
		WHERE (UPPER(ad.status) = 'ACTIVE'
			AND mp.apiprdt_id = ap.id
			AND mp.app_id = a.id
			AND mp.appcred_id = c.id
			AND UPPER(mp.status) = 'APPROVED'
			AND UPPER(a.status) = 'APPROVED'
			AND c.id = $1
			AND o.name = $2)
		UNION ALL
		SELECT
			ap.id,
			ap.api_resources,
			ap.environments,
			c.issued_at,
			c.status,
			a.callback_url,
			a.id,
			ad.name,
			ad.id,
			"company" as ctype,
			c.consumer_secret,
			a.callback_url,
			c.tenant_id
		FROM
			KMS_APP_CREDENTIAL AS c
			INNER JOIN KMS_APP AS a
				ON c.app_id = a.id
			INNER JOIN KMS_COMPANY AS ad
				ON ad.id = a.company_id
			INNER JOIN KMS_APP_CREDENTIAL_APIPRODUCT_MAPPER as mp
				ON mp.appcred_id = c.id
			INNER JOIN KMS_API_PRODUCT as ap
				ON ap.id = mp.apiprdt_id
			INNER JOIN KMS_ORGANIZATION AS o
				ON o.tenant_id = c.tenant_id
		WHERE (UPPER(ad.status) = 'ACTIVE'
			AND mp.apiprdt_id = ap.id
			AND mp.app_id = a.id
			AND mp.appcred_id = c.id
			AND UPPER(mp.status) = 'APPROVED'
			AND UPPER(a.status) = 'APPROVED'
			AND c.id = $1
			AND o.name = $2)
	;`

	/* these fields need to be nullable types for scanning.  This is because when using json snapshots,
	   and therefore being responsible for inserts, we were able to default everything to be not null.  With
	   sqlite snapshots, we are not necessarily guaranteed that
	*/
	var status, redirectionURIs, appName, devId, resName, resEnv, issuedAt, cType, cSecret, callback, aId, apiProductId, tenantId sql.NullString
	err := db.QueryRow(sSql, key, organizationName).Scan(&apiProductId, &resName, &resEnv, &issuedAt, &status,
		&redirectionURIs, &aId, &appName, &devId, &cType, &cSecret, &callback, &tenantId)
	switch {
	case err == sql.ErrNoRows:
		reason := "API Key verify failed for (" + key + ", " + organizationName + ", " + path + ")"
		errorCode := "REQ_ENTRY_NOT_FOUND"
		return errorResponsev2(reason, errorCode)

	case err != nil:
		reason := err.Error()
		errorCode := "SEARCH_INTERNAL_ERROR"
		return errorResponsev2(reason, errorCode)
	}

	/*
	 * Perform all validations related to the Query made with the data
	 * we just retrieved
	 */
	result := validatePath(resName.String, path)
	if result == false {
		reason := "Path Validation Failed (" + resName.String + " vs " + path + ")"
		errorCode := "PATH_VALIDATION_FAILED"
		return errorResponsev2(reason, errorCode)

	}

	/* Verify if the ENV matches */
	result = validateEnv(resEnv.String, environmentName)
	if result == false {
		reason := "ENV Validation Failed (" + resEnv.String + " vs " + environmentName + ")"
		errorCode := "ENV_VALIDATION_FAILED"
		return errorResponsev2(reason, errorCode)
	}

	//var expiresAt int64 = -1
	// select * from kms_attributes where tenant_id = '' and entity_id = 'key'
	clientIdAttributes := getKmsAttributes(db, tenantId.String, key)

	clientIdDetails := ClientIdDetails{
		ClientId:     key,
		ClientSecret: cSecret.String,
		Status:       status.String,
		// Attributes associated with the client Id.
		Attributes: clientIdAttributes,
	}
	if callback.String != "" {
		clientIdDetails.RedirectURIs = []string{callback.String}
	}
	var developerDetails DeveloperDetails
	var companyDetails CompanyDetails
	if cType.String == "developer" {
		developerDetails = DeveloperDetails{}
		developerAttributes := getKmsAttributes(db, tenantId.String, devId.String)
		getDevQuery := "select id, username, first_name, last_name, email, status, created_at, created_by, updated_at, updated_by from kms_developer where tenant_id = $1 and id = $2"
		developer := db.QueryRow(getDevQuery, tenantId, devId.String)
		err = developer.Scan(
			&developerDetails.Id,
			&developerDetails.UserName,
			&developerDetails.FirstName,
			&developerDetails.LastName,
			&developerDetails.Email,
			&developerDetails.Status,
			&developerDetails.CreatedAt,
			&developerDetails.CreatedBy,
			&developerDetails.LastmodifiedAt,
			&developerDetails.LastmodifiedBy,
			// TODO : check do we need all ??
			//Apps []string `json:"apps,omitempty"
			// TODO : is this being used ??
			//Company string `json:"company,omitempty"`
		)
		developerDetails.Attributes = developerAttributes
		if err != nil {
			log.Error("error getting developerDetails" , err)
		}
	} else {
		companyDetails := CompanyDetails{}
		companyAttributes := getKmsAttributes(db, tenantId.String, devId.String)
		getCompanyQuery := "select id, display_name, status, created_at, created_by, updated_at. updated_by from kms_company where tenant_id = $1 and id = $2"
		company := db.QueryRow(getCompanyQuery)
		err = company.Scan(
			&companyDetails.Id,
			&companyDetails.DisplayName,
			&companyDetails.Status,
			&companyDetails.CreatedAt,
			&companyDetails.CreatedBy,
			&companyDetails.LastmodifiedAt,
			&companyDetails.LastmodifiedBy,
			// TODO : check do we need all ??
			//Apps []string `json:"apps,omitempty"`
		)
		companyDetails.Attributes = companyAttributes
		if err != nil {
			log.Error("error getting companyDetails " , err)
		}
	}

	appAttributes := getKmsAttributes(db, tenantId.String, aId.String)
	appDetails := AppDetails{}
	getAppDetails := "select id, name, access_type, callback_url, display_name, status, app_family, company_id, created_at, created_by, updated_at, updated_by from kms_app where tenant_id = $1 and id = $2 "
	apps := db.QueryRow(getAppDetails, tenantId, aId.String)

	err = apps.Scan(
		&appDetails.Id,
		&appDetails.Name,
		&appDetails.AccessType,
		&appDetails.CallbackUrl,
		&appDetails.DisplayName,
		&appDetails.Status,
		&appDetails.AppFamily,
		&appDetails.Company,
		&appDetails.CreatedAt,
		&appDetails.CreatedBy,
		&appDetails.LastmodifiedAt,
		&appDetails.LastmodifiedBy,
	)
	if err != nil {
		log.Error("error getting apps ",err)
	}
	appDetails.Attributes = appAttributes

	apiProductAttributes := getKmsAttributes(db, tenantId.String, apiProductId.String)
	apiProductDetails := ApiProductDetails{}
	getApiProductsQuery := "select id, name, display_name, quota, COALESCE(quota_interval, '') , quota_time_unit, created_at, created_by, updated_at, updated_by, proxies, environments from kms_api_product where tenant_id = $1 and id = $2 "
	apiProducts := db.QueryRow(getApiProductsQuery, tenantId, apiProductId.String)
	var proxies, environments string
	err = apiProducts.Scan(
		&apiProductDetails.Id,
		&apiProductDetails.Name,
		&apiProductDetails.DisplayName,
		&apiProductDetails.QuotaLimit,
		&apiProductDetails.QuotaInterval,
		&apiProductDetails.QuotaTimeunit,
		&apiProductDetails.CreatedAt,
		&apiProductDetails.CreatedBy,
		&apiProductDetails.LastmodifiedAt,
		&apiProductDetails.LastmodifiedBy,
		&proxies,
		&environments,
	)
	if err != nil {
		log.Error("error getting apiProductDetails " , err)
	}

	if err := json.Unmarshal([]byte(proxies), &apiProductDetails.Apiproxies); err != nil {
		log.Debug("unmarshall error for proxies, sending as is " , err)
		apiProductDetails.Apiproxies = []string{ proxies }
	}
	if err := json.Unmarshal([]byte(environments), &apiProductDetails.Environments); err != nil {
		log.Debug("unmarshall error for proxies, sending as is " , err)
		apiProductDetails.Environments = []string{ environments }
	}

	apiProductDetails.Attributes = apiProductAttributes

	resp := VerifyApiKeySuccessResponse{
		ClientId: clientIdDetails,
		Organization: organizationName,
		Environment: resEnv.String,
		Developer:   developerDetails,
		Company:     companyDetails,
		App:         appDetails,
		ApiProduct:  apiProductDetails,
		// Identifier of the authorization code. This will be unique for each request.
		Identifier: key,        // TODO : what is this ?????
		Kind:       "Collection", // TODO : what is this ????

	}

	return json.Marshal(resp)
}
func getKmsAttributes(db apid.DB, tenantId string, entityId string) []Attribute {
	sql := "select name, value, type from kms_attributes where tenant_id = $1 and entity_id = $2"
	attributesForQuery := []Attribute{}
	attributes, err := db.Query(sql, tenantId, entityId)

	if err != nil {
		log.Error("Error while fetching attributes for tenant id : %s and entityId : %s", tenantId, entityId, err)
		return attributesForQuery
	}
	for attributes.Next() {
		att := Attribute{}
		attributes.Scan(
			&att.Name,
			&att.Value,
		)
		if att.Name != "" {
			attributesForQuery = append(attributesForQuery, att)
		}
	}
	return attributesForQuery
}

func errorResponsev2(reason, errorCode string) ([]byte, error) {
	if errorCode == "SEARCH_INTERNAL_ERROR" {
		log.Error(reason)
	} else {
		log.Debug(reason)
	}
	resp := ErrorResponse{
		ResponseCode:    errorCode,
		ResponseMessage: reason,
	}
	return json.Marshal(resp)
}
