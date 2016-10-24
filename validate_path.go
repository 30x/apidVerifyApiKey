package apidVerifyApiKey

import (
	"regexp"
	"strings"
)

/*
 * Check for the base path (API_Product) match with the path
 * received in the Request, via the customized regex, where
 * "**" gets de-normalized as ".*" and "*" as everything till
 * the next "/".
 */
func validatePath(basePath, requestBase string) bool {

	s := strings.TrimPrefix(basePath, "{")
	s = strings.TrimSuffix(s, "}")
	fs := strings.Split(s, ",")
	for _, a := range fs {
		str1 := strings.Replace(a, "**", "(.*)", -1)
		str2 := strings.Replace(a, "*", "([^/]+)", -1)
		if a != str1 {
			reg, _ := regexp.Compile(str1)
			res := reg.MatchString(requestBase)
			if res == true {
				return true
			}
		} else if a != str2 {
			reg, _ := regexp.Compile(str2)
			res := reg.MatchString(requestBase)
			if res == true {
				return true
			}
		} else if requestBase == a {
			return true
		}

		/*
		 * FIXME: SINGLE_FORWARD_SLASH_PATTERN not supported yet
		 */
	}

	/* if the i/p resource is empty, no checks need to be made */
	return len(fs) == 0
}
