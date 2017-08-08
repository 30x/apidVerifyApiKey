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
func (a *apiManager) handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	verifyApiKeyReq, err := validateRequest(r.Body, w)
	if err != nil {
		return
	}

	b, err := a.verifyAPIKey(verifyApiKeyReq)

	if err != nil {
		respStatusCode, atoierr := strconv.Atoi(err.Error())
		if atoierr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(respStatusCode)
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
	// TODO : make this method of verifyApiKeyReq struct
	// TODO : move validation to verifyApiKey struct validate method
	if verifyApiKeyReq.Action == "" || verifyApiKeyReq.ApiProxyName == "" || verifyApiKeyReq.OrganizationName == "" || verifyApiKeyReq.EnvironmentName == "" || verifyApiKeyReq.Key == "" {
		// TODO : set correct missing fields in error response
		errorResponse, _ := json.Marshal(errorResponse("Bad_REQUEST", "Missing element", http.StatusBadRequest))
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
	case err != nil && err.Error() == "InvalidApiKey":
		reason := "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode := "oauth.v2.InvalidApiKey"
		errResponse := errorResponse(reason, errorCode, http.StatusOK)
		return json.Marshal(errResponse)

	case err != nil:
		reason := err.Error()
		errorCode := "SEARCH_INTERNAL_ERROR"
		errResponse := errorResponse(reason, errorCode, http.StatusInternalServerError)
		return json.Marshal(errResponse)
	}

	dataWrapper.verifyApiKeySuccessResponse.ApiProduct = shortListApiProduct(dataWrapper.apiProducts, verifyApiKeyReq)
	/*
	 * Perform all validations
	 */
	errResponse := apiM.performValidations(dataWrapper)
	if errResponse != nil {
		return json.Marshal(&errResponse)
	}

	apiM.enrichAttributes(&dataWrapper)

	setDevOrCompanyInResponseBasedOnCtype(dataWrapper.ctype, dataWrapper.tempDeveloperDetails, &dataWrapper.verifyApiKeySuccessResponse)

	resp := dataWrapper.verifyApiKeySuccessResponse

	return json.Marshal(resp)
}

