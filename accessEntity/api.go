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
	"github.com/apid/apidApiMetadata/common"
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
	IdentifierOrganization   = "organization"
)

const (
	TypeDeveloper   = "developer"
	TypeCompany     = "company"
	TypeApp         = "app"
	TypeConsumerKey = "consumerkey"
)

const (
	AppTypeDeveloper = "DEVELOPER"
	AppTypeCompany   = "COMPANY"
)

const (
	StatusApproved = "APPROVED"
	StatusRevoked  = "REVOKED"
	StatusExpired  = "EXPIRED"
)

const headerRequestId = "X-Gateway-Request-Id"

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
		ResponseMessage: "Invalid Identifiers",
		StatusCode:      http.StatusBadRequest,
	}

	ErrNotFound = &common.ErrorResponse{
		ResponseCode:    strconv.Itoa(NOT_FOUND),
		ResponseMessage: "Resource Not Found",
		StatusCode:      http.StatusNotFound,
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
	// Server DB Error
	DB_ERROR
	// Invalid/Wrong Data in DB data. This probably means something wrong happened in upstream PG/Transicator.
	DATA_ERROR
	// 404
	NOT_FOUND
	// json Marshal Error
	JSON_MARSHAL_ERROR
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
	services.API().HandleFunc(a.AccessEntityPath+EndpointAppCredentials, a.HandleAppCredentials).Methods("GET")
	a.apiInitialized = true
	log.Debug("API endpoints initialized")
}

func (a *ApiManager) handleEndpoint(endpoint string, w http.ResponseWriter, r *http.Request) {
	ids, org, err := extractIdentifiers(r.URL.Query())
	if err != nil {
		writeJson(http.StatusBadRequest,
			common.ErrorResponse{
				ResponseCode:    strconv.Itoa(INVALID_PARAMETERS),
				ResponseMessage: err.Error(),
				StatusCode:      http.StatusBadRequest,
			}, w, r)
	}
	var res interface{}
	var errRes *common.ErrorResponse
	switch endpoint {
	case EndpointApp:
		res, errRes = a.getApp(org, ids)
	case EndpointApiProduct:
		res, errRes = a.getApiProduct(org, ids)
	case EndpointCompany:
		res, errRes = a.getCompany(org, ids)
	case EndpointCompanyDeveloper:
		res, errRes = a.getCompanyDeveloper(org, ids)
	case EndpointDeveloper:
		res, errRes = a.getDeveloper(org, ids)
	case EndpointAppCredentials:
		res, errRes = a.getAppCredential(org, ids)
	}

	if errRes != nil {
		writeJson(errRes.StatusCode, errRes, w, r)
		return
	}
	writeJson(http.StatusOK, res, w, r)
}

func (a *ApiManager) HandleApps(w http.ResponseWriter, r *http.Request) {
	a.handleEndpoint(EndpointApp, w, r)
}

func (a *ApiManager) HandleApiProducts(w http.ResponseWriter, r *http.Request) {
	a.handleEndpoint(EndpointApiProduct, w, r)
}

func (a *ApiManager) HandleCompanies(w http.ResponseWriter, r *http.Request) {
	a.handleEndpoint(EndpointCompany, w, r)
}

func (a *ApiManager) HandleCompanyDevelopers(w http.ResponseWriter, r *http.Request) {
	a.handleEndpoint(EndpointCompanyDeveloper, w, r)
}
func (a *ApiManager) HandleDevelopers(w http.ResponseWriter, r *http.Request) {
	a.handleEndpoint(EndpointDeveloper, w, r)
}

func (a *ApiManager) HandleAppCredentials(w http.ResponseWriter, r *http.Request) {
	a.handleEndpoint(EndpointAppCredentials, w, r)
}

func extractIdentifiers(pars map[string][]string) (map[string]string, string, error) {
	m := make(map[string]string)
	orgs := pars[IdentifierOrganization]
	if len(orgs) == 0 {
		return nil, "", fmt.Errorf("no org specified")
	}
	org := orgs[0]
	for k, v := range pars {
		k = strings.ToLower(k)
		if Identifiers[k] {
			if len(v) == 1 {
				m[k] = v[0]
			} else {
				return nil, org, fmt.Errorf("each identifier must have only 1 value")
			}
		}
	}
	return m, org, nil
}

