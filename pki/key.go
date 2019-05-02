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
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// StoreKey persists a private key on disk at the path.
func StoreKey(key crypto.PrivateKey, path string) error {
	var (
		err  error
		der  []byte
		file *os.File
	)

	if der, err = x509.MarshalPKCS8PrivateKey(key); err != nil {
		return errors.Wrap(err, "unable to marshal EC private key to DER")
	}

	if file, err = os.Create(path); err != nil {
		return errors.Wrap(err, "unable to create private key file")
	}
	defer file.Close()

	if err = file.Chmod(0600); err != nil {
		return errors.Wrap(err, "unable to set private key file permissions")
	}

	if err = pem.Encode(
		file,
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: der,
		},
	); err != nil {
		return errors.Wrap(err, "unable to store private key")
	}

	return nil
}

// LoadKey loads a private key from disk at path.
func LoadKey(path string) (crypto.PrivateKey, error) {
	var (
		err   error
		bytes []byte
		block *pem.Block
	)

	if bytes, err = ioutil.ReadFile(path); err != nil {
		return nil, errors.Wrap(err, "unable to read key file")
	}

	if block, _ = pem.Decode(bytes); block == nil {
		return nil, errors.New("unable to decode key")
	}

	return x509.ParsePKCS8PrivateKey(block.Bytes)
}
