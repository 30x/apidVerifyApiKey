package apidVerifyApiKey

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/30x/apid"
	"github.com/30x/apid/factory"
	"io/ioutil"
	"net/http/httptest"
	"net/http"
	"os"
	"strconv"
	"github.com/apigee-labs/transicator/common"
)

var (
	testTempDir string
	testServer  *httptest.Server
)

var _ = BeforeSuite(func() {
	apid.Initialize(factory.DefaultServicesFactory())

	config := apid.Config()

	var err error
	testTempDir, err = ioutil.TempDir("", "api_test")
	Expect(err).NotTo(HaveOccurred())

	config.Set("data_path", testTempDir)

	apid.InitializePlugins()

	db, err := apid.Data().DB()
	Expect(err).NotTo(HaveOccurred())
	setDB(db)
	createTables(db)
	insertTestData(db)
	testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == apiPath {
			handleRequest(w, req)
		}
	}))
})

var _ = AfterSuite(func() {
	apid.Events().Close()
	if testServer != nil {
		testServer.Close()
	}
	os.RemoveAll(testTempDir)
})

func TestVerifyAPIKey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "VerifyAPIKey Suite")
}

func insertTestData(db apid.DB) {

	txn, err := db.Begin()
	Expect(err).ShouldNot(HaveOccurred())

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "api_product_" + result,
			Type:  1,
		}
		srvItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "{/**, /test}",
			Type:  1,
		}
		srvItems["api_resources"] = scv

		scv = &common.ColumnVal{
			Value: "{Env_0, Env_1}",
			Type:  1,
		}
		srvItems["environments"] = scv

		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv
		rows = append(rows, srvItems)
		res := insertAPIproducts(rows, txn)
		Expect(res).Should(BeTrue())
	}

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "developer_id_" + result,
			Type:  1,
		}
		srvItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "test@apigee.com",
			Type:  1,
		}
		srvItems["email"] = scv

		scv = &common.ColumnVal{
			Value: "Active",
			Type:  1,
		}
		srvItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "Apigee",
			Type:  1,
		}
		srvItems["firstName"] = scv

		scv = &common.ColumnVal{
			Value: "Google",
			Type:  1,
		}
		srvItems["lastName"] = scv

		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv

		rows = append(rows, srvItems)
		res := insertDevelopers(rows, txn)
		Expect(res).Should(BeTrue())
	}

	var j, k int
	for i := 0; i < 10; i++ {
		resulti := strconv.FormatInt(int64(i), 10)
		for j = k; j < 10+k; j++ {
			var rows []common.Row

			srvItems := common.Row{}
			resultj := strconv.FormatInt(int64(j), 10)

			scv := &common.ColumnVal{
				Value: "application_id_" + resultj,
				Type:  1,
			}
			srvItems["id"] = scv

			scv = &common.ColumnVal{
				Value: "developer_id_" + resulti,
				Type:  1,
			}
			srvItems["developer_id"] = scv

			scv = &common.ColumnVal{
				Value: "approved",
				Type:  1,
			}
			srvItems["status"] = scv

			scv = &common.ColumnVal{
				Value: "http://apigee.com",
				Type:  1,
			}
			srvItems["callback_url"] = scv

			scv = &common.ColumnVal{
				Value: "Org_0",
				Type:  1,
			}
			srvItems["_apid_scope"] = scv

			scv = &common.ColumnVal{
				Value: "tenant_id_xxxx",
				Type:  1,
			}
			srvItems["tenant_id"] = scv
			rows = append(rows, srvItems)
			res := insertApplications(rows, txn)
			Expect(res).Should(BeTrue())
		}
		k = j
	}

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "app_credential_" + result,
			Type:  1,
		}
		srvItems["id"] = scv

		scv = &common.ColumnVal{
			Value: "application_id_" + result,
			Type:  1,
		}
		srvItems["app_id"] = scv

		scv = &common.ColumnVal{
			Value: "approved",
			Type:  1,
		}
		srvItems["status"] = scv

		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv
		rows = append(rows, srvItems)
		res := insertCredentials(rows, txn)
		Expect(res).Should(BeTrue())
	}

	for i := 0; i < 10; i++ {
		var rows []common.Row
		srvItems := common.Row{}
		result := strconv.FormatInt(int64(i), 10)

		scv := &common.ColumnVal{
			Value: "api_product_" + result,
			Type:  1,
		}
		srvItems["apiprdt_id"] = scv

		scv = &common.ColumnVal{
			Value: "application_id_" + result,
			Type:  1,
		}
		srvItems["app_id"] = scv

		scv = &common.ColumnVal{
			Value: "app_credential_" + result,
			Type:  1,
		}
		srvItems["appcred_id"] = scv
		scv = &common.ColumnVal{
			Value: "approved",
			Type:  1,
		}
		srvItems["status"] = scv
		scv = &common.ColumnVal{
			Value: "Org_0",
			Type:  1,
		}
		srvItems["_apid_scope"] = scv

		scv = &common.ColumnVal{
			Value: "tenant_id_xxxx",
			Type:  1,
		}
		srvItems["tenant_id"] = scv
		rows = append(rows, srvItems)
		res := insertAPIProductMappers(rows, txn)
		Expect(res).Should(BeTrue())
	}

	txn.Commit()
}
