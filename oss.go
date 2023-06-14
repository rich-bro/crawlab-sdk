package sdk

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/ngaut/log"
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
		log.Error("oss参数获取失败")
	}

	ossClient, err := oss.New("oss-us-west-1.aliyuncs.com", ossAccessKey, ossAccessSecret)
	if err != nil {
		log.Error(err)
		return err
	}

	OssBucket, err = ossClient.Bucket(ossBucketName)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
