package interfaces

import "github.com/rich-bro/crawlab-sdk/entity"

type ResultService interface {
	SaveItem(item ...entity.Result) error
	SaveItems(item []entity.Result)
}
