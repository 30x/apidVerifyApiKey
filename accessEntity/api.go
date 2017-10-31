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

package accessEntity

import (
	"encoding/json"
	"fmt"
	"github.com/apid/apidVerifyApiKey/common"
	"net/http"
	"strconv"
	"strings"
)

const (
	AccessEntityPath         = "/entities"
	EndpointApp              = "/apps"
	EndpointApiProduct       = "/apiproducts"
	EndpointCompany          = "/companies"
	EndpointCompanyDeveloper = "/companydevelopers"
	EndpointDeveloper        = "/developers"
	EndpointAppCredentials   = "/appcredentials"
)

const (
	IdentifierAppId          = "appid"
	IdentifierApiProductName = "apiproductname"
	IdentifierAppName        = "appname"
	IdentifierApiResource    = "apiresource"
	IdentifierDeveloperId    = "developerid"
	IdentifierDeveloperEmail = "developeremail"
	IdentifierConsumerKey    = "consumerkey"
	IdentifierCompanyName    = "companyname"
)

var (
	Identifiers = map[string]bool{
		"appid":          true,
		"apiproductname": true,
		"appname":        true,
		"apiresource":    true,
		"developerid":    true,
		"companyname":    true,
		"developeremail": true,
		"consumerkey":    true,
	}

	ErrInvalidPar = &common.ErrorResponse{
		ResponseCode:    strconv.Itoa(INVALID_PARAMETERS),
		ResponseMessage: "invalid identifiers",
		StatusCode:      http.StatusBadRequest,
	}

	IdentifierTree = map[string]map[string][]string{
		EndpointApiProduct: {
			IdentifierApiProductName: {},
			IdentifierAppId:          {IdentifierApiResource},
			IdentifierAppName:        {IdentifierApiResource, IdentifierDeveloperEmail, IdentifierDeveloperId, IdentifierCompanyName},
			IdentifierConsumerKey:    {IdentifierApiResource},
		},
		EndpointApp: {
			IdentifierAppId:       {},
			IdentifierAppName:     {IdentifierDeveloperEmail, IdentifierDeveloperId, IdentifierCompanyName},
			IdentifierConsumerKey: {},
		},
		EndpointCompany: {
			IdentifierAppId:       {},
			IdentifierCompanyName: {},
			IdentifierConsumerKey: {},
		},
		EndpointCompanyDeveloper: {
			IdentifierCompanyName: {},
		},
		EndpointAppCredentials: {
			IdentifierConsumerKey: {},
		},
		EndpointDeveloper: {
			IdentifierDeveloperEmail: {},
			IdentifierAppId:          {},
			IdentifierDeveloperId:    {},
			IdentifierConsumerKey:    {},
		},
	}
)

const (
	INVALID_PARAMETERS = iota
	DB_ERROR
	DATA_ERROR
)

type ApiManager struct {
	DbMan            DbManagerInterface
	AccessEntityPath string
	apiInitialized   bool
}

func (a *ApiManager) InitAPI() {
	if a.apiInitialized {
		return
	}
	services.API().HandleFunc(a.AccessEntityPath+EndpointApp, a.HandleApps).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointApiProduct, a.HandleApiProducts).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointCompany, a.HandleCompanies).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointCompanyDeveloper, a.HandleCompanyDevelopers).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointDeveloper, a.HandleDevelopers).Methods("GET")
	services.API().HandleFunc(a.AccessEntityPath+EndpointAppCredentials, a.HandleConsumers).Methods("GET")
	a.apiInitialized = true
	log.Debug("API endpoints initialized")
}

func (a *ApiManager) HandleApps(w http.ResponseWriter, r *http.Request) {

}

func (a *ApiManager) HandleApiProducts(w http.ResponseWriter, r *http.Request) {
	ids, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		common.WriteError(w, err.Error(), INVALID_PARAMETERS, http.StatusBadRequest)
	}
	details, errRes := a.getApiProduct(ids)
	if errRes != nil {
		writeJson(errRes, w, r)
		return
	}
	writeJson(details, w, r)
}

func (a *ApiManager) HandleCompanies(w http.ResponseWriter, r *http.Request) {

}
func (a *ApiManager) HandleCompanyDevelopers(w http.ResponseWriter, r *http.Request) {

}
func (a *ApiManager) HandleDevelopers(w http.ResponseWriter, r *http.Request) {

}
func (a *ApiManager) HandleConsumers(w http.ResponseWriter, r *http.Request) {

}

