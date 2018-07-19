// +build ignore

/*
 * Minio Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2015-2017 Minio, Inc.
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

package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/encrypt"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-bucketname and my-objectname
	// are dummy values, please replace them with original values.

	// Requests are always secure (HTTPS) by default. Set secure=false to enable insecure (HTTP) access.
	// This boolean value is the last argument for New().

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	s3Client, err := minio.New("localhost:9000", "minio", "minio123", true)
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	s3Client.TraceOn(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
	/*
		// should fail
		stat, err := s3Client.StatObject("test", "sses3encrypted-obj2", minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encrypt.NewSSE()}})
		if err != nil {
			fmt.Println("stat of sse-s3 enc object::", stat, err)

		}
		fmt.Println("stat:: ", stat)
	*/
	stat, err := s3Client.StatObject("test", "lssec2ssec", minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encrypt.DefaultPBKDF([]byte("peeeassword"), []byte("salt"))}})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err)
	}
	fmt.Println("stat1:: ", stat)
	/*
		stat, err = s3Client.StatObject("test", "plain", minio.StatObjectOptions{})
		if err != nil {
			fmt.Println("stat2 of unencryted object::", stat, err)

		}
		log.Println("stat2 ::", stat)
	*/
}
