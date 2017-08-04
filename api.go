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
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type apiManagerInterface interface {
	InitAPI()
	handleRequest(w http.ResponseWriter, r *http.Request)
	verifyAPIKey(verifyApiKeyReq VerifyApiKeyRequest) ([]byte, error)
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
func (apiManager *apiManager) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	verifyApiKeyReq, err := validateRequest(r.Body, w)
	if err != nil {
		return
	}

	b, err := apiManager.verifyAPIKey(verifyApiKeyReq)

	if err != nil {
		respStatusCode, atoierr := strconv.Atoi(err.Error())
		if atoierr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(respStatusCode)
		}
		// TODO : discuss and finalize on error codes.
		w.WriteHeader(http.StatusBadRequest)
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
	if verifyApiKeyReq.Action == "" || verifyApiKeyReq.ApiProxyName == "" || verifyApiKeyReq.OrganizationName == "" || verifyApiKeyReq.EnvironmentName == "" || verifyApiKeyReq.Key == "" {
		// TODO : set correct fields in error response
		errorResponse := errorResponse("Bad_REQUEST", "Missing element")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse)
		return verifyApiKeyReq, errors.New("Bad_REQUEST")
	}
	return verifyApiKeyReq, nil
}

// returns []byte to be written to client
func (apiM apiManager) verifyAPIKey(verifyApiKeyReq VerifyApiKeyRequest) ([]byte, error) {

	dataWrapper := VerifyApiKeyRequestResponseDataWrapper{
		verifyApiKeyRequest: verifyApiKeyReq,
	}
	dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientId = verifyApiKeyReq.Key

	err := apiM.dbMan.getApiKeyDetails(&dataWrapper)

	switch {
	case err == sql.ErrNoRows:
		reason := "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode := "oauth.v2.InvalidApiKey"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusBadRequest))

	case err != nil:
		reason := err.Error()
		errorCode := "SEARCH_INTERNAL_ERROR"
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusInternalServerError))
	}

	/*
	 * Perform all validations
	 */
	errResponse, err := apiM.performValidations(dataWrapper)
	if errResponse != nil {
		return errResponse, err
	}

	apiM.enrichAttributes(&dataWrapper)

	if dataWrapper.ctype == "developer" {
		dataWrapper.verifyApiKeySuccessResponse.Developer = dataWrapper.tempDeveloperDetails
	} else {
		dataWrapper.verifyApiKeySuccessResponse.Company = CompanyDetails{
			Id:             dataWrapper.tempDeveloperDetails.Id,
			DisplayName:    dataWrapper.tempDeveloperDetails.UserName,
			Status:         dataWrapper.tempDeveloperDetails.Status,
			CreatedAt:      dataWrapper.tempDeveloperDetails.CreatedAt,
			CreatedBy:      dataWrapper.tempDeveloperDetails.CreatedBy,
			LastmodifiedAt: dataWrapper.tempDeveloperDetails.LastmodifiedAt,
			LastmodifiedBy: dataWrapper.tempDeveloperDetails.LastmodifiedBy,
			Attributes:     dataWrapper.tempDeveloperDetails.Attributes,
		}
	}

	resp := dataWrapper.verifyApiKeySuccessResponse

	return json.Marshal(resp)
}

func (apiM apiManager) performValidations(dataWrapper VerifyApiKeyRequestResponseDataWrapper) ([]byte, error) {
	clientIdDetails := dataWrapper.verifyApiKeySuccessResponse.ClientId
	verifyApiKeyReq := dataWrapper.verifyApiKeyRequest
	appDetails := dataWrapper.verifyApiKeySuccessResponse.App
	tempDeveloperDetails := dataWrapper.tempDeveloperDetails
	cType := dataWrapper.ctype
	apiProductDetails := dataWrapper.verifyApiKeySuccessResponse.ApiProduct

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
		if cType == "company" {
			errorCode = "keymanagement.service.CompanyStatusNotActive"
		}
		return errorResponse(reason, errorCode), errors.New(strconv.Itoa(http.StatusUnauthorized))
	}

	result := validatePath(apiProductDetails.Resources, verifyApiKeyReq.UriPath)
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

func (a *apiManager) enrichAttributes(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) {
	clientIdAttributes := a.dbMan.getKmsAttributes(dataWrapper.tenant_id, dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientId)
	developerAttributes := a.dbMan.getKmsAttributes(dataWrapper.tenant_id, dataWrapper.tempDeveloperDetails.Id)
	appAttributes := a.dbMan.getKmsAttributes(dataWrapper.tenant_id, dataWrapper.verifyApiKeySuccessResponse.App.Id)
	apiProductAttributes := a.dbMan.getKmsAttributes(dataWrapper.tenant_id, dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Id)

	dataWrapper.verifyApiKeySuccessResponse.ClientId.Attributes = clientIdAttributes
	dataWrapper.verifyApiKeySuccessResponse.App.Attributes = appAttributes
	dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Attributes = apiProductAttributes
	dataWrapper.tempDeveloperDetails.Attributes = developerAttributes
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
