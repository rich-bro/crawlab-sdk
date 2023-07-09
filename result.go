package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/rich-bro/crawlab-sdk/entity"
	"github.com/rich-bro/crawlab-sdk/interfaces"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"strings"
)

var RS *ResultService

type ResultService struct {
	// internals
	sub grpc.TaskService_SubscribeClient
}

var thinktankVerifyKeys map[string]interface{}

func init() {
	thinktankVerifyKeys = map[string]interface{}{
		"title":             []string{"empty"},
		"site_name":         []string{"empty"},
		"site_name_cn":      []string{"empty"},
		"content":           []string{"empty"},
		"source":            []string{"empty"},
		"files":             []string{"json"},
		"images":            []string{"json"},
		"videos":            []string{"json"},
		"audios":            []string{"json"},
		"links":             []string{"json"},
		"domain":            []string{"empty"},
		"keywords":          []string{"json"},
		"lang":              []string{"empty"},
		"country_cn":        []string{"empty"},
		"country_code":      []string{"empty"},
		"created_at":        []string{"empty", "int", "length:13"},
		"updated_at":        []string{"empty", "int", "length:13"},
		"created_time":      []string{"empty", "int", "length:10"},
		"oss_files":         []string{"json"},
		"oss_images":        []string{"json"},
		"topics":            []string{"json"},
		"tags":              []string{"json"},
		"authors":           []string{"json", "filed:author_id,author_name,arthor_url"},
		"timezone":          []string{"empty"},
		"timezone_location": []string{"empty"},
	}
}

func verify(items []entity.Result) error {
	for _, item := range items {
		for k, v := range item {

			if thinktankVerifyKeys[k] != nil {
				vfuncs := thinktankVerifyKeys[k].([]string)
				for _, vfunc := range vfuncs {
					if len(strings.Split(vfunc, ":")) > 1 {
						switch strings.Split(vfunc, ":")[0] {
						case "filed":
							for _, field := range strings.Split(strings.Split(vfunc, ":")[1], ",") {
								if !gjson.Parse(v.(string)).Get(field).Exists() {
									return errors.New(fmt.Sprintf("ERROR: %s:%s not Exist!", k, field))
								}
							}
						case "length":
							lenCount, _ := strconv.Atoi(strings.Split(vfunc, ":")[1])
							if len(strconv.FormatInt(v.(int64), 10)) != lenCount {
								return errors.New(fmt.Sprintf("ERROR: %s length must be %d", k, lenCount))
							}
						}
					} else {
						switch vfunc {
						case "empty":
							if len(v.(string)) == 0 {
								return errors.New(fmt.Sprintf("ERROR: %s cannot be empty!", k))
							}
						case "json":
							if len(v.(string)) != 0 {
								_, err := json.Marshal(v.(string))
								if err != nil {
									return errors.New(fmt.Sprintf("ERROR: %s json string parse fail!", k))
								}
							}

						case "int":
							switch v.(type) {
							case int64:
							case int:
							default:
								return errors.New(fmt.Sprintf("ERROR: %s field type is not int!", k))
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (svc *ResultService) SaveItem(items ...entity.Result) error {
	err := verify(items)
	if err != nil {
		return err
	}
	svc.save(items)
	return nil
}

func (svc *ResultService) SaveItems(items []entity.Result) {
	svc.save(items)
}

func (svc *ResultService) save(items []entity.Result) {
	var _items []entity.Result
	for i, item := range items {
		_items = append(_items, item)
		if i > 0 && i%50 == 0 {
			svc._save(_items)
			_items = []entity.Result{}
		}
	}
	if len(_items) > 0 {
		svc._save(_items)
	}
}

func (svc *ResultService) _save(items []entity.Result) {
	// skip if no task id specified
	if GetTaskId().IsZero() {
		return
	}

	var records []interface{}
	for _, item := range items {
		item["_tid"] = GetTaskId()
		records = append(records, item)
	}
	data, err := json.Marshal(&entity.StreamMessageTaskData{
		TaskId:  GetTaskId(),
		Records: records,
	})
	if err != nil {
		trace.PrintError(err)
		return
	}
	if err := svc.sub.Send(&grpc.StreamMessage{
		Code: grpc.StreamMessageCode_INSERT_DATA,
		Data: data,
	}); err != nil {
		trace.PrintError(err)
		return
	}
}

func (svc *ResultService) init() (err error) {
	c := GetClient()
	taskClient := c.GetTaskClient()
	svc.sub, err = taskClient.Subscribe(context.Background())
	if err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func GetResultService(opts ...ResultServiceOption) interfaces.ResultService {
	if RS != nil {
		return RS
	}

	// service
	svc := &ResultService{}

	// apply options
	for _, opt := range opts {
		opt(svc)
	}

	// initialize
	if err := svc.init(); err != nil {
		panic(err)
	}

	RS = svc

	return svc
}

func SaveFileToOss(task entity.OssTask) error {
	bucket, err := GetOssBucket()
	if err != nil {
		return err
	}

	switch task.Type {
	case 1:
		if task.FilePath == "" || task.OssPath == "" {
			return errors.New("file path or oss path is null")
		}
		err = bucket.PutObjectFromFile(task.OssPath, task.FilePath)
	case 2:
		if task.FileIOReader == nil {
			return errors.New("file is reader is null")
		}
		err = bucket.PutObject(task.OssPath, task.FileIOReader)
	default:
		err = errors.New("not match type")
	}

	if err != nil {
		return err
	}

	return nil

}

func OssVisitLink(ossPath string, expiredTs int64) (string, error) {
	bucket, err := GetOssBucket()
	if err != nil {
		return "", err
	}

	url, err := bucket.SignURL(ossPath, http.MethodGet, expiredTs)

	if err != nil {
		return "", err
	}
	return url, nil
}

func SaveItem(items ...entity.Result) {
	GetResultService().SaveItem(items...)
}

func SaveItems(items []entity.Result) {
	GetResultService().SaveItems(items)
}
