package model

import (
	"authorityServer/DB/pgsql"
	menuModel "authorityServer/service/admin/menu/model"
	menuService "authorityServer/service/admin/menu/service"

	userModel "authorityServer/service/admin/user/model"
	userService "authorityServer/service/admin/user/service"
	"authorityServer/util/database"

	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

type GRoleMenu struct {
	database.BaseColumns

	RoleId string `json:"roleId" db:"role_id"`
	MenuId string `json:"menuId" db:"menu_id"`
}

type SaveQuery struct {
	RoleId  string `db:"role_id"`
	MenuIds []string
}

func (self *GRoleMenu) Save(query *SaveQuery) (string, error) {
	var err error
	if strings.TrimSpace(query.RoleId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	for _, v := range query.MenuIds {
		stmt, err := tx.PrepareNamed(saveSql())
		if err != nil {
			return "", err
		}
		var _query = &GRoleMenu{
			MenuId: v,
			RoleId: query.RoleId,
		}
		_query.BaseColumns.Init()
		log.Println(stmt.QueryString, *_query)

		_, err = stmt.Exec(_query)
		if err != nil {
			return "", err
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return self.ID, nil
}

type Query struct {
	pgsql.BaseQuery
	RoleIds pq.StringArray `db:"role_ids"`

	UserId string
}

type ListModel struct {
	GRoleMenu
	ParentId string `json:"parentId" db:"parent_id"`
	Name     string `json:"name" db:"name"`
	Path     string `json:"path" db:"path"`
	Icon     string `json:"icon" db:"icon"`
}

func (self *GRoleMenu) List(query *Query) ([]*ListModel, error) {
	if query == nil {
		query = new(Query)
	}
	var list = make([]*ListModel, 0)

	if len(query.RoleIds) == 0 {
		// 如果roleid为空，去userId对应的权限，
		//给别人分配权限的时候只能以自己拥有的权限为基准
		if query.UserId == "" {
			return nil, fmt.Errorf("操作错误")
		}
		user, err := userService.Get(&userModel.GetQuery{ID: query.UserId})
		if err != nil {
			return nil, err
		}

		if user.Account != "admin" {
			// 究极管理员无需判断，最高权限
			for _, v := range user.Roles {
				query.RoleIds = append(query.RoleIds, v.RoleId)
			}
			if len(query.RoleIds) == 0 {
				return nil, fmt.Errorf("操作错误")
			}
		} else {
			ms, _, err := menuService.List(&menuModel.Query{})
			if err != nil {
				return nil, err
			}

			for _, v := range ms {
				list = append(list, &ListModel{
					Name:     v.Name,
					ParentId: v.ParentId,
					Path:     v.Path,
					GRoleMenu: GRoleMenu{
						MenuId: v.ID,
					},
				})
			}
			return list, nil
		}
	}

	db := pgsql.Open()
	fullSql := listSql(query)
	stmt, err := db.PrepareNamed(fullSql)
	if err != nil {
		return nil, err
	}
	log.Println(stmt.QueryString)
	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var data = new(ListModel)
		err = rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return list, nil
}

type DeleteQuery struct {
	RoleId  string         `db:"role_id"`
	MenuIds pq.StringArray `db:"menu_ids"`
}

func (self *GRoleMenu) Delete(query *DeleteQuery) error {
	var err error
	if strings.TrimSpace(query.RoleId) == "" {
		return fmt.Errorf("操作错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed(deleteSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