func extractIdentifiers(pars map[string][]string) (map[string]string, error) {
	m := make(map[string]string)
	for k, v := range pars {
		k = strings.ToLower(k)
		if Identifiers[k] {
			if len(v) == 1 {
				m[k] = v[0]
			} else {
				return nil, fmt.Errorf("each identifier must have only 1 value")
			}
		}
	}
	return m, nil
}

func (a *ApiManager) getApiProduct(ids map[string]string) (*ApiProductDetails, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointApiProduct, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], keyVals[2], keyVals[3]

	prods, err := a.DbMan.GetApiProducts(priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getApiProduct: %v", err)
		return nil, newDbError(err)
	}

	var prod *common.ApiProduct
	var attrs []common.Attribute
	if len(prods) > 0 {
		prod = &prods[0]
		attrs = a.DbMan.GetKmsAttributes(prod.TenantId, prod.Id)[prod.Id]
	}

	return makeApiProductDetails(prod, attrs, priKey, priVal, secKey, secVal)
}

func makeApiProductDetails(prod *common.ApiProduct, attrs []common.Attribute, priKey, priVal, secKey, secVal string) (*ApiProductDetails, *common.ErrorResponse) {
	var a *ApiProductDetails
	if prod != nil {
		quotaLimit, err := strconv.Atoi(prod.Quota)
		if err != nil {
			return nil, newDataError(err)
		}
		a = &ApiProductDetails{
			APIProxies:     common.JsonToStringArray(prod.Proxies),
			APIResources:   common.JsonToStringArray(prod.ApiResources),
			ApprovalType:   prod.ApprovalType,
			Attributes:     attrs,
			CreatedAt:      prod.CreatedAt,
			CreatedBy:      prod.CreatedBy,
			Description:    prod.Description,
			DisplayName:    prod.DisplayName,
			Environments:   common.JsonToStringArray(prod.Environments),
			ID:             prod.Id,
			LastModifiedAt: prod.UpdatedAt,
			LastModifiedBy: prod.UpdatedBy,
			Name:           prod.Name,
			QuotaInterval:  prod.QuotaInterval,
			QuotaLimit:     int64(quotaLimit),
			QuotaTimeUnit:  prod.QuotaTimeUnit,
			Scopes:         common.JsonToStringArray(prod.Scopes),
		}
	} else {
		a = new(ApiProductDetails)
	}

	setResIdentifiers(a, priKey, priVal, secKey, secVal)
	return a, nil
}

func setResIdentifiers(a *ApiProductDetails, priKey, priVal, secKey, secVal string) {
	a.PrimaryIdentifierType = priKey
	a.PrimaryIdentifierValue = priVal
	a.SecondaryIdentifierType = secKey
	a.SecondaryIdentifierValue = secVal
}

func parseIdentifiers(endpoint string, ids map[string]string) (valid bool, keyVals []string) {
	if len(ids) > 2 {
		return false, nil
	}
	if m := IdentifierTree[endpoint]; m != nil {
		for key, val := range ids {
			if m[key] != nil {
				keyVals = append(keyVals, key, val)
				for _, id := range m[key] {
					if ids[id] != "" {
						keyVals = append(keyVals, id, ids[id])
						return true, keyVals
					}
				}
				if len(ids) == 2 {
					return false, nil
				}
				keyVals = append(keyVals, "", "")
				return true, keyVals
			}
		}
	}
	return false, nil
}

func newDbError(err error) *common.ErrorResponse {
	return &common.ErrorResponse{
		ResponseCode:    strconv.Itoa(DB_ERROR),
		ResponseMessage: err.Error(),
		StatusCode:      http.StatusInternalServerError,
	}
}

func newDataError(err error) *common.ErrorResponse {
	return &common.ErrorResponse{
		ResponseCode:    strconv.Itoa(DATA_ERROR),
		ResponseMessage: err.Error(),
		StatusCode:      http.StatusInternalServerError,
	}
}

func writeJson(obj interface{}, w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		log.Error("unable to marshal errorResponse: " + err.Error())
		w.Write([]byte("unable to marshal errorResponse: " + err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}
