package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/ngaut/log"
	"github.com/rich-bro/crawlab-sdk/entity"
	"github.com/rich-bro/crawlab-sdk/interfaces"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
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

}

func containsAll(source map[string]int, target []string) bool {
	for _, targetStr := range target {
		//log.Debug(targetStr)
		//log.Debug(source[targetStr])
		if source[targetStr] == 0 {
			return false
		}
	}
	return true
}

func switchTable(items []entity.Result) {

	fileds := map[string]int{}
	for _, item := range items {
		for k, _ := range item {
			fileds[k] = 1
		}
		break
	}
	//log.Debug(fileds)
	//log.Debug(containsAll(fileds, []string{"title", "content", "source", "files", "keywords", "author_names", "authors", "timezone"}))
	if containsAll(fileds, []string{"title", "content", "source", "files", "keywords", "author_names", "authors", "timezone"}) {
		//报告
		log.Debug("验证：报告")
		thinktankVerifyKeys = map[string]interface{}{
			"id":                []string{"empty"},
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
			"authors":           []string{"json", "fields:author_id,author_name,author_url"},
			"timezone":          []string{"empty", `regex:[\+|-]\d{4}`},
			"timezone_location": []string{"empty"},
		}

	} else if containsAll(fileds, []string{"title", "name", "area_of_expertise", "location", "phone", "email", "education", "website"}) {
		//专家
		log.Debug("验证：专家")
		thinktankVerifyKeys = map[string]interface{}{
			"id":             []string{"empty"},
			"title":          []string{"json"},
			"name":           []string{"empty"},
			"site_name_cn":   []string{"empty"},
			"site_name":      []string{"empty"},
			"source":         []string{"empty"},
			"audios":         []string{"json"},
			"videos":         []string{"json"},
			"related_topics": []string{"json"},
			"files":          []string{"json"},
			"oss_files":      []string{"json"},
			"domain":         []string{"empty"},
			"created_at":     []string{"empty", "int", "length:13"},
			"updated_at":     []string{"empty", "int", "length:13"},
		}
	} else {
		log.Debug("未匹配，不验证")
	}

}

func verify(items []entity.Result) error {
	switchTable(items)
	for _, item := range items {
		for k, v := range item {
			if thinktankVerifyKeys[k] != nil {
				vfuncs := thinktankVerifyKeys[k].([]string)
				for _, vfunc := range vfuncs {
					log.Debug(v)

					switch v.(type) {
					case []string:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be array!", k))
					case []int:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be array!", k))
					case []int64:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be array!", k))
					case []interface{}:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be array!", k))
					case []map[string]interface{}:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be array!", k))
					case map[interface{}]interface{}:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be map!", k))
					case map[string]interface{}:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be map!", k))
					case map[string]string:
						return errors.New(fmt.Sprintf("ERROR: %s cannot be map!", k))
					}

					if len(strings.Split(vfunc, ":")) > 1 {
						switch strings.Split(vfunc, ":")[0] {
						case "fields":
							for _, field := range strings.Split(strings.Split(vfunc, ":")[1], ",") {
								errList := []error{}
								gjson.Parse(v.(string)).ForEach(func(key, value gjson.Result) bool {
									if len(value.String()) != 0 {
										if !value.Get(field).Exists() {
											errList = append(errList, errors.New(fmt.Sprintf("ERROR: %s:%s not Exist!", k, field)))
										}
									}
									return true
								})
								if len(errList) > 0 {
									return errList[0]
								}
							}
						case "length":
							///							log.Debug(vfunc)
							lenCount, _ := strconv.Atoi(strings.Split(vfunc, ":")[1])

							switch v.(type) {
							case int:
								if len(strconv.Itoa(v.(int))) != lenCount {
									return errors.New(fmt.Sprintf("ERROR: %s length must be %d", k, lenCount))
								}
							case int64:
								if len(strconv.FormatInt(v.(int64), 10)) != lenCount {
									return errors.New(fmt.Sprintf("ERROR: %s length must be %d", k, lenCount))
								}
							}
						case "regex":
							regexStr := strings.Split(vfunc, ":")[1]
							rs := regexp.MustCompile(regexStr)
							strArr := rs.FindAllString(v.(string), -1)
							if len(strArr) != 1 {
								return errors.New(fmt.Sprintf("ERROR: %s regex %s match error", k, regexStr))
							}

							if strArr[0] != v.(string) {
								return errors.New(fmt.Sprintf("ERROR: %s regex %s match error", k, regexStr))
							}
						}
					} else {
						switch vfunc {
						case "empty":
							switch v.(type) {
							case string:
								if len(v.(string)) == 0 {
									return errors.New(fmt.Sprintf("ERROR: %s cannot be empty!", k))
								}
							case int:
								if v.(int) == 0 {
									return errors.New(fmt.Sprintf("ERROR: %s cannot be empty!", k))
								}
							case int64:
								if v.(int64) == 0 {
									return errors.New(fmt.Sprintf("ERROR: %s cannot be empty!", k))
								}
							}

						case "json":
							if len(v.(string)) != 0 {
								//log.Debug(v)
								var js json.RawMessage
								err := json.Unmarshal([]byte(v.(string)), &js)
								//log.Debug(js)
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
		log.Error(err)
		return err
	} else {
		log.Debug("verify true")
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

func SaveItem(items ...entity.Result) error {
	return GetResultService().SaveItem(items...)
}

func SaveItems(items []entity.Result) {
	GetResultService().SaveItems(items)
}
