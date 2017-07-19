package apidVerifyApiKey

type VerifyApiKeyRequest struct {
	Action string `json:"action"`

	Key string `json:"key"`

	UriPath string `json:"uriPath"`

	OrganizationName string `json:"organizationName"`

	EnvironmentName string `json:"environmentName"`

	ApiProxyName string `json:"apiProxyName"`

	// when this flag is false, authentication of key and authorization for uripath is done and authorization for apiproxies and environments is skipped. Default is true.
	ValidateAgainstApiProxiesAndEnvs bool `json:"validateAgainstApiProxiesAndEnvs,omitempty"`
}
