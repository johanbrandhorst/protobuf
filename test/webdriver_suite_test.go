package webdriver_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/sclevine/agouti"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	testproto "github.com/johanbrandhorst/protobuf/test/server/proto/test"
	"github.com/johanbrandhorst/protobuf/test/shared"
)

func TestWebdriver(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Webdriver Suite")
}

var (
	testServer   *gexec.Session
	chromeDriver = agouti.ChromeDriver(
		agouti.Desired(agouti.Capabilities{
			"loggingPrefs": map[string]string{
				"browser": "INFO",
			},
			"browserName": "chrome",
		}),
		agouti.ChromeOptions(
			"args", []string{
				"--headless",
				"--disable-gpu",
				"--allow-insecure-localhost",
			},
		),
		agouti.ChromeOptions(
			// Requires Chrome 62
			"binary", "/usr/bin/google-chrome-unstable",
		),
	)
	seleniumDriver = agouti.Selenium(
		agouti.Browser("firefox"),
		agouti.Desired(agouti.NewCapabilities("acceptInsecureCerts")),
		/* Headless firefox does not yet have a way to accept unknown certificates
		agouti.Desired(agouti.Capabilities{
			"moz:firefoxOptions": map[string][]string{
				"args": []string{"-headless"},
			},
		}),
		*/
	)
	client testproto.TestServiceClient
)

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

	By("Starting the WebDrivers", func() {
		if os.Getenv("GOPHERJS_SERVER_ADDR") == "" {
			Expect(chromeDriver.Start()).NotTo(HaveOccurred())
			//Expect(seleniumDriver.Start()).NotTo(HaveOccurred())
		}
	})

	By("Dialing the gRPC Server", func() {
		grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard))
		tc, err := credentials.NewClientTLSFromFile("./insecure/localhost.crt", "")
		Expect(err).NotTo(HaveOccurred())
		cc, err := grpc.Dial("localhost"+shared.HTTP2Server,
			grpc.WithBlock(),
			grpc.WithTransportCredentials(tc),
			grpc.WithStreamInterceptor(func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				fmt.Fprintln(GinkgoWriter, "Calling", method)
				defer func() { fmt.Fprintln(GinkgoWriter, "Finished", method) }()
				return streamer(ctx, desc, cc, method, opts...)
			}),
			grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				fmt.Fprintln(GinkgoWriter, "Calling", method)
				defer func() { fmt.Fprintln(GinkgoWriter, "Finished", method) }()
				return invoker(ctx, method, req, reply, cc, opts...)
			}),
		)
		Expect(err).NotTo(HaveOccurred())
		client = testproto.NewTestServiceClient(cc)
	})
})

var _ = AfterSuite(func() {
	By("Stopping the WebDrivers", func() {
		if os.Getenv("GOPHERJS_SERVER_ADDR") == "" {
			Expect(chromeDriver.Stop()).NotTo(HaveOccurred())
			//Expect(seleniumDriver.Stop()).NotTo(HaveOccurred())
		}
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
