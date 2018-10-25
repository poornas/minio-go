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
	"time"

	"github.com/minio/minio-go/pkg/encrypt"

	"github.com/minio/minio-go"
)

const (
	serverEndpoint = "SERVER_ENDPOINT"
	accessKey      = "ACCESS_KEY"
	secretKey      = "SECRET_KEY"
	enableHTTPS    = "ENABLE_HTTPS"
	runSetup       = false
)

func main() {
	if runSetup {
		setup()
	}
	//GOOD
	testCopyUnencryptedSmallObjects() // Case1: SMALL: PLAIN -> PLAIN

	testCopyDestSSES3SmallObjects() // Case2: SMALL: PLAIN -> SSE-S3

	testCopyDestSSECSmallObjects()         // Case3: SMALL: PLAIN -> SSEC
	testCopySrcDestSSECSmallObjects()      // Case4: SMALL: SSEC -> SSEC
	testCopySrcSSECDestPlainSmallObjects() // Case 5: Small SSEC-> Plain

	testCopySrcAndDestSSES3SmallObjects() //bad Case6: SMALL: SSE-S3 -> SSE-S3
	testCopySrcS3DestPlainSmallObjects()  //good Case7 : SMALL: SSE-S3 -> PLAIN
	testCopySrcSSECDestS3SmallObjects()   //bad Case 8: Small SSEC-> SSE-S3
	testCopySrcS3DestSSECSmallObjects()   //bad Case 9: Small SSE-S3-> SSE-C

	// // 	//good
	testCopySrcAndDestSSES3LargeObjects()
	testCopySrcSSECDestS3LargeObjects()
	testCopySrcS3DestSSECLargeObjects()
	testCopyUnencryptedLargeObjects()
	testCopyDestSSES3LargeObjects()
	testCopyDestSSECLargeObjects()
	testCopySrcDestSSECLargeObjects() // Case4: LARGE: SSEC -> SSEC
	testCopySrcS3DestPlainLargeObjects()
	testCopySrcSSECDestPlainLargeObjects()

}
func setup() {
	// initialize logging params
	startTime := time.Now()
	testName := "small-SSEC->SSES3" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	//small object put
	object, err := os.Open("/home/kris/Downloads/smallfile")
	//object, err := os.Open("/home/kris/Downloads/dump/100-0.txt")
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
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	bucketname := "tbucket11"
	//objectName := "plaincsv.gz"
	objectName := "small-ssec"

	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+objectName))
	// sse-c
	if n, err := s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encryption}); err != nil {
		log.Fatal("upload of ", objectName, " failed... ", n, err)
	}

	objectName = "small-s3"
	encryption = encrypt.NewSSE()
	// // sse-s3
	n, err := s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encryption})
	if err != nil {
		log.Fatal("upload of ", objectName, " failed... ", n, err)

	}
	fmt.Println("upload size of small s3=============>", n)

	objectName = "small"
	// unencrypted
	if n, err = s3Client.PutObject(bucketname, objectName, object, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"}); err != nil {
		log.Fatal("upload of ", objectName, " failed... ", n, err)
	}
	fmt.Println("upload size of small=============>", n)

	objectName = "large-ssec"
	lobject, err := os.Open("/home/kris/Downloads/dump/100-0.txt")
	defer lobject.Close()
	if err != nil {
		log.Fatal("couldnt open file")
	}

	objectStat, err = lobject.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	encryption = encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+objectName))
	// sse-c
	if n, err := s3Client.PutObject(bucketname, objectName, lobject, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encryption}); err != nil {
		log.Fatal("upload of ", objectName, " failed... ", n, err)
	}

	objectName = "large-s3"
	encryption = encrypt.NewSSE()
	// sse-s3
	if n, err := s3Client.PutObject(bucketname, objectName, lobject, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream", ServerSideEncryption: encryption}); err != nil {
		log.Fatal("upload of ", objectName, " failed... ", n, err)

	}

	objectName = "large"
	// unencrypted
	if n, err := s3Client.PutObject(bucketname, objectName, lobject, objectStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"}); err != nil {
		log.Fatal("upload of ", objectName, " failed... ", n, err)
	}

}
func testCopySrcSSECDestS3SmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-SSEC->SSES3" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small-ssec"
	dstObject := srcObject + "-cpy-s3"
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	var srcencryption = encrypt.DefaultPBKDF([]byte(password), []byte(srcBucket+srcObject))

	src := minio.NewSourceInfo(srcBucket, srcObject, srcencryption)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, encrypt.NewSSE(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {
		fmt.Println("so get failed ----- ", opts, err)
		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopySrcS3DestSSECSmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-SSES3->SSEC" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small-s3"
	dstObject := srcObject + "-cpy-ssec"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	var dstencryption = encrypt.DefaultPBKDF([]byte(password), []byte(dstBucket+dstObject))
	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, dstencryption, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: dstencryption}})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{ServerSideEncryption: dstencryption}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}
func testCopySrcS3DestPlainSmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-SSES3->Plain" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small-s3"
	dstObject := srcObject + "-cpy-plain"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}
