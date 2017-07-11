package webdriver_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
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
			Eventually(page.FirstByClass("failed")).Should(BeFound())
			failures, err := page.FindByID("qunit-testresult").FindByClass("failed").Text()
			Expect(err).NotTo(HaveOccurred())
			if failures == "0" {
				return
			}

			// We have at least one failure - lets compile an error message
			Eventually(page.FindByID(
				"qunit-tests",
			).AllByClass(
				"fail",
			).AllByClass(
				"fail",
			)).Should(BeFound())
			messages := page.FindByID(
				"qunit-tests",
			).AllByClass(
				"fail",
			).AllByClass(
				"fail",
			)
			elements, err := messages.Elements()
			Expect(err).NotTo(HaveOccurred())
			var errMsgs []string
			for _, element := range elements {
				// Get error summary
				msg, err := element.GetElement(api.Selector{
					Using: "css selector",
					Value: ".test-message",
				})
				Expect(err).NotTo(HaveOccurred())
				errText, err := msg.GetText()
				Expect(err).NotTo(HaveOccurred())

				// Get diff
				expected, err := element.GetElements(api.Selector{
					Using: "css selector",
					Value: "del",
				})
				Expect(err).NotTo(HaveOccurred())
				var expectedText string
				if len(expected) > 0 {
					expectedText, err = expected[0].GetText()
					Expect(err).NotTo(HaveOccurred())
				}
				actual, err := element.GetElements(api.Selector{
					Using: "css selector",
					Value: "ins",
				})
				Expect(err).NotTo(HaveOccurred())
				var actualText string
				if len(actual) > 0 {
					actualText, err = actual[0].GetText()
					Expect(err).NotTo(HaveOccurred())
				}

				errMsg := errText
				if expectedText != "" && actualText != "" {
					errMsg = fmt.Sprintf(
						"%s\n\tExpected: %q\n\tActual: %q",
						errText,
						strings.TrimSuffix(expectedText, " "),
						strings.TrimSuffix(actualText, " "),
					)
				}
				errMsgs = append(errMsgs, errMsg)
			}

			// Prints each error
			Fail(strings.Join(errMsgs, "\n-----------------------------------\n"))
		})
	})
})
