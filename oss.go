package sdk

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"os"
)

var (
	OssBucket *oss.Bucket
)

func OssClientInit() error {

	ossEndpoint := os.Getenv("CRAWLAB_OSS_ENDPOINT")
	ossAccessKey := os.Getenv("CRAWLAB_OSS_ACCESS_KEY")
	ossAccessSecret := os.Getenv("CRAWLAB_OSS_SECRET")
	ossBucketName := os.Getenv("CRAWLAB_OSS_BUCKET")

	if ossEndpoint == "" || ossAccessKey == "" || ossAccessSecret == "" {
		return errors.New("oss参数获取失败")
	}

	ossClient, err := oss.New(ossEndpoint, ossAccessKey, ossAccessSecret)
	if err != nil {
		return err
	}

	OssBucket, err = ossClient.Bucket(ossBucketName)
	if err != nil {
		return err
	}

	return nil
}
