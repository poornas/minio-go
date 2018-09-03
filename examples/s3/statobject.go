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
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-bucketname and my-objectname
	// are dummy values, please replace them with original values.

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
	} //s3Client, err := minio.New("localhost:9000", accessKey, secretKey, false)
	s3Client, err := minio.New("localhost:9000", accessKey, secretKey, false)
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
			bucketname := "tbucket11"
			objectName := "t1"
			coreClient := &minio.Core{Client: s3Client}
			// Iterate over all parts and aggregate the size.
			partsInfo, err := coreClient.ListObjectParts(bucketname, objectName, "d95d937a-c32a-450f-b9c9-7e8d343c9d63", 1, 10)
			if err != nil {
				fmt.Println("errrr..", err)

			}
			fmt.Println(partsInfo)
			var size int64
			for _, partInfo := range partsInfo.ObjectParts {
				fmt.Println(partInfo)
				size += partInfo.Size
			}
			fmt.Println("total size of aprts..", size)
		}
	*/
	//password := "correct horse battery staple" // Specify your password. DO NOT USE THIS ONE - USE YOUR OWN.

	bucketname := "tbucket11"
	objectName := "original"
	// //m := map[string]string{"X-Amz-Server-Side-Encryption": "AES256"}
	//encryption := encrypt.DefaultPBKDF([]byte(password), []byte(bucketname+objectName))

	// sse-c
	// should fail
	/*
		stat, err := s3Client.StatObject(bucketname, objectName, minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encryption}})
		if err != nil {
			fmt.Println("stat of sse-s3 enc object::", stat, err)

		}
		// fmt.Println("stat:: ", stat)

		//stat, err := s3Client.StatObject("fudmod", "enc1", minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encrypt.DefaultPBKDF([]byte("peeeassword"), []byte("salt"))}})
	*/
	//sse - s3
	opts := minio.GetObjectOptions{}
	//opts.SetRange(0, -1)
	stat, err := s3Client.StatObject(bucketname, objectName, minio.StatObjectOptions{GetObjectOptions: opts})
	if err != nil {
		fmt.Println("stat1 of sse-s3 enc object::", stat, err, stat.Size, stat.Metadata)
	}
	fmt.Println("stat1:: ", stat.Size, stat.Metadata)

	/*bucketname := "tbucket11"

	  //stat, err := s3Client.StatObject(bucketname, "sses3d2sssec", minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encrypt.DefaultPBKDF([]byte("correct horse battery staple"), []byte(bucketname+"sses3d2sssec"))}})
	  stat, err := s3Client.StatObject(bucketname, "t1", minio.StatObjectOptions{})
	  //stat, err := s3Client.StatObject(bucketname, "plaind2ssecd2s", minio.StatObjectOptions{minio.GetObjectOptions{ServerSideEncryption: encrypt.DefaultPBKDF([]byte("password"), []byte("salt"))}})
	*/
	// 	if err != nil {
	// 		fmt.Println("stat2 of unencryted object::", stat, err)

	// 	}
	// 	log.Println("stat2 ::", stat)

}