func setDevOrCompanyInResponseBasedOnCtype(ctype string, tempDeveloperDetails DeveloperDetails, response *VerifyApiKeySuccessResponse) {
	if ctype == "developer" {
		response.Developer = tempDeveloperDetails
	} else {
		response.Company = CompanyDetails{
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
}

func shortListApiProduct(details []ApiProductDetails, verifyApiKeyReq VerifyApiKeyRequest) ApiProductDetails {
	var bestMathcedProduct ApiProductDetails
	rankedProducts := make(map[int][]ApiProductDetails)
	rankedProducts[2] = []ApiProductDetails{}
	rankedProducts[3] = []ApiProductDetails{}

	for _, apiProd := range details {
		if len(apiProd.Resources) == 0 || validatePath(apiProd.Resources, verifyApiKeyReq.UriPath) {
			if len(apiProd.Apiproxies) == 0 || contains(apiProd.Apiproxies, verifyApiKeyReq.ApiProxyName) {
				if len(apiProd.Environments) == 0 || contains(apiProd.Environments, verifyApiKeyReq.EnvironmentName) {
					bestMathcedProduct = apiProd
					return bestMathcedProduct
					// set rank 1 or just return
				} else {
					// set rank to 2
					rankedProducts[2] = append(rankedProducts[2], apiProd)
				}
			} else {
				// set rank to 3,
				rankedProducts[3] = append(rankedProducts[3], apiProd)
			}
		}
	}

	if len(rankedProducts[2]) > 0 {
		return rankedProducts[2][0]
	} else if len(rankedProducts[3]) > 0 {
		return rankedProducts[3][0]
	}

	return bestMathcedProduct

}

func (apiM apiManager) performValidations(dataWrapper VerifyApiKeyRequestResponseDataWrapper) *ErrorResponse {
	clientIdDetails := dataWrapper.verifyApiKeySuccessResponse.ClientId
	verifyApiKeyReq := dataWrapper.verifyApiKeyRequest
	appDetails := dataWrapper.verifyApiKeySuccessResponse.App
	tempDeveloperDetails := dataWrapper.tempDeveloperDetails
	cType := dataWrapper.ctype
	apiProductDetails := dataWrapper.verifyApiKeySuccessResponse.ApiProduct
	var reason, errorCode string

	if !strings.EqualFold("APPROVED", clientIdDetails.Status) {
		reason = "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode = "oauth.v2.ApiKeyNotApproved"
	}

	if !strings.EqualFold("APPROVED", appDetails.Status) {
		reason = "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode = "keymanagement.service.invalid_client-app_not_approved"
	}

	if !strings.EqualFold("ACTIVE", tempDeveloperDetails.Status) {
		reason = "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode = "keymanagement.service.DeveloperStatusNotActive"
		if cType == "company" {
			errorCode = "keymanagement.service.CompanyStatusNotActive"
		}
	}

	if dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Id == "" {
		reason = "Path Validation Failed. Product not resolved"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"
	}

	result := len(apiProductDetails.Resources) == 0 || validatePath(apiProductDetails.Resources, verifyApiKeyReq.UriPath)
	if !result {
		reason = "Path Validation Failed (" + strings.Join(apiProductDetails.Resources, ", ") + " vs " + verifyApiKeyReq.UriPath + ")"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"
	}

	if verifyApiKeyReq.ValidateAgainstApiProxiesAndEnvs && (len(apiProductDetails.Apiproxies) > 0 && !contains(apiProductDetails.Apiproxies, verifyApiKeyReq.ApiProxyName)) {
		reason = "Proxy Validation Failed (" + strings.Join(apiProductDetails.Apiproxies, ", ") + " vs " + verifyApiKeyReq.ApiProxyName + ")"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"
	}
	/* Verify if the ENV matches */
	if verifyApiKeyReq.ValidateAgainstApiProxiesAndEnvs && (len(apiProductDetails.Environments) > 0 && !contains(apiProductDetails.Environments, verifyApiKeyReq.EnvironmentName)) {
		reason = "ENV Validation Failed (" + strings.Join(apiProductDetails.Environments, ", ") + " vs " + verifyApiKeyReq.EnvironmentName + ")"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"

	}

	if errorCode != "" {
		log.Debug("Validation error occoured ", errorCode, " ", reason)
		ee := errorResponse(reason, errorCode, http.StatusOK)
		return &ee
	}

	return nil

}

func (a *apiManager) enrichAttributes(dataWrapper *VerifyApiKeyRequestResponseDataWrapper) {

	attributeMap := a.dbMan.getKmsAttributes(dataWrapper.tenant_id, dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientId, dataWrapper.tempDeveloperDetails.Id, dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Id, dataWrapper.verifyApiKeySuccessResponse.App.Id)

	clientIdAttributes := attributeMap[dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientId]
	developerAttributes := attributeMap[dataWrapper.tempDeveloperDetails.Id]
	appAttributes := attributeMap[dataWrapper.verifyApiKeySuccessResponse.App.Id]
	apiProductAttributes := attributeMap[dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Id]

	dataWrapper.verifyApiKeySuccessResponse.ClientId.Attributes = clientIdAttributes
	dataWrapper.verifyApiKeySuccessResponse.App.Attributes = appAttributes
	dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Attributes = apiProductAttributes
	dataWrapper.tempDeveloperDetails.Attributes = developerAttributes
}

func errorResponse(reason, errorCode string, statusCode int) ErrorResponse {
	if errorCode == "SEARCH_INTERNAL_ERROR" {
		log.Error(reason)
	} else {
		log.Debug(reason)
	}
	resp := ErrorResponse{
		ResponseCode:    errorCode,
		ResponseMessage: reason,
		StatusCode:      statusCode,
	}
	return resp
}
