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
	"github.com/30x/apid-core"
	"net/http"
	"io/ioutil"
	"errors"
	"io"
)

type apiManagerInterface interface {
	InitAPI()
	//addChangedDeployment(string)
	//distributeEvents()
}
type apiManager struct {
	dbMan               dbManagerInterface
	verifiersEndpoint string
	apiInitialized      bool
}

func (a *apiManager) InitAPI() {
	if a.apiInitialized {
		return
	}
	services.API().HandleFunc(a.verifiersEndpoint, a.handleRequest).Methods("POST")
	a.apiInitialized = true
	log.Debug("API endpoints initialized")
}

// handle client API
func (a *apiManager) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	verifyApiKeyReq, err := validateRequest(r.Body, w)
	if err != nil {
		return
	}

	b, err := verifyAPIKey(verifyApiKeyReq, a.dbMan.getDb())

	if err != nil {
		log.Errorf("error: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Debugf("handleVerifyAPIKey result %s", b)
	w.Write(b)
}

func validateRequest(requestBody io.ReadCloser, w http.ResponseWriter) (VerifyApiKeyRequest, error) {
	// 1. read request boby
	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return VerifyApiKeyRequest{}, errors.New("Bad_REQUEST")
	}
	log.Debug(string(body))
	// 2. umarshall json to struct
	var verifyApiKeyReq VerifyApiKeyRequest
	err = json.Unmarshal(body, &verifyApiKeyReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return VerifyApiKeyRequest{}, errors.New("Bad_REQUEST")
	}
	log.Debug(verifyApiKeyReq)

	// 2. verify params
	if verifyApiKeyReq.Action == "" || verifyApiKeyReq.ApiProxyName == "" || verifyApiKeyReq.EnvironmentName == "" || verifyApiKeyReq.Key == "" {
		// TODO : set correct fields in error response
		errorResponse , _ := errorResponse("Bad_REQUEST","Missing element")
		w.Write(errorResponse)
		return VerifyApiKeyRequest{}, errors.New("Bad_REQUEST")
	}
	return verifyApiKeyReq, nil
}

