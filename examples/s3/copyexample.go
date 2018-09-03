package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	minio "github.com/minio/minio-go"
)

func main() {

	s3Client, err := minio.New("localhost:9000", os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	err = s3Client.MakeBucket("my-bucketname", "")
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		data := make([]byte, 4*1024)
		size := int64(len(data))
		rand.Read(data)
		for {
			_, err := s3Client.PutObject("tbucket11", "my-objectname", bytes.NewReader(data), size, minio.PutObjectOptions{})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("1")
		}
	}()

	go func() {
		data := make([]byte, 1024)
		size := int64(len(data))
		rand.Read(data)
		for {
			_, err := s3Client.PutObject("tbucket11", "my-objectname", bytes.NewReader(data), size, minio.PutObjectOptions{})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("2")
		}
	}()

	go func() {
		// Source object
		src := minio.NewSourceInfo("tbucket11", "my-objectname", nil)
		dst, err := minio.NewDestinationInfo("my-bucketname", "my-objectname-copy", nil, nil)
		if err != nil {
			log.Fatalln(err)
		}
		for {
			// Initiate copy object.
			err := s3Client.CopyObject(dst, src)
			if err != nil {
				// log.Fatal(err)
			}
			fmt.Printf("C")
		}

	}()

	wg.Wait()
}
