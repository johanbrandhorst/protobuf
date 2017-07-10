package webdriver_test

import (
	"os/exec"
	"syscall"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var testServer *exec.Cmd
var agoutiDriver *agouti.WebDriver

var _ = BeforeSuite(func() {
	testServer = exec.Command("go", "run", "./server/main.go")
	testServer.Stderr = GinkgoWriter
	testServer.Stdout = GinkgoWriter
	// Set process group
	testServer.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	Expect(testServer.Start()).NotTo(HaveOccurred())

	agoutiDriver = agouti.ChromeDriver(
	// Unfortunately headless doesn't seem to work quite yet,
	// seems to lock up when trying to load the page.
	// (tried Google Chrome 59.0.3071.115)
	// https://developers.google.com/web/updates/2017/04/headless-chrome#drivers
	//agouti.ChromeOptions("args", []string{
	//	"--headless",
	//	"--disable-gpu",
	//	"--remote-debugging-port=9222",
	//}),
	)
	Expect(agoutiDriver.Start()).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	Expect(agoutiDriver.Stop()).NotTo(HaveOccurred())
	pgid, err := syscall.Getpgid(testServer.Process.Pid)
	Expect(err).NotTo(HaveOccurred())
	Expect(syscall.Kill(-pgid, syscall.SIGINT)).NotTo(HaveOccurred())
	testServer.Wait() // This will error, but that's expected
})
