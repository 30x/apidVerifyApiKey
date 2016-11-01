package apidVerifyApiKey

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validate Env", func() {

	It("validation1", func() {
		s := validateEnv("{foo,bar}", "foo")
		Expect(s).Should(BeTrue())
	})
	It("validation2", func() {
		s := validateEnv("{foo,bar}", "bar")
		Expect(s).Should(BeTrue())
	})
	It("validation3", func() {
		s := validateEnv("{foo,bar}", "xxx")
		Expect(s).Should(BeFalse())
	})
	It("validation4", func() {
		s := validateEnv("{}", "xxx")
		Expect(s).Should(BeFalse())
	})
})
