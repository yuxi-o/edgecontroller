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

package main_test

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/open-ness/common/log"
	"github.com/open-ness/edgecontroller/telemetry"
)

var _ = Describe("Telemetry", func() {
	Describe("Syslog", func() {
		var logger *log.Logger

		BeforeEach(func() {
			logger = &log.Logger{}

			// Dial to the Controller syslog server over TLS
			clientConf := conf.Clone()
			clientConf.ServerName = telemetry.SyslogSNI
			Expect(
				logger.ConnectSyslogTLS("localhost:6514", clientConf),
			).To(Succeed())
		})

		It("Sends one syslog line to the Controller syslog server", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logger.Alert("BAD THINGS")

			contents, err := tailF(ctx, filepath.Join(telemDir, "syslog.log"))
			Expect(err).NotTo(HaveOccurred())
			Eventually(contents).Should(Say("BAD THINGS\n"))
		})

		It("Sends multiple syslog lines to the Controller syslog server", func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logger.Alert("THING ONE")
			logger.Alert("THING TWO")

			contents, err := tailF(ctx, filepath.Join(telemDir, "syslog.log"))
			Expect(err).NotTo(HaveOccurred())
			Eventually(contents).Should(Say("THING ONE\n[^\n]+THING TWO\n"))
		})
	})

	Describe("StatsD", func() { /* TODO */ })
})

func tailF(ctx context.Context, name string) (*Buffer, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return BufferReader(ignoreEOFReader{f, ctx}), nil
}

type ignoreEOFReader struct {
	io.Reader
	ctx context.Context
}

func (r ignoreEOFReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
RetryRead:
	for err == io.EOF {
		select {
		case <-r.ctx.Done():
			break RetryRead
		case <-time.After(50 * time.Millisecond):
			n, err = r.Reader.Read(p)
		}
	}
	return n, err
}
