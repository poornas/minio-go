package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	minio "github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/encrypt"
)

func main() {
	// Note: YOUR-ACCESSKEYID, YOUR-SECRETACCESSKEY, my-testfile, my-bucketname and
	// my-objectname are dummy values, please replace them with original values.

	// New returns an Amazon S3 compatible client object. API compatibility (v2 or v4) is automatically
	// determined based on the Endpoint value.
	accessKey := "AKIAJAS63GECEYO4Y77A"
	secretKey := "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"
	bucketName := "tbucket11"
	//endPoint := "s3.amazonaws.com"
	//bucketName := "test"
	endPoint := "localhost:9000"
	objectName := "sssec-dbug"
	path := "/home/kris/Downloads/dump/large12M.txt"

	s3Client, err := minio.New(endPoint, accessKey, secretKey, true)
	if err != nil {
		log.Fatalln(err)
	}

	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	s3Client.SetCustomTransport(tr)
	s3Client.TraceOn(os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}

	key := []byte("01234567890123456789012345678901")

	encryption, err := encrypt.NewSSEC(key)
	//	encryption := encrypt.NewSSE()
	var size int64
	size, err = s3Client.FPutObject(bucketName, objectName, path, minio.PutObjectOptions{
		ServerSideEncryption: encryption,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("size", size)

	reader, err := s3Client.GetObject(bucketName, objectName, minio.GetObjectOptions{ServerSideEncryption: encryption})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	decBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalln(err)
	}
	f, err := os.Open(path)
	origBytes, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	if !bytes.Equal(decBytes, origBytes) {
		log.Fatalln("error in matching")
	}
}
