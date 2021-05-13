package service

import (
	"babygrow/DB/appGorm"
	mediaDBModel "babygrow/service/media/dbModel"
	mediaService "babygrow/service/media/service2"

	"babygrow/service/member/daily/dao"
	"babygrow/service/member/daily/dbModel"
	dailyCommentService "babygrow/service/member/dailyComment/service2"
	"babygrow/util"
	"babygrow/util/interfaces"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

type CreateModel struct {
	dbModel.GDaily
	Medias []*mediaDBModel.Media `json:"medias" gorm:"-"`
}

func Create(entity *CreateModel) (string, error) {
	if strings.TrimSpace(entity.UserId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(entity.BabyId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	db := appGorm.Open()
	id, err := dao.Insert(db, &entity.GDaily)
	if err != nil {
		return "", err
	}
	medias := mediaService.InitMedias(entity.Medias, dbModel.BusinessName, id, entity.UserId)
	if err := mediaService.MultiCreate(medias); err != nil {
		log.Println(err)
	}
	return id, nil
}

func Modify(query interfaces.QueryInterface) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.GetID()) == "" {
		return errors.New("更新条件错误")
	}
	db := appGorm.Open()
	return dao.Update(db, query)
}

func List(query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	db := appGorm.Open()
	list, total, err := dao.List(db, query)
	if err != nil {
		return nil, 0, err
	}
	// 获取评价内容
	var ids pq.StringArray = list.GetStringValues("id")
	var userIds pq.StringArray = list.GetStringValues("userId")
	userIds = util.ArrayStringDuplicateRemoval(userIds)

	// 查找对应的媒体信息
	mediaService.MergeMediaToListItem(interfaces.QueryMap{
		"businessIds": ids,
		"businesses":  pq.StringArray{dbModel.BusinessName},
	}, list, "", "id")
	err = mediaService.MergeFirstMediaToListItem(interfaces.QueryMap{
		"businessIds": userIds,
		"businesses":  pq.StringArray{"member"},
	}, list, "userId", "avatar")
	if err != nil {
		return nil, 0, err
	}

	commentQuery := interfaces.NewQueryMap()
	commentQuery.Set("diaryIds", ids)
	if comments, err := dailyCommentService.Count(commentQuery); err != nil {
		return nil, 0, err
	} else {
		for _, v := range list {
			v["commentCount"] = comments[v.GetID()]
		}
	}
	return list, total, nil
}

func Del(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}

func Get(query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	if query.GetID() == "" {
		return nil, fmt.Errorf("参数错误")
	}
	db := appGorm.Open()
	diary, err := dao.Get(db, query)
	if err != nil {
		return nil, err
	}

	// 查找对应的媒体信息
	mediaService.MergeMediaToItem(interfaces.QueryMap{
		"businessIds": pq.StringArray{diary.GetID()},
		"businesses":  pq.StringArray{dbModel.BusinessName},
	}, diary, "", "id")
	err = mediaService.MergeFirstMediaToItem(interfaces.QueryMap{
		"businessIds": diary.GetValue("userId"),
		"businesses":  pq.StringArray{"member"},
	}, diary, "userId", "avatar")
	if err != nil {
		return nil, err
	}

	commentQuery := interfaces.NewQueryMap()
	commentQuery.Set("diaryId", diary.GetID())
	count, err := dailyCommentService.Count(commentQuery)
	if err != nil {
		return nil, err
	}
	diary.Set("commentCount", count[diary.GetID()])
	return diary, nil
}
