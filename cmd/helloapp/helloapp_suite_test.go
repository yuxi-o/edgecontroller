package main_test

import (
	"bufio"
	"os/exec"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestHelloapp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HelloApp Suite")
}

var (
	service *gexec.Session
	port    int
)

var _ = BeforeSuite(func() {
	service, port = StartService()
})

var _ = AfterSuite(func() {
	if service != nil {
		service.Kill()
	}
})

// StartService starts the service on a random port.
// Returns the session and the port.
func StartService() (session *gexec.Session, port int) {
	exe, err := gexec.Build(
		"github.com/smartedgemec/controller-ce/cmd/helloapp",
	)
	Expect(err).ToNot(HaveOccurred(), "Problem building service")

	cmd := exec.Command(exe, "-port", "0")

	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred(), "Problem starting service")

	Eventually(session.Err, 3).Should(gbytes.Say("listening on port:"),
		"Service did not start in time")

	// Scan the next word for the port
	scanner := bufio.NewScanner(session.Err)
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	Expect(scanner.Err()).ToNot(HaveOccurred(), "Couldn't scan for port")
	port, err = strconv.Atoi(scanner.Text())
	Expect(err).ToNot(HaveOccurred(), "Couldn't parse port")

	return session, port
}
