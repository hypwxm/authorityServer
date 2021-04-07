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

func Modify(query interfaces.QueryMap) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.GetID()) == "" {
		return errors.New("更新条件错误")
	}
	db := appGorm.Open()
	return dao.Update(db, query)
}

func List(query interfaces.QueryMap) ([]interfaces.ModelMap, int64, error) {
	db := appGorm.Open()
	list, total, err := dao.List(db, query)
	if err != nil {
		return nil, 0, err
	}
	// 获取评价内容
	var ids []string = make([]string, 0)
	var userIds []string = make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.GetID())
		userIds = append(userIds, v["userId"].(string))
	}
	userIds = util.ArrayStringDuplicateRemoval(userIds)

	// 查找对应的媒体信息
	mediaService.ListMapWithMedia(ids, dbModel.BusinessName, list, "", "ID")
	err = mediaService.ListMapWithMedia(userIds, "member", list, "MemberMedia", "UserId")
	if err != nil {
		return nil, 0, err
	}

	for _, v := range list {
		if len(v.MemberMedia) > 0 {
			v.Avatar = v.MemberMedia[0].Url
		}
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

func Del(query interfaces.QueryMap) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}

func Get(query interfaces.QueryMap) (interfaces.ModelMap, error) {
	db := appGorm.Open()
	return dao.Get(db, query)
}