func (a *ApiManager) getCompanyDeveloper(org string, ids map[string]string) (*CompanyDevelopersSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointCompanyDeveloper, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal := keyVals[0], keyVals[1]

	devs, err := a.DbMan.GetCompanyDevelopers(org, priKey, priVal, "", "")
	if err != nil {
		log.Errorf("getCompanyDeveloper: %v", err)
		return nil, newDbError(err)
	}

	if len(devs) == 0 {
		return nil, ErrNotFound
	}

	var details []*CompanyDeveloperDetails
	for _, dev := range devs {
		comName, err := a.DbMan.GetComNames(dev.CompanyId, TypeCompany)
		if err != nil || len(comName) == 0 {
			log.Errorf("getCompanyDeveloper: %v", err)
			return nil, newDbError(err)
		}
		email, err := a.DbMan.GetDevEmailByDevId(dev.DeveloperId, org)
		if err != nil {
			log.Errorf("getCompanyDeveloper: %v", err)
			return nil, newDbError(err)
		}
		detail := makeComDevDetails(&dev, comName[0], email)
		details = append(details, detail)
	}
	return &CompanyDevelopersSuccessResponse{
		CompanyDevelopers:      details,
		Organization:           org,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
	}, nil
}

func (a *ApiManager) getDeveloper(org string, ids map[string]string) (*DeveloperSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointDeveloper, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal := keyVals[0], keyVals[1]

	devs, err := a.DbMan.GetDevelopers(org, priKey, priVal, "", "")
	if err != nil {
		log.Errorf("getDeveloper: %v", err)
		return nil, newDbError(err)
	}

	if len(devs) == 0 {
		return nil, ErrNotFound
	}
	dev := &devs[0]

	attrs := a.DbMan.GetKmsAttributes(dev.TenantId, dev.Id)[dev.Id]
	comNames, err := a.DbMan.GetComNames(dev.Id, TypeDeveloper)
	if err != nil {
		log.Errorf("getDeveloper: %v", err)
		return nil, newDbError(err)
	}
	appNames, err := a.DbMan.GetAppNames(dev.Id, TypeDeveloper)
	if err != nil {
		log.Errorf("getDeveloper: %v", err)
		return nil, newDbError(err)
	}
	details := makeDevDetails(dev, appNames, comNames, attrs)
	return &DeveloperSuccessResponse{
		Developer:              details,
		Organization:           org,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
	}, nil
}

func (a *ApiManager) getCompany(org string, ids map[string]string) (*CompanySuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointCompany, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal := keyVals[0], keyVals[1]

	coms, err := a.DbMan.GetCompanies(org, priKey, priVal, "", "")
	if err != nil {
		log.Errorf("getCompany: %v", err)
		return nil, newDbError(err)
	}

	if len(coms) == 0 {
		return nil, ErrNotFound
	}
	com := &coms[0]

	attrs := a.DbMan.GetKmsAttributes(com.TenantId, com.Id)[com.Id]
	appNames, err := a.DbMan.GetAppNames(com.Id, TypeCompany)
	if err != nil {
		log.Errorf("getCompany: %v", err)
		return nil, newDbError(err)
	}
	details := makeCompanyDetails(com, appNames, attrs)
	return &CompanySuccessResponse{
		Company:                details,
		Organization:           org,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
	}, nil
}

