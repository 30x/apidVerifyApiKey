package verifyApiKey

import (
	"github.com/apid/apid-core"
)

const (
	ApiPath = "/verifiers/apikey"
)

var (
	services apid.Services
	log      apid.LogService
)

func SetApidServices(s apid.Services, l apid.LogService) {
	services = s
	log = l
}
