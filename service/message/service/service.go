package service

import (
	"babygrow/service/message/model"
	"context"
)

func Create(entity *model.GMessage) (string, error) {
	return entity.Insert(context.Background())
}

func List(query *model.Query) ([]*model.ListModel, int64, error) {
	return model.List(query)
}

func UnreadCount(query *model.UnreadCountQuery) (int64, error) {
	return model.GetUnreadCount(query)
}
