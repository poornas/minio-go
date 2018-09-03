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
	"path"

	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/encrypt"
)

func getEncrypt(bucket, object, ssetype string) encrypt.ServerSide {
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	if ssetype == "sse-s3" {
		return encrypt.NewSSE()
	}
	if ssetype == "sse-c" {
		return encrypt.DefaultPBKDF([]byte(password), []byte(bucket+object))
	}
	return nil
}
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

	bucketname := "tbucket11"
	objectname := "t1"
	sseType := "sse-c"
	encryption := getEncrypt(bucketname, objectname, sseType)

	// test 1
	putOpts := minio.PutObjectOptions{ServerSideEncryption: encryption}
	var getOpts minio.GetObjectOptions
	if sseType == "sse-c" {
		getOpts = minio.GetObjectOptions{ServerSideEncryption: encryption}

	}
	TestPutGetStat(s3Client, bucketname, objectname, putOpts, getOpts, "/home/kris/Downloads/dump/large6M.txt", sseType, "Test multipart sse-s3", 1)
	/*
		//small object put

		//object, err := os.Open("/home/kris/Downloads/smallfile")
		//object, err := os.Open("/home/kris/Downloads/dump/100-0.txt")
		object, err := os.Open("/home/kris/Downloads/dump/large6M.txt")
		// object, err := os.Open("/home/kris/code/src/github.com/minio/mygoodcsv.csv.gz")
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

		// //m := map[string]string{"X-Amz-Server-Side-Encryption": "AES256"}
		//encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+objectName))
		// // sse-c
		//n, err := s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encryption})
		// sse-s3
		n, err := s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/gzip", ServerSideEncryption: encrypt.NewSSE()})

		//n, err := s3Client.PutObject("tt1b", "s3enc-s1mall", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.NewSSE()})

		//n, err := s3Client.PutObject("test", "sse2s3", object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encrypt.NewSSE()})

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
func TestPutGetStat(s3Client *minio.Client, bucket, objectName string, opts minio.PutObjectOptions, getOpts minio.GetObjectOptions, filename string, ssetype string, testDesc string, testNum int) {
	fmt.Println(fmt.Sprintf("----- Starting  Test: %s -- %s for %s/%s -----", testNum, testDesc, bucket, objectName))
	fmt.Println("|------- encryption type :%s ---------|", ssetype)
	object, err := os.Open(filename)

	if err != nil {
		log.Fatalln(err)
	}
	defer object.Close()
	objectStat, err := object.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	opts.UserMetadata = map[string]string{"my-custom-header": "my-custom-val"}
	n, err := s3Client.PutObject(bucket, objectName, object, objectStat.Size(), opts)
	if err != nil {
		fmt.Println(fmt.Sprintf("Upload of %s/%s failed with err %s uploaded %s bytes", bucket, objectName, err, n))
	}

	stat, err := s3Client.StatObject(bucket, objectName, minio.StatObjectOptions{GetObjectOptions: getOpts})
	//stat, err := s3Client.StatObject(bucketname, "plaind2ssecd2s", minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encrypt.DefaultPBKDF([]byte("password"), []byte("salt"))}})

	if err != nil {
		fmt.Println("stat of unencryted object::", stat, err)

	}
	log.Println("stat ::", stat)
	downloadFile := path.Join("/home/kris/Downloads", objectName, "-dl")
	if err := s3Client.FGetObject(bucket, objectName, downloadFile, getOpts); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully saved %s", downloadFile)
	fmt.Println(fmt.Sprintf("----- Ending Test: %s -- %s for %s/%s -----", testNum, testDesc, bucket, object))

}