// returns []byte to be written to client
func verifyAPIKey(verifyApiKeyReq VerifyApiKeyRequest, db apid.DB) ([]byte, error) {

	key := verifyApiKeyReq.Key
	organizationName := verifyApiKeyReq.OrganizationName
	environmentName := verifyApiKeyReq.EnvironmentName
	path := verifyApiKeyReq.UriPath
	//action := verifyApiKeyReq.Action

	sSql := `
		SELECT
			ap.api_resources,
			ap.environments,
			"developer" as ctype,
			c.tenant_id,

			c.status,
			c.consumer_secret,

			ad.id as dev_id,
			ad.username as dev_username,
			ad.first_name as dev_first_name,
			ad.last_name as dev_last_name,
			ad.email as dev_email,
			ad.status as dev_status,
			ad.created_at as dev_created_at,
			ad.created_by as dev_created_by,
			ad.updated_at as dev_updated_at,
			ad.updated_by as dev_updated_by,

			a.id as app_id,
			a.name as app_name,
			a.access_type as app_access_type,
			a.callback_url as app_callback_url,
			a.display_name as app_display_name,
			a.status as app_status,
			a.app_family as app_app_family,
			a.company_id as app_company_id,
			a.created_at as app_created_at,
			a.created_by as app_created_by,
			a.updated_at as app_updated_at,
			a.updated_by as app_updated_by,

			ap.id as prod_id,
			ap.name as prod_name,
			ap.display_name as prod_display_name,
			ap.quota as prod_quota,
			COALESCE(ap.quota_interval, '') as prod_quota_interval,
			ap.quota_time_unit as prod_quota_time_unit,
			ap.created_at as prod_created_at,
			ap.created_by as prod_created_by,
			ap.updated_at as prod_updated_at,
			ap.updated_by as prod_updated_by,
			ap.proxies as prod_proxies,
			ap.environments as prod_environments
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
			ap.api_resources,
			ap.environments,
			"company" as ctype,
			c.tenant_id,

			c.status,
			c.consumer_secret,

			ad.id as dev_id,
			ad.display_name as dev_username,
			"" as dev_first_name,
			"" as dev_last_name,
			"" as dev_email,
			ad.status as dev_status,
			ad.created_at as dev_created_at,
			ad.created_by as dev_created_by,
			ad.updated_at as dev_updated_at,
			ad.updated_by as dev_updated_by,

			a.id as app_id,
			a.name as app_name,
			a.access_type as app_access_type,
			a.callback_url as app_callback_url,
			a.display_name as app_display_name,
			a.status as app_status,
			a.app_family as app_app_family,
			a.company_id as app_company_id,
			a.created_at as app_created_at,
			a.created_by as app_created_by,
			a.updated_at as app_updated_at,
			a.updated_by as app_updated_by,

			ap.id as prod_id,
			ap.name as prod_name,
			ap.display_name as prod_display_name,
			ap.quota as prod_quota,
			COALESCE(ap.quota_interval, '') as prod_quota_interval,
			ap.quota_time_unit as prod_quota_time_unit,
			ap.created_at as prod_created_at,
			ap.created_by as prod_created_by,
			ap.updated_at as prod_updated_at,
			ap.updated_by as prod_updated_by,
			ap.proxies as prod_proxies,
			ap.environments as prod_environments

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
	var finalDeveloperDetails DeveloperDetails
	var proxies, environments string
	var resName, resEnv, cType, tenantId sql.NullString

	tempDeveloperDetails := DeveloperDetails{}
	companyDetails := CompanyDetails{}
	appDetails := AppDetails{}
	apiProductDetails := ApiProductDetails{}
	clientIdDetails := ClientIdDetails{}
	clientIdDetails.ClientId = key

	err := db.QueryRow(sSql, key, organizationName).
		Scan(
			&resName,
			&resEnv,
			&cType,
			&tenantId,
			&clientIdDetails.Status,
			&clientIdDetails.ClientSecret,

			&tempDeveloperDetails.Id,
			&tempDeveloperDetails.UserName,
			&tempDeveloperDetails.FirstName,
			&tempDeveloperDetails.LastName,
			&tempDeveloperDetails.Email,
			&tempDeveloperDetails.Status,
			&tempDeveloperDetails.CreatedAt,
			&tempDeveloperDetails.CreatedBy,
			&tempDeveloperDetails.LastmodifiedAt,
			&tempDeveloperDetails.LastmodifiedBy,

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
	switch {
		case err == sql.ErrNoRows:
			reason := "API Key verify failed for (" + key + ", " + organizationName + ", " + path + ")"
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
	result = validateEnv(resEnv.String, environmentName)
	if result == false {
		reason := "ENV Validation Failed (" + resEnv.String + " vs " + environmentName + ")"
		errorCode := "ENV_VALIDATION_FAILED"
		return errorResponse(reason, errorCode)
	}

	clientIdAttributes := getKmsAttributes(db, tenantId.String, clientIdDetails.ClientId)
	appAttributes := getKmsAttributes(db, tenantId.String, appDetails.Id)
	apiProductAttributes := getKmsAttributes(db, tenantId.String, apiProductDetails.Id)
	developerAttributes := getKmsAttributes(db, tenantId.String, tempDeveloperDetails.Id)

	clientIdDetails.Attributes = clientIdAttributes
	appDetails.Attributes = appAttributes
	apiProductDetails.Attributes = apiProductAttributes

	if appDetails.CallbackUrl != "" {
		clientIdDetails.RedirectURIs = []string{appDetails.CallbackUrl}
	}
	if err := json.Unmarshal([]byte(proxies), &apiProductDetails.Apiproxies); err != nil {
		log.Debug("unmarshall error for proxies, sending as is ", err)
		apiProductDetails.Apiproxies = []string{proxies }
	}
	if err := json.Unmarshal([]byte(environments), &apiProductDetails.Environments); err != nil {
		log.Debug("unmarshall error for proxies, sending as is ", err)
		apiProductDetails.Environments = []string{environments }
	}

	if cType.String == "developer" {
		finalDeveloperDetails := &tempDeveloperDetails
		finalDeveloperDetails.Attributes = developerAttributes
	} else {
		companyDetails := CompanyDetails{
			Id: tempDeveloperDetails.Id,
			DisplayName: tempDeveloperDetails.UserName,
			Status: tempDeveloperDetails.Status,
			CreatedAt: tempDeveloperDetails.CreatedAt,
			CreatedBy: tempDeveloperDetails.CreatedBy,
			LastmodifiedAt: tempDeveloperDetails.LastmodifiedAt,
			LastmodifiedBy: tempDeveloperDetails.LastmodifiedBy,
		}
		companyDetails.Attributes = developerAttributes
	}

	resp := VerifyApiKeySuccessResponse{
		ClientId: clientIdDetails,
		Organization: organizationName,
		Environment: resEnv.String,
		Developer:   finalDeveloperDetails,
		Company:     companyDetails,
		App:         appDetails,
		ApiProduct:  apiProductDetails,
		// Identifier of the authorization code. This will be unique for each request.
		Identifier: key, // TODO : what is this ?????
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

func errorResponse(reason, errorCode string) ([]byte, error) {
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
