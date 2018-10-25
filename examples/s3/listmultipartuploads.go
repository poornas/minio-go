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
)

func main() {
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
	//	s3Client.TraceOn(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}

	bucketname := "tbucket11"
	objectname := "lmssec"
	//sseType := "sse-s3"
	//encryption := getEncrypt(bucketname, objectname, sseType)

	// Create a done channel to control 'ListObjects' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List all multipart uploads from a bucket-name with a matching prefix.
	for multipartObject := range s3Client.ListIncompleteUploads(bucketname, objectname, true, doneCh) {
		if multipartObject.Err != nil {
			fmt.Println(multipartObject.Err)
			return
		}
		fmt.Println(multipartObject.UploadID, ":", multipartObject.Key, " :", multipartObject.Size)
	}
	return
}
