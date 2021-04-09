package service

import (
	"babygrow/DB/appGorm"
	mediaModel "babygrow/service/media/model"
	mediaService "babygrow/service/media/service"
	"babygrow/service/member/daily/dao"
	"babygrow/service/member/daily/dbModel"
	dailyCommentModel "babygrow/service/member/dailyComment/model"
	dailyCommentService "babygrow/service/member/dailyComment/service"
	"babygrow/util"
	"babygrow/util/interfaces"
	"errors"
	"fmt"
	"log"
	"strings"
)

type CreateModel struct {
	dbModel.GDaily
	Medias []*mediaModel.Media `json:"medias" gorm:"-"`
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
	var ids []string = list.GetValues("id").([]string)
	var userIds []string = list.GetValues("userId").([]string)
	userIds = util.ArrayStringDuplicateRemoval(userIds)

	// 查找对应的媒体信息
	mediaService.ListMapWithMedia(ids, dbModel.BusinessName, list, "", "id")
	err = mediaService.ListMapWithMediaFirst(userIds, "member", list, "userId", "avatar")
	if err != nil {
		return nil, 0, err
	}

	if comments, _, err := dailyCommentService.List(&dailyCommentModel.Query{
		DiaryIds: ids,
	}); err != nil {
		return nil, 0, err
	} else {
		for _, v := range list {
			v["comments"] = make([]*dailyCommentModel.ListModel, 0)
			for _, vm := range comments {
				if v.GetID() == vm.DiaryId {
					v["comments"] = append(v["comments"].([]*dailyCommentModel.ListModel), vm)
				}
			}
		}
	}
	return list, total, nil
}

func Del(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}

func Get(query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	db := appGorm.Open()
	return dao.Get(db, query)
}
