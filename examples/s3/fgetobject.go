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
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/encrypt"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-bucketname, my-objectname
	// and my-filename.csv are dummy values, please replace them with original values.

	// Requests are always secure (HTTPS) by default. Set secure=false to enable insecure (HTTP) access.
	// This boolean value is the last argument for New().

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	accessKey := "minio"
	if a, ok := os.LookupEnv("ACCESS_KEY"); ok {
		accessKey = a
	}
	secretKey := "minio123"
	if s, ok := os.LookupEnv("SECRET_KEY"); ok {
		secretKey = s
	}
	s3Client, err := minio.New("localhost:9000", accessKey, secretKey, true)

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	s3Client.TraceOn(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	//password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.
	bucketname := "tbucket11"
	objectName := "small.csv"

	//m := map[string]string{"X-Amz-Server-Side-Encryption": "AES256"}
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+objectName))
	//encryption := encrypt.DefaultPBKDF([]byte("password"), []byte("salt"))
	//encryption := encrypt.NewSSE()
	opts := minio.GetObjectOptions{}
	opts.SetRange(0, 20)
	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(bucketname, objectName, "/home/kris/Downloads/acusss.txt", opts); err != nil {

		// if err := s3Client.FGetObject(bucketname, objectName, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully saved my-filename.csv")
}
