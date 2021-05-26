package service2

import (
	"babygrow/DB/appGorm"
	mediaDBModel "babygrow/service/media/dbModel"
	mediaService "babygrow/service/media/service2"
	"log"

	"babygrow/service/member/user/dao"
	"babygrow/service/member/user/dbModel"
	"babygrow/util/interfaces"
	"errors"
	"fmt"
	"strings"

	"github.com/hypwxm/rider/utils/cryptos"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
)

type CreateModel struct {
	dbModel.GMember
	Medias []*mediaDBModel.Media `json:"medias" gorm:"-"`
}

func Create(entity *CreateModel) (string, error) {
	var err error

	// 必须有登录账号
	if strings.TrimSpace(entity.Account) == "" {
		return "", fmt.Errorf("新用户账号不能为空")
	}
	// 必须有登录密码
	if strings.TrimSpace(entity.Password) == "" {
		return "", fmt.Errorf("新用户密码不能为空")
	}
	db := appGorm.Open()

	query := interfaces.NewQueryMap()
	query.Set("account", entity.Account)
	query.Set("selects", "id")
	if mayExistsUser, err := dao.Get(db, query); err != nil {
		return "", err
	} else if mayExistsUser.GetID() != "" {
		return "", fmt.Errorf("账号已存在")
	}
	// 为新用户创建唯一盐
	entity.Salt = cryptos.RandString()
	entity.Password = SignPwd(entity.Password, entity.Salt)
	id, err := dao.Insert(db, &entity.GMember)

	if err != nil {
		return "", err
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

// 修改昵称
func ModifyNickname(query interfaces.QueryInterface) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	db := appGorm.Open()
	query.Set("selects", "nickname")
	return dao.Update(db, query)
}

// 只修改头像
func ModifyAvatar(query interfaces.QueryInterface) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.GetStringValue("userId")) == "" {
		return errors.New("更新条件错误")
	}
	var media = make([]*mediaDBModel.Media, 1)
	if err := mapstructure.Decode(query.GetValue("media"), &media); err != nil {
		return err
	}
	log.Printf("%+v", media)
	medias := mediaService.InitMedias(media, dbModel.BusinessName, query.GetStringValue("userId"), query.GetStringValue("userId"))
	return mediaService.MultiCreate(medias)
}

func List(query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	db := appGorm.Open()
	list, total, err := dao.List(db, query)
	if err != nil {
		return nil, 0, err
	}
	// 获取评价内容
	var ids pq.StringArray = list.GetStringValues("id")

	// 查找对应的媒体信息
	mediaService.MergeFirstMediaToListItem(interfaces.QueryMap{
		"businessIds": ids,
		"businesses":  pq.StringArray{dbModel.BusinessName},
	}, list, "id", "avatar")
	return list, total, nil
}

func Del(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}

func Get(query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	db := appGorm.Open()
	user, err := dao.Get(db, query)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(query.GetStringValue("password")) != "" {
		// 如果密码传过来了，是登录事件
		signedPwd := SignPwd(query.GetStringValue("password"), user.GetStringValue("salt"))
		if signedPwd != user.GetStringValue("password") {
			return nil, errors.New("密码错误")
		}
	}

	err = mediaService.MergeFirstMediaToItem(interfaces.QueryMap{
		"businessIds": pq.StringArray{user.GetID()},
		"businesses":  pq.StringArray{"member"},
	}, user, "id", "avatar")
	if err != nil {
		return nil, err
	}

	return user, nil
}
