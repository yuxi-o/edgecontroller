package statsd_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestStatsd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StatsD Suite")
}
