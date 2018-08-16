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
	//s3Client, err := minio.New("localhost:9000", accessKey, secretKey, false)
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

	// Enable trace.
	// s3Client.TraceOn(os.Stderr)

	// =========> CASE 1 :::::: encrypted obj on server side to plain object on server side SSE-S3 -> Plain
	 // GOOD -SMALL
	// Source object

	src := minio.NewSourceInfo("test", "lsses3", nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set unmodified condition, copy object unmodified since 2014 April.
	// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set matching ETag condition, copy object which matches the following ETag.
	// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

	// Set matching ETag except condition, copy object which does not match the following ETag.
	// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

	// Destination object
	dst, err := minio.NewDestinationInfo("test", "lsses32plain", nil, nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
	/* =========> CASE 2 :::::: encrypted obj on server side to encrypted object on server side
	should see different encryption key; response header; metadata
	// NOTOK - using encrypted size instead og decrypted size....

	// Source object
	src := minio.NewSourceInfo("test", "lsses3", nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set unmodified condition, copy object unmodified since 2014 April.
	// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set matching ETag condition, copy object which matches the following ETag.
	// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

	// Set matching ETag except condition, copy object which does not match the following ETag.
	// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

	// Destination object
	dst, err := minio.NewDestinationInfo("test", "lsses32sses3", encrypt.NewSSE(), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
	/*
		//=========> CASE 3 :::::: plain object on server side to encrypted object on server side
		//		 should see  encryption key; response header; metadata
	*/
	// Source object
	src := minio.NewSourceInfo("test", "lplain", nil)

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	//src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set unmodified condition, copy object unmodified since 2014 April.
	// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set matching ETag condition, copy object which matches the following ETag.
	// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

	// Set matching ETag except condition, copy object which does not match the following ETag.
	// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

	// Destination object
	dst, err := minio.NewDestinationInfo("test", "lplain2sses3", encrypt.DefaultPBKDF([]byte("password"), []byte("salt")), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")

	/*
			 =========> CASE 4 :::::: plain object on server side to SSE-C
				//GOOD

		// Source object
		src := minio.NewSourceInfo("test", "lplain", nil)

		// All following conditions are allowed and can be combined together.

		// Set modified condition, copy object modified since 2014 April.
		src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set unmodified condition, copy object unmodified since 2014 April.
		// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set matching ETag condition, copy object which matches the following ETag.
		// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

		// Set matching ETag except condition, copy object which does not match the following ETag.
		// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

		dst, err := minio.NewDestinationInfo("test", "lplain2ssec", encrypt.DefaultPBKDF([]byte("password"), []byte("salt")), nil)
		if err != nil {
			log.Fatalln(err)
		}

		// Initiate copy object.
		err = s3Client.CopyObject(dst, src)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
	/*
			 =========> CASE 5 :::::: SSE-C object on server side to SSE-S3
				//GOOD

		// Source object
		src := minio.NewSourceInfo("test", "lssec", encrypt.DefaultPBKDF([]byte("password"), []byte("salt")))

		// All following conditions are allowed and can be combined together.

		// Set modified condition, copy object modified since 2014 April.
		src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set unmodified condition, copy object unmodified since 2014 April.
		// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set matching ETag condition, copy object which matches the following ETag.
		// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

		// Set matching ETag except condition, copy object which does not match the following ETag.
		// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

		dst, err := minio.NewDestinationInfo("test", "lssec2sses3", encrypt.NewSSE(), nil)
		if err != nil {
			log.Fatalln(err)
		}

		// Initiate copy object.
		err = s3Client.CopyObject(dst, src)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
	/*
		 =========> CASE 6 :::::: SSE-S3 object on server side to SSE-C


		// Source object
		src := minio.NewSourceInfo("test", "lsses3", nil)

		// All following conditions are allowed and can be combined together.

		// Set modified condition, copy object modified since 2014 April.
		src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set unmodified condition, copy object unmodified since 2014 April.
		// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set matching ETag condition, copy object which matches the following ETag.
		// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

		// Set matching ETag except condition, copy object which does not match the following ETag.
		// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

		dst, err := minio.NewDestinationInfo("test", "lsses32ssec", encrypt.DefaultPBKDF([]byte("password"), []byte("salt")), nil)
		if err != nil {
			log.Fatalln(err)
		}

		// Initiate copy object.
		err = s3Client.CopyObject(dst, src)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
	/*
			 =========> CASE 7 :::::: SSE-S3 object on server side to plain
				// GOOD

		// Source object
		src := minio.NewSourceInfo("test", "lssec", encrypt.DefaultPBKDF([]byte("password"), []byte("salt")))

		// All following conditions are allowed and can be combined together.

		// Set modified condition, copy object modified since 2014 April.
		src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set unmodified condition, copy object unmodified since 2014 April.
		// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set matching ETag condition, copy object which matches the following ETag.
		// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

		// Set matching ETag except condition, copy object which does not match the following ETag.
		// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

		dst, err := minio.NewDestinationInfo("test", "lssec2plain", nil, nil)
		if err != nil {
			log.Fatalln(err)
		}

		// Initiate copy object.
		err = s3Client.CopyObject(dst, src)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
	/*
			 =========> CASE 8 :::::: plain object on server side to plain

			//GOOD

		// Source object
		src := minio.NewSourceInfo("test", "lplain", nil)

		// All following conditions are allowed and can be combined together.

		// Set modified condition, copy object modified since 2014 April.
		src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set unmodified condition, copy object unmodified since 2014 April.
		// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

		// Set matching ETag condition, copy object which matches the following ETag.
		// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

		// Set matching ETag except condition, copy object which does not match the following ETag.
		// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

		dst, err := minio.NewDestinationInfo("test", "lplain2plain", nil, nil)
		if err != nil {
			log.Fatalln(err)
		}

		// Initiate copy object.
		err = s3Client.CopyObject(dst, src)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
	/* 	 =========> CASE 9  :::::: sse-c object on server side to sse-c
	//GOOD

	// Source object
	src := minio.NewSourceInfo("test", "lssec", encrypt.DefaultPBKDF([]byte("password"), []byte("salt")))

	// All following conditions are allowed and can be combined together.

	// Set modified condition, copy object modified since 2014 April.
	src.SetModifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set unmodified condition, copy object unmodified since 2014 April.
	// src.SetUnmodifiedSinceCond(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

	// Set matching ETag condition, copy object which matches the following ETag.
	// src.SetMatchETagCond("31624deb84149d2f8ef9c385918b653a")

	// Set matching ETag except condition, copy object which does not match the following ETag.
	// src.SetMatchETagExceptCond("31624deb84149d2f8ef9c385918b653a")

	dst, err := minio.NewDestinationInfo("test", "lssec2ssec", encrypt.DefaultPBKDF([]byte("peeeassword"), []byte("salt")), nil)
	if err != nil {
		log.Fatalln(err)
	}

	// Initiate copy object.
	err = s3Client.CopyObject(dst, src)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Copied source object /my-sourcebucketname/my-sourceobjectname to destination /my-bucketname/my-objectname Successfully.")
	*/
}
