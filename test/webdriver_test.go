package webdriver_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("gRPC-Web Unit Tests", func() {
	var page *agouti.Page

	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(page.Destroy()).NotTo(HaveOccurred())
	})

	It("should pass", func() {
		By("loading the test page", func() {
			Expect(page.Navigate("https://localhost:10000")).NotTo(HaveOccurred())
			Expect(page).To(HaveURL("https://localhost:10000/"))
		})

		By("finding the number of failures", func() {
			Eventually(page.FindByClass("failed")).Should(BeFound())
			Expect(page.FindByClass("failed")).To(HaveText("0"))
		})
	})
})
