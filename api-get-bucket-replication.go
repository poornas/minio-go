/*
 * MinIO Go Library for Amazon S3 Compatible Cloud Storage
 * Copyright 2020 MinIO, Inc.
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

package minio

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/minio/minio-go/v7/pkg/s3utils"
)

// GetBucketReplication fetches bucket replication configuration.If config is not
// found, returns empty string with nil error.
func (c Client) GetBucketReplication(ctx context.Context, bucketName string) (string, error) {
	// Input validation.
	if err := s3utils.CheckValidBucketName(bucketName); err != nil {
		return "", err
	}
	bucketReplicationCfg, err := c.getBucketReplication(ctx, bucketName)
	if err != nil {
		errResponse := ToErrorResponse(err)
		if errResponse.Code == "ReplicationConfigurationNotFoundError" {
			return "", nil
		}
		return "", err
	}
	return bucketReplicationCfg, nil
}

// Request server for current bucket replication config.
func (c Client) getBucketReplication(ctx context.Context, bucketName string) (string, error) {
	// Get resources properly escaped and lined up before
	// using them in http request.
	urlValues := make(url.Values)
	urlValues.Set("replication", "")

	// Execute GET on bucket to get replication config.
	resp, err := c.executeMethod(ctx, "GET", requestMetadata{
		bucketName:  bucketName,
		queryValues: urlValues,
	})

	defer closeResponse(resp)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", httpRespToErrorResponse(resp, bucketName, "")
	}

	bucketReplicationBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	replicationCfg := string(bucketReplicationBuf)
	return replicationCfg, err
}
