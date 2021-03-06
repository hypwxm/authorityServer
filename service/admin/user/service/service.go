package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/service/admin/user/dao"
	"github.com/hypwxm/authorityServer/service/admin/user/dbModel"
	"github.com/hypwxm/authorityServer/util"
	"github.com/hypwxm/authorityServer/util/interfaces"
	"github.com/mitchellh/mapstructure"

	"github.com/hypwxm/rider/utils/cryptos"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const NeedPwdWords = dao.NeedPwdWords

// 默认创建一个超级管理员
func InitAdmin() error {
	admin := &CreateModel{
		GAdminUser: dbModel.GAdminUser{
			Account:   "admin",
			Password:  "123456",
			Username:  "管理员",
			CreatorId: "system",
			Creator:   "系统",
		},
	}
	_, err := Create(admin)
	return err
}

type CreateModel struct {
	dbModel.GAdminUser
	Roles []*dbModel.GUserRole `json:"roles"`
}

func Create(entity *CreateModel) (string, error) {
	var err error

	if strings.TrimSpace(entity.Account) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(entity.Password) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(entity.Username) == "" {
		return "", fmt.Errorf("操作错误")
	}
	// admin账号为系统初始化时需要进行创建，最高权限
	if entity.Account != "admin" {
		if len(entity.Roles) == 0 {
			return "", fmt.Errorf("操作错误")
		}
	}

	db := appGorm.Open()
	var id string
	err = db.Transaction(func(tx *gorm.DB) error {
		oQuery := interfaces.NewQueryMap()
		oQuery.Set("account", entity.Account)
		oUser, err := dao.Get(tx, oQuery)
		if err != nil {
			return err
		}
		if oUser.GetID() != "" {
			return fmt.Errorf("账号已存在")
		}

		entity.Salt = cryptos.RandString()
		entity.Password = util.SignPwd(entity.Password, entity.Salt)

		id, err = dao.Insert(tx, &entity.GAdminUser)
		if err != nil {
			return err
		}
		for _, v := range entity.Roles {
			v.UserId = id
		}
		if entity.Account != "admin" {
			err = dao.RolesInsert(tx, entity.Roles)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return id, err
}

// 根据条件获取单个用户
func Get(query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	db := appGorm.Open()
	user, err := dao.Get(db, query)
	if err != nil {
		return nil, err
	}
	if password := query.GetStringValue("password"); password != "" {
		// 如果密码传过来了，是登录事件
		signedPwd := util.SignPwd(password, user.GetStringValue("salt"))
		if signedPwd != user.GetStringValue("password") {
			return nil, errors.New("密码错误")
		}
	} else {
		// 如果非登录状态则获取下用户的角色信息
		if roles, err := dao.GetRolesByUserIds(db, pq.StringArray{user.GetID()}); err != nil {
			return nil, err
		} else {
			user.Set("roles", roles)
		}
	}
	return user, nil
}

func List(query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	db := appGorm.Open()
	list, total, err := dao.List(db, query)

	if err != nil {
		return nil, 0, err
	}

	// 查询角色信息
	roles, err := dao.GetRolesByUserIds(db, list.GetStringValues("id"))
	if err != nil {
		return nil, 0, err
	}

	for _, v := range list {
		role := interfaces.NewModelMapSlice(0)
		for _, vm := range roles {
			if v.GetID() == vm.GetStringValue("userId") {
				role = append(role, vm)
			}
		}
		v.Set("roles", role)

	}

	return list, total, nil

}

func Modify(query interfaces.QueryInterface) error {
	db := appGorm.Open()

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := dao.Update(tx, query); err != nil {
			return err
		}

		// 如果password有更新的话
		if pwd := query.GetStringValue("password"); pwd != "" {
			if !util.ValidatePwd(pwd) {
				return fmt.Errorf("密码太短")
			}
			user, err := dao.Get(tx, interfaces.QueryMap{"id": query.GetID()})
			if err != nil {
				return err
			}
			query.Set("password", util.SignPwd(pwd, user.GetStringValue("salt")))

		}

		// 更新操作直接把之前的角色信息删除，再重新插入
		err := dao.DeleteRoles(tx, query)
		if err != nil {
			return err
		}
		roles := make([]*dbModel.GUserRole, 0)
		err = mapstructure.Decode(query.GetValue("roles"), &roles)
		if err != nil {
			return err
		}
		err = dao.RolesInsert(tx, roles)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func Del(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}
