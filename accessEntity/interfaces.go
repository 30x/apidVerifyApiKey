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
	GetApiProducts(priKey, priVal, secKey, secVal string) (apiProducts []common.ApiProduct, err error)
}
