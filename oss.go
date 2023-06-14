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

	//用户登录名称 ray@1542487417503206.onaliyun.com
	//AccessKey ID LTAI5tDEKQzFtUw485r5sTZJ
	//AccessKey Secret grJyNYJrhbwuaV372DuHl9tPG9S6ed

	ossClient, err := oss.New("oss-us-west-1.aliyuncs.com", "LTAI5tH19diwzbw55TrCXkJK", "gKCCzpqaqIfpcrO34uMFxXu0QMcFHY")
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

//oss-us-west-1
//oss-accelerate
//oss-us-west-1.aliyuncs.com
//oss-accelerate.aliyuncs.com
//client, err := oss.New("oss-us-west-1.aliyuncs.com", "LTAI5tH19diwzbw55TrCXkJK", "gKCCzpqaqIfpcrO34uMFxXu0QMcFHY",oss.Proxy("http://127.0.0.1:2087"))
