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

package pki

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// StoreCertificate persists a certificate chain to disk.
func StoreCertificate(path string, certs ...*x509.Certificate) error {
	if len(certs) == 0 {
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "unable to create file")
	}
	defer file.Close()

	for _, cert := range certs {
		if err = pem.Encode(
			file,
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: cert.Raw,
			},
		); err != nil {
			return errors.Wrapf(err, "unable to store certificate at %s", path)
		}
	}

	if err = file.Chmod(0600); err != nil {
		return errors.Wrap(err, "unable to set certificate file permissions")
	}

	return nil
}

// LoadCertificate loads a certificate from disk.
func LoadCertificate(path string) (*x509.Certificate, error) {
	var (
		err   error
		bytes []byte
		block *pem.Block
	)

	if bytes, err = ioutil.ReadFile(filepath.Clean(path)); err != nil {
		return nil, errors.Wrap(err, "unable to read certificate file")
	}

	if block, _ = pem.Decode(bytes); block == nil {
		return nil, errors.New("unable to decode certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}
