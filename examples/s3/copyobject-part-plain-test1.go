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
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go"
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
	s3Client, err := minio.NewCore("localhost:9000", accessKey, secretKey, true)

	if err != nil {
		log.Fatal("Error:", err)
	}

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	// TEST#1 :  Create a multipart object in custom backend format. Then copy it over without any additional parts into another using COPY OBJECT PART
	//password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.
	bucketName := "tbucket11"
	objectName := "plain2"
	// Make a buffer with 5MB of data
	buf := bytes.Repeat([]byte("abcde"), 1024*1024)

	// Save the data
	//objectName := randString(60, rand.NewSource(time.Now().UnixNano()), "")
	objInfo, err := s3Client.PutObject(bucketName, objectName, bytes.NewReader(buf), int64(len(buf)), "", "", map[string]string{
		"Content-Type": "binary/octet-stream",
	}, nil)
	if err != nil {
		log.Fatal("Error:", err, bucketName, objectName)
	}

	if objInfo.Size != int64(len(buf)) {
		log.Fatalf("Error: number of bytes does not match, want %v, got %v\n", len(buf), objInfo.Size)
	}
	//sse - s3
	srcInfo, err := s3Client.StatObject(bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", srcInfo, err, srcInfo.Size, srcInfo.Metadata)
	}
	//m := map[string]string{"X-Amz-Server-Side-Encryption": "AES256"}
	//encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+objectName))
	//encryption := encrypt.DefaultPBKDF([]byte("password"), []byte("salt"))
	//	encryption := encrypt.NewSSE()
	/*
		opts := minio.GetObjectOptions{}
		opts.SetRange(0, 20)
		//	opts.ServerSideEncryption = encryption
		if err := s3Client.FGetObject(bucketName, objectName, "/home/kris/Downloads/ss3.txt", opts); err != nil {

			// if err := s3Client.FGetObject(bucketname, objectName, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
			log.Fatalln(err)
		}
		log.Println("Successfully saved my-filename.csv")

	*/
	destBucketName := bucketName
	destObjectName := objectName + "-dest1"

	// Make a buffer with 5MB of data
	//buf := bytes.Repeat([]byte("abcde"), 1024*1024)

	uploadID, err := s3Client.NewMultipartUpload(destBucketName, destObjectName, minio.PutObjectOptions{}) //ServerSideEncryption: encrypt.NewSSE()})
	if err != nil {
		log.Fatal("NMU Error:", err, bucketName, objectName)
	}
	fmt.Println("started NMU,,,,")
	// Content of the destination object will be two copies of
	// `objectName` concatenated, followed by first byte of
	// `objectName`.
	metadata := map[string]string{} //"X-Amz-Server-Side-Encryption": "AES256"}
	// First of three parts
	fstPart, err := s3Client.CopyObjectPart(bucketName, objectName, destBucketName, destObjectName, uploadID, 1, 0, -1, metadata)
	if err != nil {
		log.Fatal("COP Error:", err, destBucketName, destObjectName)
	}
	fmt.Println("copied part,,,,")
	// Second of three parts
	sndPart, err := s3Client.CopyObjectPart(bucketName, objectName, destBucketName, destObjectName, uploadID, 2, 0, -1, nil)
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}

	// Last of three parts
	lstPart, err := s3Client.CopyObjectPart(bucketName, objectName, destBucketName, destObjectName, uploadID, 3, 0, 1, nil)
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}
	// Complete the multipart upload
	err = s3Client.CompleteMultipartUpload(destBucketName, destObjectName, uploadID, []minio.CompletePart{fstPart, sndPart, lstPart})
	if err != nil {
		log.Fatal("CMU Error:", err, destBucketName, destObjectName)
	}
	fmt.Println("completed upload,,,,")

	// Stat the object and check its length matches
	objInfo, err = s3Client.StatObject(destBucketName, destObjectName, minio.StatObjectOptions{})
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}
	fmt.Println("destination object info...", objInfo)
	if objInfo.Size != (5*1024*1024)*2+1 {
		log.Fatal("Destination object has incorrect size!")
	}

	// Now we read the data back
	getOpts := minio.GetObjectOptions{}
	getOpts.SetRange(0, srcInfo.Size)
	r, _, err := s3Client.GetObject(destBucketName, destObjectName, getOpts)
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}
	fmt.Println("Test1: read entire dest object successfully")

	// Now we read the data back
	getOpts = minio.GetObjectOptions{}
	getOpts.SetRange(0, 5*1024*1024-1)
	r, _, err = s3Client.GetObject(destBucketName, destObjectName, getOpts)
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}
	getBuf := make([]byte, 5*1024*1024)
	_, err = io.ReadFull(r, getBuf)
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}
	if !bytes.Equal(getBuf, buf) {
		log.Fatal("Got unexpected data in first 5MB")
	}

	getOpts.SetRange(5*1024*1024, 0)
	r, _, err = s3Client.GetObject(destBucketName, destObjectName, getOpts)
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}
	getBuf = make([]byte, 5*1024*1024+1)
	_, err = io.ReadFull(r, getBuf)
	if err != nil {
		log.Fatal("Error:", err, destBucketName, destObjectName)
	}
	if !bytes.Equal(getBuf[:5*1024*1024], buf) {
		log.Fatal("Got unexpected data in second 5MB")
	}
	if getBuf[5*1024*1024] != buf[0] {
		log.Fatal("Got unexpected data in last byte of copied object!")
	}
}
