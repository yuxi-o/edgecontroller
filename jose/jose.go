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

package jose

import (
	"crypto"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// JWSTokenIssuer issues and validates JSON web signature tokens.
type JWSTokenIssuer struct {
	Key          crypto.PrivateKey
	KeyAlgorithm string
}

// Issue issues a new JWT token signed with the authority key and valid for one
// day. The signed JWT token is returned in the RFC 7519 compact serialization
// format.
func (s *JWSTokenIssuer) Issue() (string, error) {
	signer, err := jose.NewSigner(
		jose.SigningKey{
			Key:       s.Key,
			Algorithm: jose.SignatureAlgorithm(s.KeyAlgorithm),
		},
		new(jose.SignerOptions).WithType("JWT"))
	if err != nil {
		return "", errors.Wrap(err, "unable to create token signer")
	}

	claims := jwt.Claims{
		Expiry: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 1 day
	}

	return jwt.Signed(signer).Claims(claims).CompactSerialize()
}

// Validate validates the JWT token was signed with the authority key and has
// not yet expired. The signed JWT token is expected to be in the RFC 7519
// compact serialization format.
func (s *JWSTokenIssuer) Validate(t string) error {
	token, err := jwt.ParseSigned(t)
	if err != nil {
		return errors.Wrap(err, "unable to parse token")
	}

	key, ok := s.Key.(crypto.Signer)
	if !ok {
		return errors.Wrap(err, "invalid signing key")
	}

	var claims jwt.Claims
	err = token.Claims(key.Public(), &claims)
	if err != nil {
		return errors.Wrap(err, "unable to deserialize token claims")
	}

	return claims.Validate(jwt.Expected{Time: time.Now()})
}
