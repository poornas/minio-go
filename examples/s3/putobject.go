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
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-testfile, my-bucketname and
	// my-objectname are dummy values, please replace them with original values.

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
	s3Client, err := minio.New("localhost:9000", accessKey, secretKey, false)
	//s3Client, err := minio.New("localhost:9001", accessKey, secretKey, false)

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	s3Client.TraceOn(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
	//small object put
	//object, err := os.Open("/home/kris/Downloads/smallfile")
	object, err := os.Open("/home/kris/Downloads/dump/large12M.txt")
	//object, err := os.Open("/home/kris/Downloads/dump/largefile")
	//object, err := os.Open("/home/kris/code/src/github.com/minio/mygoodcsv.csv.gz")
	//object, err := os.Open("/home/kris/Downloads/test.csv")

	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()
	objectStat, err := object.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	//password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	bucketname := "tbucket11"
	//objectName := "plaincsv.gz"
	objectName := "original"

	// //m := map[string]string{"X-Amz-Server-Side-Encryption": "AES256"}
	//encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+objectName))
	// // sse-c
	//n, err := s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.NewSSE)})
	// sse-s3
	n, err := s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/gzip", ServerSideEncryption: encrypt.NewSSE()})

	//n, err := s3Client.PutObject("tt1b", "s3enc-s1mall"b, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.NewSSE()})

	//n, err := s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Uploaded", "my-objectname", " of size: ", n, "Successfully.")

	/* large object put */
	//	object, err := os.Open("/home/kris/Downloads/smallfile")
	/*object, err := os.Open("/home/kris/Downloads/wso2is-5.6.0.zip")

	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()
	objectStat, err := object.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	*/
	//m := map[string]string{"X-Amz-Server-Side-Encryption": "AES256"}
	// n, err := s3Client.PutObject("test", "lsses3", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.NewSSE()})
	//n, err := s3Client.PutObject("test", "lssec", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.DefaultPBKDF([]byte("password"), []byte("salt"))})
	//n, err := s3Client.PutObject("test", "lplain", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})

	/* fake a sse-c and sse-s3 header at the same time
	m := map[string]string{"X-Amz-Server-Side-Encryption": "AES256"}
		n, err := s3Client.PutObject("test", "sses3", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.NewSSE()})
	n, err := s3Client.PutObject("test", "ssecandsses3", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.DefaultPBKDF([]byte("password"), []byte("salt")), UserMetadata: m})

	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Uploaded", "my-objectname", " of size: ", n, "Successfully.")
	*/
}
