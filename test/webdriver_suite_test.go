package webdriver_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var testServer *gexec.Session
var agoutiDriver *agouti.WebDriver

var _ = BeforeSuite(func() {
	var binPath string
	By("Building the server", func() {
		var err error
		binPath, err = gexec.Build("./server/main.go")
		Expect(err).NotTo(HaveOccurred())
	})

	By("Running the server", func() {
		var err error
		testServer, err = gexec.Start(exec.Command(binPath), GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	By("Starting the WebDriver", func() {
		agoutiDriver = agouti.ChromeDriver(
		// Unfortunately headless doesn't seem to work quite yet,
		// seems lock up loading the page.
		// (tried Google Chrome 59.0.3071.115)
		// https://developers.google.com/web/updates/2017/04/headless-chrome#drivers
		//agouti.ChromeOptions("args", []string{
		//	"--headless",
		//	"--disable-gpu",
		//}),
		)
		Expect(agoutiDriver.Start()).NotTo(HaveOccurred())
	})
})

var _ = AfterSuite(func() {
	By("Stopping the WebDriver", func() {
		Expect(agoutiDriver.Stop()).NotTo(HaveOccurred())
	})

	By("Stopping the server", func() {
		testServer.Terminate()
		testServer.Wait()
		Expect(testServer).To(gexec.Exit())
	})

	By("Cleaning up built artifacts", func() {
		gexec.CleanupBuildArtifacts()
	})
})
