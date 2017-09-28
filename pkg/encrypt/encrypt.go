/*
 * Minio Go Library for Amazon S3 Compatible Cloud Storage (C) 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package encrypt implements a generic interface to encrypt any stream of data.
// currently this package implements two types of encryption
// - Symmetric encryption using AES.
// - Asymmetric encrytion using RSA.
package encrypt

import (
	"errors"
	"fmt"
	"io"
)

const (
	// AesCbcPkcs5 specifies the client-side-encryption algorithm AES-CBC with
	// PKCS5 padding. This algorithm is implemented for AWS compability but is
	// not recommended because of security issues.
	AesCbcPkcs5 = "AES/CBC/PKCS5"

	// DareHmacSha256 specifies the client-side-encryption algorithm DARE with
	// a HMAC-SHA256 KDF scheme. This algorithm provides tamper-proof encryption
	// and is recommended over any current AWS S3 client-side-encryption algorithm.
	DareHmacSha256 = "DARE-HAMC-SHA256"
)

const (
	cseIV        = "X-Amz-Meta-X-Amz-Iv"
	cseKey       = "X-Amz-Meta-X-Amz-Key-v2"
	cseAlgorithm = "X-Amz-Meta-X-Amz-Cek-Alg"
)

// Cipher is a generic interface for en/decrypting streams using
// S3 client/server side encryption.
type Cipher interface {
	Seal(header map[string]string, src io.Reader) (io.ReadCloser, error)

	Open(header map[string]string, src io.Reader) (io.ReadCloser, error)

	Overhead(size int64) int64
}

// NewCipher creates a new cipher using the provided algorithm and key.
func NewCipher(algorithm string, key Key) (Cipher, error) {
	switch algorithm {
	default:
		return nil, fmt.Errorf("algorithm '%s' is not supported", algorithm)
	case AesCbcPkcs5:
		return aesCbcPkcs5{key: key}, nil
	case DareHmacSha256:
		symKey, ok := key.(*SymmetricKey)
		if !ok {
			return nil, errors.New("encryption key must be a symmetric key")
		}
		if len(symKey.masterKey) != 32 {
			return nil, errors.New("encryption key must be 256 bit long")
		}
		d := dareHmacSha256{}
		copy(d[:], symKey.masterKey)
		return d, nil
	}
}
