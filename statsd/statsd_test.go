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
package statsd_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/alexcesaro/statsd.v2"
)

var _ = Describe("StatsD Service", func() {
	var (
		cli *statsd.Client
		r   *rand.Rand
	)

	BeforeEach(func() {
		var err error

		cli, err = statsd.New()
		Expect(err).ToNot(HaveOccurred())

		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	AfterEach(func() {
		cli.Close()
	})

	Describe("Listening", func() {
		Describe("Success", func() {
			It("Should write received stats to a log file on host", func() {
				By("Generating a random value")
				testVal := r.Int31()

				By("Sending value as a gauge stat to the StatsD service")
				cli.Gauge("test", testVal)

				By("Verifying the gauge stat was written to the log file")
				Eventually(func() error {
					stats, err := ioutil.ReadFile("stats.log")
					if err != nil {
						return err
					}

					if !strings.Contains(string(stats), fmt.Sprint(testVal)) {
						return fmt.Errorf("unable to find: %d", testVal)
					}

					return nil
				}, 5*time.Second, time.Second).Should(BeNil())
			})
		})
	})
})
