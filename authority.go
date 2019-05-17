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

package cce

import "crypto/x509"

// AuthorityService manages digital certificates.
type AuthorityService interface {
	// CAChain returns the certificate authority chain, starting with the
	// issuing CA and ending with the root CA (inclusive).
	CAChain() ([]*x509.Certificate, error)
	// SignCSR signs a ASN.1 DER encoded certificate signing request.
	SignCSR(der []byte, template *x509.Certificate) (*x509.Certificate, error)
}