func (a *ApiManager) getApiProduct(org string, ids map[string]string) (*ApiProductSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointApiProduct, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], "", ""
	if len(keyVals) > 2 {
		secKey, secVal = keyVals[2], keyVals[3]
	}
	prods, err := a.DbMan.GetApiProducts(org, priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getApiProduct: %v", err)
		return nil, newDbError(err)
	}

	var attrs []common.Attribute
	if len(prods) == 0 {
		return nil, ErrNotFound
	}
	prod := &prods[0]
	attrs = a.DbMan.GetKmsAttributes(prod.TenantId, prod.Id)[prod.Id]
	details, errRes := makeApiProductDetails(prod, attrs)
	if errRes != nil {
		return nil, errRes
	}

	return &ApiProductSuccessResponse{
		ApiProduct:               details,
		Organization:             org,
		PrimaryIdentifierType:    priKey,
		PrimaryIdentifierValue:   priVal,
		SecondaryIdentifierType:  secKey,
		SecondaryIdentifierValue: secVal,
	}, nil
}

func (a *ApiManager) getAppCredential(org string, ids map[string]string) (*AppCredentialSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointApiProduct, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal := keyVals[0], keyVals[1]

	appCreds, err := a.DbMan.GetAppCredentials(org, priKey, priVal, "", "")
	if err != nil {
		log.Errorf("getAppCredential: %v", err)
		return nil, newDbError(err)
	}

	if len(appCreds) == 0 {
		return nil, ErrNotFound
	}
	appCred := &appCreds[0]
	attrs := a.DbMan.GetKmsAttributes(appCred.TenantId, appCred.Id)[appCred.Id]
	apps, err := a.DbMan.GetApps(org, IdentifierAppId, appCred.AppId, "", "")
	if err != nil {
		log.Errorf("getAppCredential: %v", err)
		return nil, newDbError(err)
	}

	if len(apps) == 0 {
		log.Errorf("getAppCredential: No App with id=%v", appCred.AppId)
		return &AppCredentialSuccessResponse{
			AppCredential: nil,
			Organization:  org,
		}, nil
	}
	app := &apps[0]
	cd, errRes := a.getCredDetails(appCred, app.Status)
	if errRes != nil {
		return nil, errRes
	}
	devStatus := ""
	if app.DeveloperId != "" {
		devStatus, err = a.DbMan.GetStatus(app.DeveloperId, AppTypeDeveloper)
		if err != nil {
			log.Errorf("getAppCredential error get status: %v", err)
			return nil, newDbError(err)
		}
	}
	cks := makeConsumerKeyStatusDetails(app, cd, devStatus)
	details := makeAppCredentialDetails(appCred, cks, []string{app.CallbackUrl}, attrs)
	return &AppCredentialSuccessResponse{
		AppCredential:          details,
		Organization:           org,
		PrimaryIdentifierType:  priKey,
		PrimaryIdentifierValue: priVal,
	}, nil
}

func (a *ApiManager) getApp(org string, ids map[string]string) (*AppSuccessResponse, *common.ErrorResponse) {
	valid, keyVals := parseIdentifiers(EndpointApp, ids)
	if !valid {
		return nil, ErrInvalidPar
	}
	priKey, priVal, secKey, secVal := keyVals[0], keyVals[1], "", ""
	if len(keyVals) > 2 {
		secKey, secVal = keyVals[2], keyVals[3]
	}

	apps, err := a.DbMan.GetApps(org, priKey, priVal, secKey, secVal)
	if err != nil {
		log.Errorf("getApp: %v", err)
		return nil, newDbError(err)
	}

	var app *common.App
	var attrs []common.Attribute

	if len(apps) == 0 {
		return nil, ErrNotFound
	}

	app = &apps[0]
	attrs = a.DbMan.GetKmsAttributes(app.TenantId, app.Id)[app.Id]
	prods, err := a.DbMan.GetApiProductNames(app.Id, TypeApp)
	if err != nil {
		log.Errorf("getApp error getting productNames: %v", err)
		return nil, newDbError(err)
	}
	parStatus, err := a.DbMan.GetStatus(app.ParentId, app.Type)
	if err != nil {
		log.Errorf("getApp error getting parent status: %v", err)
		return nil, newDbError(err)
	}
	creds, err := a.DbMan.GetAppCredentials(org, IdentifierAppId, app.Id, "", "")
	if err != nil {
		log.Errorf("getApp error getting parent status: %v", err)
		return nil, newDbError(err)
	}
	var credDetails []*CredentialDetails
	for _, cred := range creds {
		detail, errRes := a.getCredDetails(&cred, app.Status)
		if errRes != nil {
			return nil, errRes
		}
		credDetails = append(credDetails, detail)
	}

	parent, errRes := a.getAppParent(app.ParentId, app.Type)
	if errRes != nil {
		return nil, errRes
	}
	details, errRes := makeAppDetails(app, parent, parStatus, prods, credDetails, attrs)
	if errRes != nil {
		return nil, errRes
	}
	return &AppSuccessResponse{
		App:                      details,
		Organization:             org,
		PrimaryIdentifierType:    priKey,
		PrimaryIdentifierValue:   priVal,
		SecondaryIdentifierType:  secKey,
		SecondaryIdentifierValue: secVal,
	}, nil
}