func testCopySrcAndDestSSES3SmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-SSES3->SSES3" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small-s3"
	dstObject := srcObject + "-cpy-sses3"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, encrypt.NewSSE(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopyDestSSES3SmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-plain->SSES3" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small"
	dstObject := srcObject + "-cpy-sses3"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, encrypt.NewSSE(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("****%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat success:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}
func testCopySrcDestSSECSmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-SSEC->SSEC" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small-ssec"
	dstObject := srcObject + "-cpy-ssec"

	// All following conditions are allowed and can be combined together.

	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.
	srcencryption := encrypt.DefaultPBKDF([]byte(password), []byte(srcBucket+srcObject))

	var dstencryption = encrypt.DefaultPBKDF([]byte(password), []byte(dstBucket+dstObject))

	src := minio.NewSourceInfo(srcBucket, srcObject, srcencryption)

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, dstencryption, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("abt to copty")
	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))
	/*
		stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: dstencryption}})
		if err != nil {
			fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
		}
		fmt.Println("stat1:: ", stat.Size, stat.Metadata)

		opts := minio.GetObjectOptions{ServerSideEncryption: dstencryption}
		//	opts.SetRange(0, 20)
		//	opts.ServerSideEncryption = encryption
		if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

			// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
			log.Fatalln(err)
		}
	*/
}

func testCopySrcSSECDestPlainSmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-SSEC->SSEC" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small-ssec"
	dstObject := srcObject + "-cpy-plain"

	// All following conditions are allowed and can be combined together.

	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.
	srcencryption := encrypt.DefaultPBKDF([]byte(password), []byte(srcBucket+srcObject))

	src := minio.NewSourceInfo(srcBucket, srcObject, srcencryption)

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("abt to copty")
	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}

}

func testCopyDestSSECSmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-plain->SSEC" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small"
	dstObject := srcObject + "-cpy-ssec"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(dstBucket+dstObject))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, encryption, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encryption}})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{ServerSideEncryption: encryption}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopyUnencryptedSmallObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:small-plain->plain" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "small"
	dstObject := srcObject + "-cpy-plain"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}
func testCopySrcDestSSECLargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-SSEC->SSEC" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		fmt.Println("121")

		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large-ssec"
	dstObject := srcObject + "-cpy-ssec"
	fmt.Println("11")
	// All following conditions are allowed and can be combined together.

	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.
	srcencryption := encrypt.DefaultPBKDF([]byte(password), []byte(srcBucket+srcObject))

	var dstencryption = encrypt.DefaultPBKDF([]byte(password), []byte(dstBucket+dstObject))

	src := minio.NewSourceInfo(srcBucket, srcObject, srcencryption)

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, dstencryption, nil)
	if err != nil {
		fmt.Println("err.1")
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: dstencryption}})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{ServerSideEncryption: dstencryption}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopyUnencryptedLargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-plain->plain" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large"
	dstObject := srcObject + "-cpy-plain"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopyDestSSES3LargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-plain->SSES3" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large"
	dstObject := srcObject + "-cpy-sses3"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, encrypt.NewSSE(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopyDestSSECLargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-plain->SSEC" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large"
	dstObject := srcObject + "-cpy-ssec"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	encryption := encrypt.DefaultPBKDF([]byte(password), []byte(dstBucket+dstObject))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, encryption, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encryption}})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{ServerSideEncryption: encryption}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopySrcSSECDestPlainLargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-SSEC->plain" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large-ssec"
	dstObject := srcObject + "-cpy-plain"

	// All following conditions are allowed and can be combined together.

	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.
	srcencryption := encrypt.DefaultPBKDF([]byte(password), []byte(srcBucket+srcObject))

	src := minio.NewSourceInfo(srcBucket, srcObject, srcencryption)

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("abt to copty")
	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}

}

func testCopySrcS3DestPlainLargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-SSES3->Plain" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large-s3"
	dstObject := srcObject + "-cpy-plain"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}
func testCopySrcAndDestSSES3LargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-SSES3->SSES3" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large-s3"
	dstObject := srcObject + "-cpy-sses3"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, encrypt.NewSSE(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopySrcSSECDestS3LargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-SSEC->SSES3" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large-ssec"
	dstObject := srcObject + "-cpy-s3"
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	var srcencryption = encrypt.DefaultPBKDF([]byte(password), []byte(srcBucket+srcObject))

	src := minio.NewSourceInfo(srcBucket, srcObject, srcencryption)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(dstBucket, dstObject, encrypt.NewSSE(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}

func testCopySrcS3DestSSECLargeObjects() {
	// initialize logging params
	startTime := time.Now()
	testName := "GWON:large-SSES3->SSEC" //getFuncName()
	function := "CopyObject(destination, source)"
	args := map[string]interface{}{}

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	//s3Client.TraceOn(os.Stdout)

	if err != nil {
		log.Println(testName, function, args, startTime, "", "Minio client object creation failed", err)
		return
	}

	srcBucket := "tbucket11"
	dstBucket := "tbucket11"

	srcObject := "large-s3"
	dstObject := srcObject + "-cpy-ssec"
	src := minio.NewSourceInfo(srcBucket, srcObject, nil)
	password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	var dstencryption = encrypt.DefaultPBKDF([]byte(password), []byte(dstBucket+dstObject))
	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Destination object
	dst, err := minio.NewDestinationInfo(srcBucket, dstObject, dstencryption, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("%s : Copied %s/%s to %s/%s successfully", testName, srcBucket, srcObject, dstBucket, dstObject))

	stat, err := s3Client.StatObject(dstBucket, dstObject, minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: dstencryption}})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	opts := minio.GetObjectOptions{ServerSideEncryption: dstencryption}
	//	opts.SetRange(0, 20)
	//	opts.ServerSideEncryption = encryption
	if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/ss3.txt", opts); err != nil {

		// if err := s3Client.FGetObject(dstBucket, dstObject, "/home/kris/Downloads/osses1d2.txt", opts); err != nil {
		log.Fatalln(err)
	}
}
