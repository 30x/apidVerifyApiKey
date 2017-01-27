package apidVerifyApiKey

import (
	"github.com/apigee-labs/transicator/common"
	"strconv"
)

func convertSuffix(i int) string{
	return strconv.FormatInt(int64(i), 10)
}

func generateTestApiProduct(suffix int) common.Row{
	return common.Row{
		"id": {
			Value: "api_product_" + convertSuffix(suffix),
		},
		"api_resources": {
			Value: "{/**, /test}",
		},
		"environments": {
			Value: "{Env_0, Env_1}",
		},
		"tenant_id": {
			Value: "tenant_id_xxxx",
		},
		"_change_selector": {
			Value: "Org_0",
		},
	}
}

func generateTestDeveloper(suffix int) common.Row {
	 return common.Row{
		"id": {
			Value: "developer_id_" + convertSuffix(suffix),
		},
		"status": {
			Value: "Active",
		},
		"email": {
			Value: "test@apigee.com",
		},
		"first_name": {
			Value: "Apigee",
		},
		"last_name": {
			Value: "Google",
		},
		"tenant_id": {
			Value: "tenant_id_xxxx",
		},
		"_change_selector": {
			Value: "Org_0",
		},
	}
}

func generateTestCompany(suffix int) common.Row {
	return common.Row{
		"id": {
			Value: "company_id_" + convertSuffix(suffix),
		},
		"status": {
			Value: "Active",
		},
		"name": {
			Value: "Apigee Corporation",
		},
		"display_name": {
			Value: "Apigee",
		},
		"tenant_id": {
			Value: "tenant_id_xxxx",
		},
		"_change_selector": {
			Value: "Org_0",
		},
	}
}

func generateTestCompanyDeveloper(suffix int) common.Row {
	return common.Row{
		"developer_id": {
			Value: "developer_id_" + convertSuffix(suffix),
		},
		"tenant_id": {
			Value: "tenant_id_0",
		},
		"_change_selector": {
			Value: "test_org0",
		},
		"company_id": {
			Value: "company_id_" + convertSuffix(suffix),
		},
	}
}

func generateTestApp(suffix1, suffix2 int) common.Row {
	return common.Row{
		"id": {
			Value: "application_id_" + convertSuffix(suffix1),
		},
		"developer_id": {
			Value: "developer_id_" + convertSuffix(suffix2),
		},
		"status": {
			Value: "Approved",
		},
		"tenant_id": {
			Value: "tenant_id_xxxx",
		},
		"callback_url": {
			Value: "http://apigee.com",
		},
		"_change_selector": {
			Value: "Org_0",
		},
	}
}

func generateTestAppCreds(suffix int) common.Row {
	return common.Row{
		"id": {
			Value: "app_credential_" + convertSuffix(suffix),
		},
		"app_id": {
			Value: "application_id_" + convertSuffix(suffix),
		},
		"status": {
			Value: "Approved",
		},
		"tenant_id": {
			Value: "tenant_id_xxxx",
		},
		"callback_url": {
			Value: "http://apigee.com",
		},
		"_change_selector": {
			Value: "Org_0",
		},
	}
}

func generateTestApiProductMapper(suffix int) common.Row {
	return common.Row{
		"apiprdt_id": {
			Value: "api_product_" + convertSuffix(suffix),
		},
		"status": {
			Value: "Approved",
		},
		"app_id": {
			Value: "application_id_" + convertSuffix(suffix),
		},
		"appcred_id": {
			Value: "app_credential_" + convertSuffix(suffix),
		},
		"tenant_id": {
			Value: "tenant_id_xxxx",
		},
		"_change_selector": {
			Value: "Org_0",
		},
	}
}