func (a *ApiManager) getAppParent(id string, parentType string) (string, *common.ErrorResponse) {
	switch parentType {
	case AppTypeDeveloper:
		return id, nil
	case AppTypeCompany:
		names, err := a.DbMan.GetComNames(id, TypeCompany)
		if err != nil {
			return "", newDbError(err)
		}
		if len(names) == 0 {
			log.Warnf("getAppParent: No company with id=%v", id)
			return "", nil
		}
		return names[0], nil
	}
	return "", nil
}

func makeConsumerKeyStatusDetails(app *common.App, c *CredentialDetails, devStatus string) *ConsumerKeyStatusDetails {
	return &ConsumerKeyStatusDetails{
		AppCredential:   c,
		AppID:           c.AppID,
		AppName:         app.Name,
		AppStatus:       app.Status,
		AppType:         app.Type,
		DeveloperID:     app.DeveloperId,
		DeveloperStatus: devStatus,
		IsValidKey:      strings.EqualFold(c.Status, StatusApproved),
	}
}

func makeAppCredentialDetails(ac *common.AppCredential, cks *ConsumerKeyStatusDetails, redirectUrl []string, attrs []common.Attribute) *AppCredentialDetails {
	return &AppCredentialDetails{
		AppID:             ac.AppId,
		AppName:           cks.AppName,
		Attributes:        attrs,
		ConsumerKey:       ac.Id,
		ConsumerKeyStatus: cks,
		ConsumerSecret:    ac.ConsumerSecret,
		DeveloperID:       cks.DeveloperID,
		RedirectUris:      redirectUrl,
		Scopes:            common.JsonToStringArray(ac.Scopes),
		Status:            ac.Status,
	}
}

