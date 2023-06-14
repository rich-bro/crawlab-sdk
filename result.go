package sdk

import (
	"context"
	"encoding/json"
	"errors"
	grpc "github.com/crawlab-team/crawlab-grpc"
	"github.com/crawlab-team/go-trace"
	"github.com/rich-bro/crawlab-sdk/entity"
	"github.com/rich-bro/crawlab-sdk/interfaces"
	"net/http"
)

var RS *ResultService

type ResultService struct {
	// internals
	sub grpc.TaskService_SubscribeClient
}

func (svc *ResultService) SaveItem(items ...entity.Result) {
	svc.save(items)
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
	err := OssClientInit()
	if err != nil {
		return err
	}

	switch task.Type {
	case 1:
		if task.FilePath == "" || task.OssPath == "" {
			return errors.New("file path or oss path is null")
		}
		err = OssBucket.PutObjectFromFile(task.OssPath, task.FilePath)
	case 2:
		if task.FileIOReader == nil {
			return errors.New("file is reader is null")
		}
		err = OssBucket.PutObject(task.OssPath, task.FileIOReader)
	default:
		err = errors.New("not match type")
	}

	if err != nil {
		return err
	}

	return nil

}

func OssVisitLink(ossPath string, expiredTs int64) (string, error) {
	url, err := OssBucket.SignURL(ossPath, http.MethodGet, expiredTs)

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
