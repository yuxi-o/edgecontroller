// Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package syslog_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"log/syslog"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Syslog Service", func() {
	var (
		sysLog      io.Writer
		logFile     string
		testMessage string
	)

	BeforeEach(func() {
		var err error

		sysLog, err = syslog.Dial(
			"udp",
			":514",
			syslog.LOG_INFO,
			"test")
		Expect(err).ToNot(HaveOccurred())

		logFile = filepath.Join("logs", "messages-kv.log")

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		testMessage = fmt.Sprintf("Test message %d", r.Int())
	})

	Describe("Listening", func() {
		Describe("Success", func() {
			It("Should write received messages to a log file on host", func() {
				By("Sending a syslog message to syslog service")
				_, err := fmt.Fprint(sysLog, testMessage)
				Expect(err).ToNot(HaveOccurred())

				By("Verifying the syslog message was written to the log file")
				Eventually(func() error {
					logs, err := ioutil.ReadFile(logFile)
					if err != nil {
						return err
					}

					if !strings.Contains(string(logs), testMessage) {
						return fmt.Errorf("unable to find: %s", testMessage)
					}

					return nil
				}, 5*time.Second, time.Second).Should(BeNil())
			})
		})
	})
})
