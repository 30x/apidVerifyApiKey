package apidVerifyApiKey

// Fields related to consumer key
type ClientIdDetails struct {
	ClientId string `json:"clientId,omitempty"`

	ClientSecret string `json:"clientSecret,omitempty"`

	RedirectURIs []string `json:"redirectURIs,omitempty"`

	Status string `json:"status,omitempty"`

	// Attributes associated with the client Id.
	Attributes []Attribute `json:"attributes,omitempty"`
}
