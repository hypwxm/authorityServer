package service

import (
	"authorityServer/DB/appGorm"
	menuService "authorityServer/service/admin/menu/service"
	"authorityServer/service/admin/rolePermission/menu/dao"
	"authorityServer/service/admin/rolePermission/menu/dbModel"
	"authorityServer/service/admin/rolePermission/menu/model"

	"authorityServer/util/interfaces"
	"fmt"
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

func Del(query *model.DeleteQuery) error {
	return new(model.GRoleMenu).Delete(query)
}

func List(query interfaces.QueryMap) (interfaces.ModelMapSlice, int64, error) {
	if roleIds := query.ToStringArray("roleIds"); len(roleIds) == 0 {
		// 如果roleid为空，去userId对应的权限，
		//给别人分配权限的时候只能以自己拥有的权限为基准
		if query.GetStringValue("userId") == "" {
			return nil, fmt.Errorf("操作错误")
		}
		user, err := userService.Get(&userModel.GetQuery{ID: query.UserId})
		if err != nil {
			return nil, err
		}

		if user.Account != "admin" {
			// 究极管理员无需判断，最高权限
			for _, v := range user.Roles {
				query.Set("roleIds", append(roleIds, v.RoleId))
			}
			if len(roleIds) == 0 {
				return nil, fmt.Errorf("操作错误")
			}
		} else {
			ms, _, err := menuService.List(nil)
			if err != nil {
				return nil, err
			}
			return ms, nil
		}
	}

	db := appGorm.Open()
	return dao.List(db, query)
}
