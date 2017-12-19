package common

import (
	"github.com/apid/apid-core"
	"github.com/apid/apid-core/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"testing"
)

var testTempDirBase string

var _ = BeforeSuite(func() {
	apid.Initialize(factory.DefaultServicesFactory())
	SetApidServices(apid.AllServices(), apid.Log().ForModule("apidApiMetadata"))
	var err error
	testTempDirBase, err = ioutil.TempDir("", "verify_apikey_")
	Expect(err).Should(Succeed())
})

func TestCommon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ApiMetadata Common Suite")
}

var _ = AfterSuite(func() {
	Expect(os.RemoveAll(testTempDirBase)).Should(Succeed())
})
