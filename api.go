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
	"errors"
	"github.com/30x/apid-core"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type apiManagerInterface interface {
	InitAPI()
	//addChangedDeployment(string)
	//distributeEvents()
}
type apiManager struct {
	dbMan             dbManagerInterface
	verifiersEndpoint string
	apiInitialized    bool
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
		respStatusCode, atoierr := strconv.Atoi(err.Error())
		if atoierr != nil {
			w.WriteHeader(respStatusCode)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	log.Debugf("handleVerifyAPIKey result %s", b)
	w.Write(b)
	return
}

func validateRequest(requestBody io.ReadCloser, w http.ResponseWriter) (VerifyApiKeyRequest, error) {
	// 1. read request boby
	var verifyApiKeyReq VerifyApiKeyRequest
	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return verifyApiKeyReq, errors.New("Bad_REQUEST")
	}
	log.Debug(string(body))
	// 2. umarshall json to struct
	err = json.Unmarshal(body, &verifyApiKeyReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return verifyApiKeyReq, errors.New("Bad_REQUEST")
	}
	log.Debug(verifyApiKeyReq)

	// 2. verify params
	if verifyApiKeyReq.Action == "" || verifyApiKeyReq.ApiProxyName == "" || verifyApiKeyReq.EnvironmentName == "" || verifyApiKeyReq.Key == "" {
		// TODO : set correct fields in error response
		errorResponse := errorResponse("Bad_REQUEST", "Missing element")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse)
		return verifyApiKeyReq, errors.New("Bad_REQUEST")
	}
	return verifyApiKeyReq, nil
}

// returns []byte to be written to client
func verifyAPIKey(verifyApiKeyReq VerifyApiKeyRequest, db apid.DB) ([]byte, error) {

	/* these fields need to be nullable types for scanning.  This is because when using json snapshots,
	   and therefore being responsible for inserts, we were able to default everything to be not null.  With
	   sqlite snapshots, we are not necessarily guaranteed that
	*/
	var finalDeveloperDetails DeveloperDetails
	var companyDetails CompanyDetails
	var cType, tenantId sql.NullString

	tempDeveloperDetails := DeveloperDetails{}
	appDetails := AppDetails{}
	apiProductDetails := ApiProductDetails{}
	clientIdDetails := ClientIdDetails{}
	clientIdDetails.ClientId = verifyApiKeyReq.Key

	err := getApiKeyDetails(db, verifyApiKeyReq, &cType, &tenantId, &tempDeveloperDetails, &appDetails, &apiProductDetails, &clientIdDetails)

	switch {
	case err == sql.ErrNoRows:
		reason := "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode := "oauth.v2.InvalidApiKey"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))

	case err != nil:
		reason := err.Error()
		errorCode := "SEARCH_INTERNAL_ERROR"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusInternalServerError))
	}

	/*
	 * Perform all validations related to the Query made with the data
	 * we just retrieved
	 */
	errResponse, err := performValidations(verifyApiKeyReq, clientIdDetails, appDetails, tempDeveloperDetails, apiProductDetails, cType)

	if errResponse != nil {
		return errResponse, err
	}

	enrichAttributes(db, tenantId.String, &clientIdDetails, &appDetails, &tempDeveloperDetails, &apiProductDetails)

	if cType.String == "developer" {
		finalDeveloperDetails = tempDeveloperDetails
	} else {
		companyDetails = CompanyDetails{
			Id:             tempDeveloperDetails.Id,
			DisplayName:    tempDeveloperDetails.UserName,
			Status:         tempDeveloperDetails.Status,
			CreatedAt:      tempDeveloperDetails.CreatedAt,
			CreatedBy:      tempDeveloperDetails.CreatedBy,
			LastmodifiedAt: tempDeveloperDetails.LastmodifiedAt,
			LastmodifiedBy: tempDeveloperDetails.LastmodifiedBy,
			Attributes:     tempDeveloperDetails.Attributes,
		}

	}

	resp := VerifyApiKeySuccessResponse{
		ClientId:     clientIdDetails,
		Organization: verifyApiKeyReq.OrganizationName,
		Environment:  verifyApiKeyReq.EnvironmentName,
		Developer:    finalDeveloperDetails,
		Company:      companyDetails,
		App:          appDetails,
		ApiProduct:   apiProductDetails,
		// Identifier of the authorization code. This will be unique for each request.
		Identifier: verifyApiKeyReq.Key, // TODO : what is this ?????
		Kind:       "Collection",        // TODO : what is this ????

	}

	return json.Marshal(resp)
}

