package apidVerifyApiKey

import "strings"

/*
 * Ensure the ENV matches.
 */
func validateEnv(envLocal string, envInPath string) bool {

	s := strings.TrimPrefix(envLocal, "{")
	s = strings.TrimSuffix(s, "}")
	fs := strings.Split(s, ",")
	for _, a := range fs {
		if a == envInPath {
			return true
		}
	}
	return false
}
