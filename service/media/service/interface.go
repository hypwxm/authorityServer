package service

import "babygrow/service/media/model"

type MapperInterface interface {
	GetId() string
	GetMedia() []*model.Media
}

func Mapper(list []MapperInterface, query *model.Query) (interface{}, int, error) {
	// 查找对应的媒体信息
	medias, count, err := List(query)

	if err != nil {
		return nil, 0, err
	}

	for _, v := range list {
		media := v.GetMedia()
		if media == nil {
			break
		}
		for _, vm := range medias {
			if v.GetId() == vm.BusinessId {
				media = append(media, vm)
			}
		}
	}
	return list, count, err
}