func makeApiProductDetails(prod *common.ApiProduct, attrs []common.Attribute) (*ApiProductDetails, *common.ErrorResponse) {
	var a *ApiProductDetails
	if prod != nil {
		var quotaLimit int
		var err error
		if prod.Quota != "" {
			quotaLimit, err = strconv.Atoi(prod.Quota)
			if err != nil {
				return nil, newDataError(err)
			}
		}

		a = &ApiProductDetails{
			ApiProxies:     common.JsonToStringArray(prod.Proxies),
			ApiResources:   common.JsonToStringArray(prod.ApiResources),
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
	return a, nil
}

func makeAppDetails(app *common.App, parent string, parentStatus string, prods []string, creds []*CredentialDetails, attrs []common.Attribute) (*AppDetails, *common.ErrorResponse) {
	var a *AppDetails
	if app != nil {
		a = &AppDetails{
			AccessType:      app.AccessType,
			ApiProducts:     prods,
			AppCredentials:  creds,
			AppFamily:       app.AppFamily,
			AppParentID:     parent,
			AppParentStatus: parentStatus,
			AppType:         app.Type,
			Attributes:      attrs,
			CallbackUrl:     app.CallbackUrl,
			CreatedAt:       app.CreatedAt,
			CreatedBy:       app.CreatedBy,
			DisplayName:     app.DisplayName,
			Id:              app.Id,
			LastModifiedAt:  app.UpdatedAt,
			LastModifiedBy:  app.UpdatedBy,
			Name:            app.Name,
			Status:          app.Status,
		}
	} else {
		a = new(AppDetails)
	}
	return a, nil
}

func makeCompanyDetails(com *common.Company, appNames []string, attrs []common.Attribute) *CompanyDetails {
	return &CompanyDetails{
		Apps:           appNames,
		Attributes:     attrs,
		CreatedAt:      com.CreatedAt,
		CreatedBy:      com.CreatedBy,
		DisplayName:    com.DisplayName,
		ID:             com.Id,
		LastModifiedAt: com.UpdatedAt,
		LastModifiedBy: com.UpdatedBy,
		Name:           com.Name,
		Status:         com.Status,
	}
}

func makeDevDetails(dev *common.Developer, appNames []string, comNames []string, attrs []common.Attribute) *DeveloperDetails {
	return &DeveloperDetails{
		Apps:           appNames,
		Attributes:     attrs,
		Companies:      comNames,
		CreatedAt:      dev.CreatedAt,
		CreatedBy:      dev.CreatedBy,
		Email:          dev.Email,
		FirstName:      dev.FirstName,
		ID:             dev.Id,
		LastModifiedAt: dev.UpdatedAt,
		LastModifiedBy: dev.UpdatedBy,
		LastName:       dev.LastName,
		Password:       dev.Password,
		Status:         dev.Status,
		UserName:       dev.UserName,
	}
}

func makeComDevDetails(comDev *common.CompanyDeveloper, comName, devEmail string) *CompanyDeveloperDetails {
	return &CompanyDeveloperDetails{
		CompanyName:    comName,
		CreatedAt:      comDev.CreatedAt,
		CreatedBy:      comDev.CreatedBy,
		DeveloperEmail: devEmail,
		LastModifiedAt: comDev.UpdatedAt,
		LastModifiedBy: comDev.UpdatedBy,
		Roles:          common.JsonToStringArray(comDev.Roles),
	}
}

func (a *ApiManager) getCredDetails(cred *common.AppCredential, appStatus string) (*CredentialDetails, *common.ErrorResponse) {

	refs, err := a.DbMan.GetApiProductNames(cred.Id, TypeConsumerKey)
	if err != nil {
		log.Errorf("Error when getting product reference list")
		return nil, newDbError(err)
	}
	return &CredentialDetails{
		ApiProductReferences: refs,
		AppID:                cred.AppId,
		AppStatus:            appStatus,
		Attributes:           a.DbMan.GetKmsAttributes(cred.TenantId, cred.Id)[cred.Id],
		ConsumerKey:          cred.Id,
		ConsumerSecret:       cred.ConsumerSecret,
		ExpiresAt:            cred.ExpiresAt,
		IssuedAt:             cred.IssuedAt,
		MethodType:           cred.MethodType,
		Scopes:               common.JsonToStringArray(cred.Scopes),
		Status:               cred.Status,
	}, nil
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

func writeJson(code int, obj interface{}, w http.ResponseWriter, r *http.Request) {

	requestId := r.Header.Get(headerRequestId)
	bytes, err := json.Marshal(obj)
	// JSON error
	if err != nil {
		code = http.StatusInternalServerError
		log.Errorf("unable to marshal errorResponse for request_id=[%s]: %v", requestId, err)
		jsonError := &common.ErrorResponse{
			ResponseCode:    strconv.Itoa(JSON_MARSHAL_ERROR),
			ResponseMessage: fmt.Sprintf("JSON Marshal Error %v for object: %v", err, obj),
			StatusCode:      http.StatusInternalServerError,
		}
		if bytes, err = json.Marshal(jsonError); err != nil {
			log.Errorf("unable to marshal JSON error response for request_id=[%s]: %v", requestId, err)
			w.Header().Set("Content-Type", "text/plain")
			bytes = []byte("unable to marshal errorResponse : " + err.Error())
		} else {
			w.Header().Set("Content-Type", "application/json")
		}
	} else { // success
		w.Header().Set("Content-Type", "application/json")
	}
	w.WriteHeader(code)
	log.Debugf("Sending response_code=%d for request_id=[%s]: %s", code, requestId, bytes)
	w.Write(bytes)
}
