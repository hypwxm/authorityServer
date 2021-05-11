package service

import (
	"babygrow/DB/appGorm"
	mediaDBModel "babygrow/service/media/dbModel"
	mediaService "babygrow/service/media/service2"

	"babygrow/service/member/dailyComment/dao"
	"babygrow/service/member/dailyComment/dbModel"
	"babygrow/util"
	"babygrow/util/interfaces"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

type CreateModel struct {
	dbModel.GDailyComment
	Medias []*mediaDBModel.Media `json:"medias" gorm:"-"`
}

func Create(entity *CreateModel) (string, error) {
	if strings.TrimSpace(entity.UserId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(entity.BabyId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(entity.DiaryId) == "" {
		return "", fmt.Errorf("操作错误")
	}

	db := appGorm.Open()
	id, err := dao.Insert(db, &entity.GDailyComment)
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
	err = mediaService.MergeMediaToListItem(interfaces.QueryMap{
		"businessIds": ids,
		"businesses":  pq.StringArray{dbModel.BusinessName},
	}, list, "medias", "id")
	if err != nil {
		return nil, 0, err
	}
	err = mediaService.MergeFirstMediaToListItem(interfaces.QueryMap{
		"businessIds": userIds,
		"businesses":  pq.StringArray{"member"},
	}, list, "userId", "avatar")
	if err != nil {
		return nil, 0, err
	}

	commentQuery := interfaces.NewQueryMap()
	commentQuery.Set("commentIds", ids)
	if comments, err := Count(commentQuery); err != nil {
		return nil, 0, err
	} else {
		for _, v := range list {
			v["commentCount"] = comments[v.GetStringValue("id")]
		}
	}

	return list, total, nil
}

// 返回 {[dairyId1]: count, [dairyId2]: count}
func Count(query interfaces.QueryInterface) (map[string]int64, error) {
	db := appGorm.Open()
	return dao.Count(db, query)
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
	comment, err := dao.Get(db, query)
	if err != nil {
		return nil, err
	}

	// 查找对应的媒体信息
	mediaService.MergeMediaToItem(interfaces.QueryMap{
		"businessIds": pq.StringArray{comment.GetID()},
		"businesses":  pq.StringArray{dbModel.BusinessName},
	}, comment, "", "id")
	err = mediaService.MergeFirstMediaToItem(interfaces.QueryMap{
		"businessIds": comment.GetValue("userId"),
		"businesses":  pq.StringArray{"member"},
	}, comment, "userId", "avatar")
	if err != nil {
		return nil, err
	}

	return comment, nil
}
