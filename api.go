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
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiManagerInterface interface {
	InitAPI()
	handleRequest(w http.ResponseWriter, r *http.Request)
	verifyAPIKey(verifyApiKeyReq VerifyApiKeyRequest) (*VerifyApiKeySuccessResponse, *ErrorResponse)
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

	var returnValue interface{}

	verifyApiKeyReq, err := validateRequest(r.Body, w)
	if err != nil {
		errorResponse, jsonErr := json.Marshal(errorResponse("Bad_REQUEST", err.Error(), http.StatusBadRequest))
		if jsonErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(jsonErr.Error()))
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse)
		return
	}

	verifyApiKeyResponse, errorResponse := a.verifyAPIKey(verifyApiKeyReq)

	if errorResponse != nil {
		setResponseHeader(errorResponse, w)
		returnValue = errorResponse
	} else {
		returnValue = verifyApiKeyResponse
	}
	b, _ := json.Marshal(returnValue)
	log.Debugf("handleVerifyAPIKey result %s", b)
	w.Write(b)

}

func setResponseHeader(errorResponse *ErrorResponse, w http.ResponseWriter) {
	if errorResponse.StatusCode != 0 {
		w.WriteHeader(errorResponse.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func validateRequest(requestBody io.ReadCloser, w http.ResponseWriter) (VerifyApiKeyRequest, error) {
	defer requestBody.Close()
	// 1. read request boby
	var verifyApiKeyReq VerifyApiKeyRequest
	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return verifyApiKeyReq, err
	}
	log.Debug("request body: ", string(body))
	// 2. umarshall json to struct
	err = json.Unmarshal(body, &verifyApiKeyReq)
	if err != nil {
		return verifyApiKeyReq, err
	}
	log.Debug(verifyApiKeyReq)

	// 2. verify params
	if isValid, err := verifyApiKeyReq.validate(); !isValid {
		return verifyApiKeyReq, err
	}
	return verifyApiKeyReq, nil
}

// returns []byte to be written to client
func (apiM apiManager) verifyAPIKey(verifyApiKeyReq VerifyApiKeyRequest) (*VerifyApiKeySuccessResponse, *ErrorResponse) {

	dataWrapper := VerifyApiKeyRequestResponseDataWrapper{
		verifyApiKeyRequest: verifyApiKeyReq,
	}
	dataWrapper.verifyApiKeySuccessResponse.ClientId.ClientId = verifyApiKeyReq.Key
	dataWrapper.verifyApiKeySuccessResponse.Environment = verifyApiKeyReq.EnvironmentName

	err := apiM.dbMan.getApiKeyDetails(&dataWrapper)

	switch {
	case err != nil && err.Error() == "InvalidApiKey":
		reason := "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode := "oauth.v2.InvalidApiKey"
		errResponse := errorResponse(reason, errorCode, http.StatusOK)
		return nil, &errResponse

	case err != nil:
		reason := err.Error()
		errorCode := "SEARCH_INTERNAL_ERROR"
		errResponse := errorResponse(reason, errorCode, http.StatusInternalServerError)
		return nil, &errResponse
	}

	dataWrapper.verifyApiKeySuccessResponse.ApiProduct = shortListApiProduct(dataWrapper.apiProducts, verifyApiKeyReq)
	/*
	 * Perform all validations
	 */
	errResponse := apiM.performValidations(dataWrapper)
	if errResponse != nil {
		return nil, errResponse
	}

	apiM.enrichAttributes(&dataWrapper)

	setDevOrCompanyInResponseBasedOnCtype(dataWrapper.ctype, dataWrapper.tempDeveloperDetails, &dataWrapper.verifyApiKeySuccessResponse)

	return &dataWrapper.verifyApiKeySuccessResponse, nil
}

func setDevOrCompanyInResponseBasedOnCtype(ctype string, tempDeveloperDetails DeveloperDetails, response *VerifyApiKeySuccessResponse) {
	if ctype == "developer" {
		response.Developer = tempDeveloperDetails
	} else {
		response.Company = CompanyDetails{
			Id:             tempDeveloperDetails.Id,
			Name:           tempDeveloperDetails.FirstName,
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
	rankedProducts := make([][]ApiProductDetails, 2)

	for _, apiProd := range details {
		if len(apiProd.Resources) == 0 || validatePath(apiProd.Resources, verifyApiKeyReq.UriPath) {
			if len(apiProd.Apiproxies) == 0 || contains(apiProd.Apiproxies, verifyApiKeyReq.ApiProxyName) {
				if len(apiProd.Environments) == 0 || contains(apiProd.Environments, verifyApiKeyReq.EnvironmentName) {
					bestMathcedProduct = apiProd
					return bestMathcedProduct
					// set rank 1 or just return
				} else {
					// set rank to 2
					rankedProducts[0] = append(rankedProducts[0], apiProd)
				}
			} else {
				// set rank to 3,
				rankedProducts[1] = append(rankedProducts[1], apiProd)
			}
		}
	}

	if len(rankedProducts[0]) > 0 {
		return rankedProducts[0][0]
	}

	if len(rankedProducts[1]) > 0 {
		return rankedProducts[1][0]
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
		log.Debug("Validation error occoured ", errorCode, " ", reason)
		ee := errorResponse(reason, errorCode, http.StatusOK)
		return &ee
	}

	if !strings.EqualFold("APPROVED", appDetails.Status) {
		reason = "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode = "keymanagement.service.invalid_client-app_not_approved"
		log.Debug("Validation error occoured ", errorCode, " ", reason)
		ee := errorResponse(reason, errorCode, http.StatusOK)
		return &ee
	}

	if !strings.EqualFold("ACTIVE", tempDeveloperDetails.Status) {
		reason = "API Key verify failed for (" + verifyApiKeyReq.Key + ", " + verifyApiKeyReq.OrganizationName + ")"
		errorCode = "keymanagement.service.DeveloperStatusNotActive"
		if cType == "company" {
			errorCode = "keymanagement.service.CompanyStatusNotActive"
		}
		log.Debug("Validation error occoured ", errorCode, " ", reason)
		ee := errorResponse(reason, errorCode, http.StatusOK)
		return &ee
	}

	if dataWrapper.verifyApiKeySuccessResponse.ApiProduct.Id == "" {
		reason = "Path Validation Failed. Product not resolved"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"
		log.Debug("Validation error occoured ", errorCode, " ", reason)
		ee := errorResponse(reason, errorCode, http.StatusOK)
		return &ee
	}

	result := len(apiProductDetails.Resources) == 0 || validatePath(apiProductDetails.Resources, verifyApiKeyReq.UriPath)
	if !result {
		reason = "Path Validation Failed (" + strings.Join(apiProductDetails.Resources, ", ") + " vs " + verifyApiKeyReq.UriPath + ")"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"
		log.Debug("Validation error occoured ", errorCode, " ", reason)
		ee := errorResponse(reason, errorCode, http.StatusOK)
		return &ee
	}

	if verifyApiKeyReq.ValidateAgainstApiProxiesAndEnvs && (len(apiProductDetails.Apiproxies) > 0 && !contains(apiProductDetails.Apiproxies, verifyApiKeyReq.ApiProxyName)) {
		reason = "Proxy Validation Failed (" + strings.Join(apiProductDetails.Apiproxies, ", ") + " vs " + verifyApiKeyReq.ApiProxyName + ")"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"
		log.Debug("Validation error occoured ", errorCode, " ", reason)
		ee := errorResponse(reason, errorCode, http.StatusOK)
		return &ee
	}
	/* Verify if the ENV matches */
	if verifyApiKeyReq.ValidateAgainstApiProxiesAndEnvs && (len(apiProductDetails.Environments) > 0 && !contains(apiProductDetails.Environments, verifyApiKeyReq.EnvironmentName)) {
		reason = "ENV Validation Failed (" + strings.Join(apiProductDetails.Environments, ", ") + " vs " + verifyApiKeyReq.EnvironmentName + ")"
		errorCode = "oauth.v2.InvalidApiKeyForGivenResource"
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
