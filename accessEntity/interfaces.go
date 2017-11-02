package accessEntity

import (
	"github.com/apid/apidVerifyApiKey/common"
	"net/http"
)

type ApiManagerInterface interface {
	common.ApiManagerInterface
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

type DbManagerInterface interface {
	common.DbManagerInterface
	GetApiProducts(org, priKey, priVal, secKey, secVal string) (apiProducts []common.ApiProduct, err error)
	GetApps(org, priKey, priVal, secKey, secVal string) (apps []common.App, err error)
	GetCompanies(org, priKey, priVal, secKey, secVal string) (companies []common.Company, err error)
	GetCompanyDevelopers(org, priKey, priVal, secKey, secVal string) (companyDevelopers []common.CompanyDeveloper, err error)
	GetAppCredentials(org, priKey, priVal, secKey, secVal string) (appCredentials []common.AppCredential, err error)
	GetDevelopers(org, priKey, priVal, secKey, secVal string) (developers []common.Developer, err error)
	GetApiProductNamesByAppId(appId string) ([]string, error)
	GetAppNamesByComId(comId string) ([]string, error)
	GetComNamesByDevId(devId string) ([]string, error)
	GetAppNamesByDevId(devId string) ([]string, error)
	GetComNameByComId(comId string) (string, error)
	GetDevEmailByDevId(devId string) (string, error)
	GetStatus(id, t string) (string, error)
	GetOrgName(tenantId string) (string, error)
}
