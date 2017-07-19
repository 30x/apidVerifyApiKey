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
	"github.com/30x/apid-core"
)

// handle client API
func handleRequestv2(w http.ResponseWriter, r *http.Request) {

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
	elems := []string{"action", "key", "uriPath", "organizationName", "environmentName", "apiProxyName", "validateAgainstApiProxiesAndEnvs"}
	for _, elem := range elems {
		if f.Get(elem) == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Missing element: %s", elem)))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	b, err := verifyAPIKeyv2(f)
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
func verifyAPIKeyv2(f url.Values) ([]byte, error) {

	key := f.Get("key")
	organizationName := f.Get("organizationName")
	environmentName := f.Get("environmentName")
	path := f.Get("uriPath")
	action := f.Get("action")

	if key == "" || organizationName == "" || environmentName == "" || path == "" || action != "verify" {
		log.Debug("Input params Invalid/Incomplete")
		reason := "Input Params Incomplete or Invalid"
		errorCode := "INCORRECT_USER_INPUT"
		return errorResponse(reason, errorCode)
	}

	db := getDB()
	tenantId := f.Get("tenantId")

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
			a.callback_url
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
			a.callback_url
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
	var status, redirectionURIs, appName, devId, resName, resEnv, issuedAt, cType, cSecret, callback, aId, apiProductId sql.NullString
	err := db.QueryRow(sSql, key, tenantId).Scan(&apiProductId, &resName, &resEnv, &issuedAt, &status,
		&redirectionURIs, &aId, &appName, &devId, &cType, &cSecret, &callback)
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
	clientIdAttributes := getKmsAttributes(db, tenantId, key)

	clientIdDetails := ClientIdDetails{
		ClientId: key,

		ClientSecret: cSecret.String,

		//RedirectURIs : { callback.String },

		Status: status.String,

		// Attributes associated with the client Id.
		Attributes: clientIdAttributes,
	}
	var developerDetails DeveloperDetails
	var companyDetails CompanyDetails
	if cType.String == "developer" {
		// TODO : get developer details
		// select * from kms_developer where tenant_id = '' and id = ' devId.String'
		developerAttributes := getKmsAttributes(db, tenantId, devId.String)
		developerDetails = DeveloperDetails{
			Id: devId.String,
			//UserName : "",
			//FirstName "",
			//LastName string `json:"lastName,omitempty"`
			//Email string `json:"email,omitempty"`
			//Status string `json:"status,omitempty"`
			//
			// TODO : check do we need all ??
			//Apps []string `json:"apps,omitempty"`
			//
			//CreatedAt int64 `json:"created_at,omitempty"`
			//
			//CreatedBy string `json:"created_by,omitempty"`
			//
			//LastmodifiedAt int64 `json:"lastmodified_at,omitempty"`
			//
			//LastmodifiedBy string `json:"lastmodified_by,omitempty"`
			//
			//Company string `json:"company,omitempty"`
			//
			//// Attributes associated with the developer.
			// TODO : fetch this
			// select * from kms_attributes where tenant_id = '' and entity_id = 'key'
			Attributes : developerAttributes,
		}
	} else {
		// TODO :get company details
		// select * from kms_company where tenant_id = '' and id = ' devId.String'
		companyAttributes := getKmsAttributes(db, tenantId, devId.String)
		companyDetails = CompanyDetails{
			Id: devId.String,
			//
			//Name string `json:"name,omitempty"`
			//
			//DisplayName string `json:"displayName,omitempty"`
			//
			//Status string `json:"status,omitempty"`
			//
			// TODO : check do we need all ??
			//Apps []string `json:"apps,omitempty"`
			//
			//CreatedAt int64 `json:"created_at,omitempty"`
			//
			//CreatedBy string `json:"created_by,omitempty"`
			//
			//LastmodifiedAt int64 `json:"lastmodified_at,omitempty"`
			//
			//LastmodifiedBy string `json:"lastmodified_by,omitempty"`
			//
			//// Attributes associated with the company.
			// TODO : fetch this
			// select * from kms_attributes where tenant_id = '' and entity_id = 'key'
			Attributes: companyAttributes,
		}
	}

	// TODO :get app details
	// select * from kms_app where tenant_id = '' and id = 'aId.String'
	appAttributes := getKmsAttributes(db, tenantId, aId.String)
	appDetails := AppDetails{
		Id: aId.String,
		//
		//Name string `json:"name,omitempty"`
		//
		//AccessType string `json:"accessType,omitempty"`
		//
		//CallbackUrl string `json:"callbackUrl,omitempty"`
		//
		//DisplayName string `json:"displayName,omitempty"`
		//
		//Status string `json:"status,omitempty"`
		//
		// TODO : apiproducts - is this specific to credential or app ?? - do we have to get all products for an app or credential
		//Apiproducts []string `json:"apiproducts,omitempty"`
		//
		//AppFamily string `json:"appFamily,omitempty"`
		//
		//CreatedAt int64 `json:"created_at,omitempty"`
		//
		//CreatedBy string `json:"created_by,omitempty"`
		//
		//LastmodifiedAt int64 `json:"lastmodified_at,omitempty"`
		//
		//LastmodifiedBy string `json:"lastmodified_by,omitempty"`
		//
		//Company string `json:"company,omitempty"`
		//
		//// Attributes associated with the app.
		// TODO : fetch this
		// select * from kms_attributes where tenant_id = '' and entity_id = 'key'
		Attributes: appAttributes,
	}

	// TODO : fetch and populate
	// select * from kms_developer where tenant_id = '' and id = ' devId.String'
	apiProductAttributes := getKmsAttributes(db, tenantId, apiProductId.String)
	apiProductDetails := ApiProductDetails{

		Id: apiProductId.String,

		//Name string `json:"name,omitempty"`

		//DisplayName string `json:"displayName,omitempty"`
		//
		//QuotaLimit int64 `json:"quota.limit,omitempty"`
		//
		//QuotaInterval int64 `json:"quota.interval,omitempty"`
		//
		//QuotaTimeunit int64 `json:"quota.timeunit,omitempty"`
		//
		//Status string `json:"status,omitempty"`
		//
		//CreatedAt int64 `json:"created_at,omitempty"`
		//
		//CreatedBy string `json:"created_by,omitempty"`
		//
		//LastmodifiedAt int64 `json:"lastmodified_at,omitempty"`
		//
		//LastmodifiedBy string `json:"lastmodified_by,omitempty"`
		//
		//Company string `json:"company,omitempty"`
		//
		//Environments []string `json:"environments,omitempty"`
		//
		//Apiproxies []string `json:"apiproxies,omitempty"`
		//
		//TODO : fetch apiProduct attributes
		// Attributes associated with the apiproduct.
		// TODO : fetch this
		// select * from kms_attributes where tenant_id = '' and entity_id = 'key'
		Attributes: apiProductAttributes,
	}

	resp := VerifyApiKeySuccessResponse{
		ClientId: clientIdDetails,
		// Self string `json:"self,omitempty"`

		// Organization Identifier/Name
		Organization: "test", // TODO : find where to get this info from and fix.

		// Environment Identifier/Name
		Environment: resEnv.String,

		Developer: developerDetails,

		Company: companyDetails,

		App: appDetails,

		ApiProduct: apiProductDetails,

		// Identifier of the authorization code. This will be unique for each request.
		Identifier: "id", // TODO : what is this ?????

		Kind: "your_kind", // TODO : what is this ????

	}

	return json.Marshal(resp)
}
func getKmsAttributes(db apid.DB, tenantId string, entityId string) []Attribute {
	sql := "select name, value, type from kms_attributes where tenant_id = $1 and entity_id = $2"
	attributesForQuery := []Attribute{}
	attributes, err := db.Query(sql,tenantId, entityId)

	if(err == nil){
		log.Error("Error while fetching attributes for tenant id : [{}] and entityId : [{}]", tenantId, entityId ,err)
		return attributesForQuery;
	}
	for attributes.Next() {
		att := Attribute{}
		attributes.Scan(
			&att.Name,
			&att.Value,
			&att.Kind,
		)
		attributesForQuery = append(attributesForQuery, att)
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
