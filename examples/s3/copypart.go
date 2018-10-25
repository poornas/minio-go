package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const bucket string = "kannappan102"
const object1 string = "part1"
const object2 string = "part2"
const targetobject string = "object4"

func setupClient() *s3.S3 {

	accessKey := "AKIAJAS63GECEYO4Y77A"                     //"minio"
	secretKey := "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" //"minio123"
	sdkEndpoint := "https://localhost:9000"                 //"https://s3.amazonaws.com"               //"https://localhost:9000"                 //               //

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	newSession := session.New()
	s3Config := &aws.Config{
		Credentials:      creds,
		Endpoint:         aws.String(sdkEndpoint),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(false),
	}

	// Create an S3 service object in the default region.
	return s3.New(newSession, s3Config)
}

func main() {

	s3Client := setupClient()

	if s3Client == nil {
		fmt.Println("s3Client is nil, connection failed")
		return
	}

	// Create MultiPart
	resp, err := createmultipartupload(s3Client)
	if err != nil {
		fmt.Println("CreateMultipartUpload", err)
	}
	fmt.Println("uploadID", *resp.UploadId)
	up1, err := s3Client.UploadPartCopy(&s3.UploadPartCopyInput{
		PartNumber: aws.Int64(1),
		UploadId:   resp.UploadId,
		Bucket:     aws.String(bucket),
		CopySource: aws.String(path.Join(bucket, object1)),
		Key:        aws.String(targetobject),
		//CopySourceSSECustomerAlgorithm: aws.String("AES256"),
		//CopySourceSSECustomerKey:       aws.String("01234567890123456789012345678901"),
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String("01234567890123456789012345678901"),
	})
	if err != nil {
		fmt.Println("copyobjectPart1", err)
		os.Exit(1)
	}
	up2, err := s3Client.UploadPartCopy(&s3.UploadPartCopyInput{
		PartNumber: aws.Int64(2),
		UploadId:   resp.UploadId,
		Bucket:     aws.String(bucket),
		CopySource: aws.String(path.Join(bucket, object2)),
		Key:        aws.String(targetobject),
		//CopySourceSSECustomerAlgorithm: aws.String("AES256"),
		//CopySourceSSECustomerKey:       aws.String("01234567890123456789012345678901"),
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String("01234567890123456789012345678901"),
	})
	if err != nil {
		fmt.Println("copyobjectPart2", err)
		os.Exit(1)
	}
	var completedParts []*s3.CompletedPart
	completedParts = append(completedParts, &s3.CompletedPart{
		ETag:       up1.CopyPartResult.ETag,
		PartNumber: aws.Int64(int64(1))})
	completedParts = append(completedParts, &s3.CompletedPart{
		ETag:       up2.CopyPartResult.ETag,
		PartNumber: aws.Int64(int64(2))})
	lo, err := s3Client.ListMultipartUploads(&s3.ListMultipartUploadsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(targetobject),
		//UploadIdMarker: resp.UploadId,
	})

	if err != nil {
		fmt.Println("listmultiparterror", err)
	}
	fmt.Println("listmultipartop", lo.String())

	// List incomplete parts for the current uploadid
	lp, err := s3Client.ListParts(&s3.ListPartsInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(targetobject),
		UploadId: resp.UploadId,
	})

	if err != nil {
		fmt.Println("listParts", err)
	}
	fmt.Println("listParts", lp.String())

	completeResponse, err := completeMultipartUpload(s3Client, resp, completedParts)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Successfully uploaded file: %s\n", completeResponse.String())

}
func createmultipartupload(s3Client *s3.S3) (*s3.CreateMultipartUploadOutput, error) {
	//	h := md5.New()
	//	io.WriteString(h, "01234567890123456789012345678901")
	//	md5Sum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return s3Client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(targetobject),
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String("01234567890123456789012345678901"),
		//		SSECustomerKeyMD5:    aws.String(md5Sum),
	})
}
func completeMultipartUpload(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, completedParts []*s3.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	return svc.CompleteMultipartUpload(completeInput)
}

func uploadPart(svc *s3.S3, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) (*s3.CompletedPart, error) {
	tryNum := 1
	h := md5.New()
	io.WriteString(h, "01234567890123456789012345678901")
	md5Sum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	partInput := &s3.UploadPartInput{
		Body:                 bytes.NewReader(fileBytes),
		Bucket:               resp.Bucket,
		Key:                  resp.Key,
		PartNumber:           aws.Int64(int64(partNumber)),
		UploadId:             resp.UploadId,
		ContentLength:        aws.Int64(int64(len(fileBytes))),
		SSECustomerAlgorithm: aws.String("AES256"),
		SSECustomerKey:       aws.String("01234567890123456789012345678901"),
		SSECustomerKeyMD5:    aws.String(md5Sum),
	}
	maxRetries := 10
	for tryNum <= maxRetries {
		uploadResult, err := svc.UploadPart(partInput)
		if err != nil {
			if tryNum == maxRetries {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			return &s3.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int64(int64(partNumber)),
			}, nil
		}
	}
	return nil, nil
}

func abortMultipartUpload(svc *s3.S3, uploadId *string) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *uploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(targetobject),
		UploadId: uploadId,
	}
	_, err := svc.AbortMultipartUpload(abortInput)
	return err
}
