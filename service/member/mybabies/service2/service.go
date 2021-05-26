package service

import (
	"babygrow/DB/appGorm"
	mediaDBModel "babygrow/service/media/dbModel"
	mediaService "babygrow/service/media/service2"
	"babygrow/service/member/mybabies/dao"
	"babygrow/service/member/mybabies/daoApply"
	"babygrow/service/member/mybabies/daomb"
	"babygrow/service/member/mybabies/dbModel"
	userService "babygrow/service/member/user/service2"
	"babygrow/util"
	"babygrow/util/interfaces"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type CreateModel struct {
	dbModel.GMyBabies
	Medias []*mediaDBModel.Media `json:"medias" gorm:"-"`
}

func Create(entity *CreateModel) (string, error) {
	if strings.TrimSpace(entity.Name) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(entity.Birthday) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(entity.Gender) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(entity.UserID) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	db := appGorm.Open()
	var id string
	var err error
	db.Transaction(func(tx *gorm.DB) error {
		id, err = dao.Insert(tx, &entity.GMyBabies)
		if err != nil {
			return err
		}
		memberBaby := new(dbModel.GMemberBabyRelation)
		memberBaby.RoleName = entity.RoleName
		memberBaby.BabyId = entity.ID
		memberBaby.UserId = entity.UserID
		_, err = daomb.Insert(tx, memberBaby)
		if err != nil {
			return err
		}

		// 返回 nil 提交事务
		return nil
	})

	medias := mediaService.InitMedias(entity.Medias, dbModel.BusinessName, id, entity.UserID)
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
	entity, err := dao.Get(db, query)
	if err != nil {
		return nil, err
	}

	// 查找对应的媒体信息
	err = mediaService.MergeFirstMediaToItem(interfaces.QueryMap{
		"businessIds": pq.StringArray{entity.GetStringValue("id")},
		"businesses":  pq.StringArray{"member"},
	}, entity, "id", "avatar")
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func GetBabyRelations(query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	if query.GetStringValue("babyId") == "" {
		return nil, 0, nil
	}
	db := appGorm.Open()
	return daomb.List(db, query)
}

// 前段传过来的是账号，所以这里需要用账号去查询用户id
func CreateBabyRelations(query interfaces.QueryInterface) (string, error) {
	// 用账号名去查询用户信息，拿到用户id
	user, err := userService.Get(query)
	if err != nil {
		return "", err
	}
	query.Set("userId", user.GetID())
	db := appGorm.Open()

	entity := new(dbModel.GMemberBabyRelation)
	mapstructure.Decode(query, entity)
	return daomb.Insert(db, entity)
}

func ApplyJoinFamily(query interfaces.QueryInterface) (string, error) {
	// 先判断下是否已经是成员了
	q := interfaces.NewQueryMap()
	q.Set("userId", query.GetValue("userId"))
	q.Set("babyId", query.GetValue("babyId"))
	if _, c, err := GetBabyRelations(q); err != nil {
		return "", err
	} else if c > 0 {
		return "", fmt.Errorf("已经是成员了")
	}

	// 创建申请记录，需要发起人进行审核
	db := appGorm.Open()

	entity := new(dbModel.GMemberBabyRelationApply)
	mapstructure.Decode(query, entity)
	return daoApply.Insert(db, entity)
}

// 删除关系
func DelRelations(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return daomb.Delete(db, query)
}

// 获取申请记录
func GetApplyMsg(query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	db := appGorm.Open()
	return daomb.List(db, query)
}
