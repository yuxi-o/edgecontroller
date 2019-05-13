package syslog_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSyslogNg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Syslog Suite")
}
