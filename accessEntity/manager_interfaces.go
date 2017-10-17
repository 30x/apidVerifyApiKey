package accessEntity

import (
	"net/http"
	"github.com/apid/apid-core"
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