package interfaces

import "github.com/rich-bro/crawlab-sdk/entity"

type ResultService interface {
	SaveItem(item ...entity.Result)
	SaveItems(item []entity.Result)
}
