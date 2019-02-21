package main_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Acceptance: Service", func() {
	It("Should say hello on HTTP GET request", func() {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
		Expect(err).ToNot(HaveOccurred(), "Problem sending request")
		defer resp.Body.Close()

		Expect(resp.StatusCode).Should(Equal(http.StatusOK))

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred(), "Problem reading response body")
		Expect(string(bodyBytes)).Should(ContainSubstring("hello"))
	})
})
