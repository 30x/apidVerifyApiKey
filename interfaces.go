package apidVerifyApiKey

import (
	"github.com/apid/apid-core"
	"net/http"
)

type ApiManagerInterface interface {
	InitAPI()
	HandleRequest(w http.ResponseWriter, r *http.Request)
}

type DbManagerInterface interface {
	SetDbVersion(string)
	GetDb() apid.DB
	GetDbVersion() string
}
