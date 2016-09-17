package apidVerifyApiKey

import "encoding/json"

/*
 * Ensure the ENV matches.
 */
func validateEnv(envLocal, envInPath string) bool {

	var ePaths []string
	json.Unmarshal([]byte(envLocal), &ePaths)
	for _, a := range ePaths {
		if a == envInPath {
			return true
		}
	}
	return false
}
