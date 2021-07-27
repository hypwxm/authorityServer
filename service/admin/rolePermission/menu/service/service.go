package service

import (
	"github.com/hypwxm/authorityServer/DB/appGorm"
	menuService "github.com/hypwxm/authorityServer/service/admin/menu/service"
	"github.com/hypwxm/authorityServer/service/admin/rolePermission/menu/dao"
	"github.com/hypwxm/authorityServer/service/admin/rolePermission/menu/dbModel"
	userService "github.com/hypwxm/authorityServer/service/admin/user/service"

	"fmt"

	"github.com/hypwxm/authorityServer/util/interfaces"
)

func Save(query interfaces.QueryInterface) ([]*dbModel.GRoleMenu, error) {
	var err error
	var roleId string
	if roleId = query.GetStringValue("roleId"); roleId == "" {
		return nil, fmt.Errorf("操作错误")
	}
	db := appGorm.Open()

	list := make([]*dbModel.GRoleMenu, 0)

	menuIds := query.ToStringArray("menuIds")
	for _, v := range menuIds {
		var et = &dbModel.GRoleMenu{
			MenuId: v,
			RoleId: roleId,
		}
		list = append(list, et)
	}

	list, err = dao.MultiInsert(db, list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func Del(query interfaces.QueryInterface) error {
	db := appGorm.Open()
	return dao.Delete(db, query)
}

func List(query interfaces.QueryInterface) (interfaces.ModelMapSlice, error) {
	if roleIds := query.ToStringArray("roleIds"); len(roleIds) == 0 {
		// 如果roleid为空，去userId对应的权限，
		//给别人分配权限的时候只能以自己拥有的权限为基准
		if query.GetStringValue("userId") == "" {
			return nil, fmt.Errorf("操作错误")
		}
		user, err := userService.Get(interfaces.QueryMap{"id": query.GetStringValue("userId")})
		if err != nil {
			return nil, err
		}

		if user.GetStringValue("account") == "admin" {
			// 究极管理员无需判断，最高权限
			ms, _, err := menuService.List(interfaces.QueryMap{})
			if err != nil {
				return nil, err
			}
			return ms, nil
		} else {
			roles := user.GetValue("roles").(interfaces.ModelMapSlice)
			roleIds := make([]string, 0)
			for _, v := range roles {
				roleIds = append(roleIds, v.GetStringValue("id"))
			}
			if len(roleIds) == 0 {
				return nil, fmt.Errorf("操作错误")
			}
			query.Set("roleIds", roleIds)
		}
	}

	db := appGorm.Open()
	return dao.List(db, query)
}