func performValidations(verifyApiKeyReq VerifyApiKeyRequest, clientIdDetails ClientIdDetails, appDetails AppDetails, tempDeveloperDetails DeveloperDetails, apiProductDetails ApiProductDetails, cType sql.NullString) ([]byte, error) {
	if !strings.EqualFold("APPROVED", clientIdDetails.Status) {
		reason := "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode := "oauth.v2.ApiKeyNotApproved"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	if !strings.EqualFold("APPROVED", appDetails.Status) {
		reason := "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode := "keymanagement.service.invalid_client-app_not_approved"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	if !strings.EqualFold("ACTIVE", tempDeveloperDetails.Status) {
		reason := "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode := "keymanagement.service.DeveloperStatusNotActive"
		if cType.String == "company" {
			errorCode = "keymanagement.service.CompanyStatusNotActive"
		}
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	result := validatePathRegex(apiProductDetails.Resources, verifyApiKeyReq.UriPath)
	if result == false {
		reason := "Path Validation Failed (" + strings.Join(apiProductDetails.Resources, ", ") + " vs " + verifyApiKeyReq.UriPath + ")"
		errorCode := "oauth.v2.InvalidApiKeyForGivenResource"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	/* Verify if the ENV matches */
	if verifyApiKeyReq.ValidateAgainstApiProxiesAndEnvs && !contains(apiProductDetails.Environments, verifyApiKeyReq.EnvironmentName) {
		reason := "ENV Validation Failed (" + strings.Join(apiProductDetails.Environments, ", ") + " vs " + verifyApiKeyReq.EnvironmentName + ")"
		errorCode := "oauth.v2.InvalidApiKeyForGivenResource"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	if verifyApiKeyReq.ValidateAgainstApiProxiesAndEnvs && !contains(apiProductDetails.Apiproxies, verifyApiKeyReq.ApiProxyName) {
		reason := "Proxy Validation Failed (" + strings.Join(apiProductDetails.Apiproxies, ", ") + " vs " + verifyApiKeyReq.ApiProxyName + ")"
		errorCode := "oauth.v2.InvalidApiKeyForGivenResource"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	return nil, nil

}

func contains(givenArray []string, searchString string) bool {
	for _, element := range givenArray {
		if element == searchString {
			return true
		}
	}
	return false
}

func enrichAttributes(db apid.DB, tenantId string, clientIdDetails *ClientIdDetails, appDetails *AppDetails, tempDeveloperDetails *DeveloperDetails, apiProductDetails *ApiProductDetails) {
	clientIdAttributes := getKmsAttributes(db, tenantId, clientIdDetails.ClientId)
	developerAttributes := getKmsAttributes(db, tenantId, tempDeveloperDetails.Id)
	appAttributes := getKmsAttributes(db, tenantId, appDetails.Id)
	apiProductAttributes := getKmsAttributes(db, tenantId, apiProductDetails.Id)

	clientIdDetails.Attributes = clientIdAttributes
	appDetails.Attributes = appAttributes
	apiProductDetails.Attributes = apiProductAttributes
	tempDeveloperDetails.Attributes = developerAttributes
}
func getKmsAttributes(db apid.DB, tenantId string, entityId string) []Attribute {

	var attName, attValue sql.NullString
	sql := "select name, value from kms_attributes where tenant_id = $1 and entity_id = $2"
	attributesForQuery := []Attribute{}
	attributes, err := db.Query(sql, tenantId, entityId)
	if err != nil {
		log.Error("Error while fetching attributes for tenant id : %s and entityId : %s", tenantId, entityId, err)
		return attributesForQuery
	}

	for attributes.Next() {
		err := attributes.Scan(
			&attName,
			&attValue,
		)
		if err != nil {
			log.Error("error fetching attributes for entityid ", entityId, err)
		}
		if attName.String != "" {
			att := Attribute{Name: attName.String, Value: attValue.String}
			attributesForQuery = append(attributesForQuery, att)
		}
	}
	log.Debug("attributes returned for query ", sql, " are ", attributesForQuery, tenantId, entityId)
	return attributesForQuery
}

func errorResponse(reason, errorCode string) []byte {
	if errorCode == "SEARCH_INTERNAL_ERROR" {
		log.Error(reason)
	} else {
		log.Debug(reason)
	}
	resp := ErrorResponse{
		ResponseCode:    errorCode,
		ResponseMessage: reason,
	}
	ret, _ := json.Marshal(resp)
	return ret
}

func getApiKeyDetails(db apid.DB, verifyApiKeyReq VerifyApiKeyRequest, cType, tenantId *sql.NullString, tempDeveloperDetails *DeveloperDetails, appDetails *AppDetails, apiProductDetails *ApiProductDetails, clientIdDetails *ClientIdDetails) error {

	var proxies, environments, resources string
	sSql := `
		SELECT
			COALESCE("developer","") as ctype,
			COALESCE(c.tenant_id,""),

			COALESCE(c.status,""),
			COALESCE(c.consumer_secret,""),

			COALESCE(ad.id,"") as dev_id,
			COALESCE(ad.username,"") as dev_username,
			COALESCE(ad.first_name,"") as dev_first_name,
			COALESCE(ad.last_name,"") as dev_last_name,
			COALESCE(ad.email,"") as dev_email,
			COALESCE(ad.status,"") as dev_status,
			COALESCE(ad.created_at,"") as dev_created_at,
			COALESCE(ad.created_by,"") as dev_created_by,
			COALESCE(ad.updated_at,"") as dev_updated_at,
			COALESCE(ad.updated_by,"") as dev_updated_by,

			COALESCE(a.id,"") as app_id,
			COALESCE(a.name,"") as app_name,
			COALESCE(a.access_type,"") as app_access_type,
			COALESCE(a.callback_url,"") as app_callback_url,
			COALESCE(a.display_name,"") as app_display_name,
			COALESCE(a.status,"") as app_status,
			COALESCE(a.app_family,"") as app_app_family,
			COALESCE(a.company_id,"") as app_company_id,
			COALESCE(a.created_at,"") as app_created_at,
			COALESCE(a.created_by,"") as app_created_by,
			COALESCE(a.updated_at,"") as app_updated_at,
			COALESCE(a.updated_by,"") as app_updated_by,

			COALESCE(ap.id,"") as prod_id,
			COALESCE(ap.name,"") as prod_name,
			COALESCE(ap.display_name,"") as prod_display_name,
			COALESCE(ap.quota,"") as prod_quota,
			COALESCE(ap.quota_interval, 0) as prod_quota_interval,
			COALESCE(ap.quota_time_unit,"") as prod_quota_time_unit,
			COALESCE(ap.created_at,"") as prod_created_at,
			COALESCE(ap.created_by,"") as prod_created_by,
			COALESCE(ap.updated_at,"") as prod_updated_at,
			COALESCE(ap.updated_by,"") as prod_updated_by,
			COALESCE(ap.proxies,"") as prod_proxies,
			COALESCE(ap.environments,"") as prod_environments,
			COALESCE(ap.api_resources,"") as prod_resources
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
		WHERE 	(mp.apiprdt_id = ap.id
			AND mp.app_id = a.id
			AND mp.appcred_id = c.id
			AND c.id = $1
			AND o.name = $2)
		UNION ALL
		SELECT
			COALESCE("company","") as ctype,
			COALESCE(c.tenant_id,""),

			COALESCE(c.status,""),
			COALESCE(c.consumer_secret,""),

			COALESCE(ad.id,"") as dev_id,
			COALESCE(ad.display_name,"") as dev_username,
			COALESCE("","") as dev_first_name,
			COALESCE("","") as dev_last_name,
			COALESCE("","") as dev_email,
			COALESCE(ad.status,"") as dev_status,
			COALESCE(ad.created_at,"") as dev_created_at,
			COALESCE(ad.created_by,"") as dev_created_by,
			COALESCE(ad.updated_at,"") as dev_updated_at,
			COALESCE(ad.updated_by,"") as dev_updated_by,

			COALESCE(a.id,"") as app_id,
			COALESCE(a.name,"") as app_name,
			COALESCE(a.access_type,"") as app_access_type,
			COALESCE(a.callback_url,"") as app_callback_url,
			COALESCE(a.display_name,"") as app_display_name,
			COALESCE(a.status,"") as app_status,
			COALESCE(a.app_family,"") as app_app_family,
			COALESCE(a.company_id,"") as app_company_id,
			COALESCE(a.created_at,"") as app_created_at,
			COALESCE(a.created_by,"") as app_created_by,
			COALESCE(a.updated_at,"") as app_updated_at,
			COALESCE(a.updated_by,"") as app_updated_by,

			COALESCE(ap.id,"") as prod_id,
			COALESCE(ap.name,"") as prod_name,
			COALESCE(ap.display_name,"") as prod_display_name,
			COALESCE(ap.quota,"") as prod_quota,
			COALESCE(ap.quota_interval,0) as prod_quota_interval,
			COALESCE(ap.quota_time_unit,"") as prod_quota_time_unit,
			COALESCE(ap.created_at,"") as prod_created_at,
			COALESCE(ap.created_by,"") as prod_created_by,
			COALESCE(ap.updated_at,"") as prod_updated_at,
			COALESCE(ap.updated_by,"") as prod_updated_by,
			COALESCE(ap.proxies,"") as prod_proxies,
			COALESCE(ap.environments,"") as prod_environments,
			COALESCE(ap.api_resources,"") as prod_resources

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
		WHERE   (mp.apiprdt_id = ap.id
			AND mp.app_id = a.id
			AND mp.appcred_id = c.id
			AND c.id = $1
			AND o.name = $2)
	;`

	//cid,csecret,did,dusername,dfirstname,dlastname,demail,dstatus,dcreated_at,dcreated_by,dlast_modified_at,dlast_modified_by, aid,aname,aaccesstype,acallbackurl,adisplay_name,astatus,aappfamily, acompany,acreated_at,acreated_by,alast_modified_at,alast_modified_by,pid,pname,pdisplayname,pquota_limit,pqutoainterval,pquotatimeout,pcreated_at,pcreated_by,plast_modified_at,plast_modified_by sql.NullString

	err := db.QueryRow(sSql, verifyApiKeyReq.Key, verifyApiKeyReq.OrganizationName).
		Scan(
			cType,
			tenantId,
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
			&resources,
		)

	if err != nil {
		log.Error("error fetching verify apikey details", err)
	}

	if err := json.Unmarshal([]byte(proxies), &apiProductDetails.Apiproxies); err != nil {
		log.Debug("unmarshall error for proxies, performing custom unmarshal ", proxies, err)

		apiProductDetails.Apiproxies = splitMalformedJson(proxies)

	}
	if err := json.Unmarshal([]byte(environments), &apiProductDetails.Environments); err != nil {
		log.Debug("unmarshall error for proxies, performing custom unmarshal ", environments, err)
		apiProductDetails.Environments = splitMalformedJson(environments)

	}
	if err := json.Unmarshal([]byte(resources), &apiProductDetails.Resources); err != nil {
		log.Debug("unmarshall error for proxies, performing custom unmarshal ", resources, err)
		apiProductDetails.Resources = splitMalformedJson(resources)

	}

	if appDetails.CallbackUrl != "" {
		clientIdDetails.RedirectURIs = []string{appDetails.CallbackUrl}
	}

	return err
}

func splitMalformedJson(fjson string) []string {
	var fs []string
	s := strings.TrimPrefix(fjson, "{")
	s = strings.TrimSuffix(s, "}")
	if utf8.RuneCountInString(s) > 0 {
		fs = strings.Split(s, ",")
	}
	return fs
}
