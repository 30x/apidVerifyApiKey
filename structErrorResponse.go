package apidVerifyApiKey

// Error response returned
type ErrorResponse struct {
	ResponseCode string `json:"response_code,omitempty"`

	ResponseMessage string `json:"response_message,omitempty"`

	Kind string `json:"kind,omitempty"`
}
