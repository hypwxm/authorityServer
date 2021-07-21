package model

import (
	"authorityServer/DB/pgsql"
	roleModel "authorityServer/service/admin/role/model"
	mediaModel "authorityServer/service/media/model"
	mediaService "authorityServer/service/media/service"

	"authorityServer/util"
	"authorityServer/util/database"

	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hypwxm/rider/utils/cryptos"
	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const BusinessName = "g_admin_user"

type GAdminUserRole struct {
	UserId string `json:"userId" db:"user_id"`
	OrgId  string `json:"orgId" db:"org_id"`
	RoleId string `json:"roleId" db:"role_id"`
}

type GAdminUserOrg struct {
	UserId string `json:"userId" db:"user_id"`
	OrgId  string `json:"orgId" db:"org_id"`
}

type UserAndRoleModel struct {
	roleModel.GAdminRole
	UserId string `json:"userId" db:"user_id"`
	RoleId string `json:"roleId" db:"role_id"`
}

type GAdminUser struct {
	database.BaseColumns

	CreatorId string `json:"creatorId" db:"creator_id"`
	Creator   string `json:"creator" db:"creator"`

	Account    string              `json:"account" db:"account"`
	Password   string              `json:"password" db:"password"`
	Username   string              `json:"username" db:"username"`
	ContactWay string              `json:"contactWay" db:"contact_way"`
	Post       string              `json:"post" db:"post"`
	Salt       string              `json:"salt" db:"salt"`
	Media      []*mediaModel.Media `json:"media" db:"-"`

	Sort int `json:"sort" db:"sort"`

	Roles []*UserAndRoleModel `json:"roles" db:"roles"`
}

func insertRoles(orgs []*UserAndRoleModel, tx *sqlx.Tx) error {
	for _, v := range orgs {
		stmt, err := tx.PrepareNamed(roleInsertSql())
		log.Println(stmt.QueryString)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(v)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
根据用户ids拿到对应角色列表
*/

func GetRolesByUserIds(ids []string) ([]*UserAndRoleModel, error) {
	db := pgsql.Open()
	whereSQL := roleListSql()
	stmt, err := db.PrepareNamed(whereSQL)
	if err != nil {
		return nil, err
	}
	log.Println(stmt.QueryString)

	arr := pq.StringArray{}
	for _, v := range ids {
		arr = append(arr, v)
	}

	rows, err := stmt.Queryx(map[string]interface{}{
		"user_ids": arr,
	})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list = make([]*UserAndRoleModel, 0)
	for rows.Next() {
		var item = new(UserAndRoleModel)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

func (self *GAdminUser) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Account) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.Password) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.Username) == "" {
		return "", fmt.Errorf("操作错误")
	}
	// admin账号为系统初始化时需要进行创建，最高权限
	if self.Account != "admin" {
		if len(self.Roles) == 0 {
			return "", fmt.Errorf("操作错误")
		}
	}

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	// 插入判断用户登录账号是否已经存在
	stmt, err := tx.PrepareNamed(insertSql())
	if err != nil {
		return "", err
	}
	log.Println(stmt.QueryString)
	var lastId string
	self.BaseColumns.Init()

	self.Salt = cryptos.RandString()
	self.Password = util.SignPwd(self.Password, self.Salt)

	err = stmt.Get(&lastId, self)
	if err != nil {
		return "", err
	}

	for _, v := range self.Roles {
		v.UserId = lastId
	}

	err = insertRoles(self.Roles, tx)
	if err != nil {
		return "", err
	}

	// 先把媒体文件插入数据库
	medias := mediaService.InitMedias(self.Media, BusinessName, lastId, self.CreatorId)
	err = mediaService.MultiCreate(medias)
	if err != nil {
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return self.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

type GetModel struct {
	GAdminUser
	Avatar string `json:"avatar"`
}

func (self *GAdminUser) GetByID(query *GetQuery) (*GetModel, error) {
	db := pgsql.Open()
	stmt, err := db.PrepareNamed(getByIdSql())
	if err != nil {
		return nil, err
	}
	var entity = new(GetModel)
	err = stmt.Get(entity, query)
	if err != nil {
		return nil, err
	}

	// 查找对应的媒体信息
	medias, _, err := mediaService.List(&mediaModel.Query{
		BusinessIds: []string{entity.ID},
		Businesses:  []string{BusinessName},
	})

	if err != nil {
		return nil, err
	}

	// 查询角色信息
	roles, err := GetRolesByUserIds([]string{entity.ID})
	if err != nil {
		return nil, err
	}

	entity.Media = medias
	if len(medias) > 0 {
		entity.Avatar = medias[0].Url
	}

	entity.Roles = roles

	return entity, nil
}

type Query struct {
	pgsql.BaseQuery
	Keywords string `db:"keywords"`
	Status   int    `db:"status"`
	OrgId    string `db:"org_id"`

	// 查询多个组织，当这个参数的长度大于0时，orgId将会失效
	OrgIds pq.StringArray `db:"org_ids"`

	RoleIds pq.StringArray `db:"role_ids"`
}

type ListModel struct {
	GAdminUser
	OrgId  string `json:"orgId" db:"org_id"`
	Avatar string `json:"avatar"`
}

func (self *GAdminUser) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	whereSql, fullSql := listSql(query)
	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}
	stmt, err := db.PrepareNamed(fullSql)
	if err != nil {
		return nil, 0, err
	}
	log.Println(stmt.QueryString)

	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list = make([]*ListModel, 0)
	var ids []string = make([]string, 0)
	for rows.Next() {
		var item = new(ListModel)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
		ids = append(ids, item.ID)
	}

	// 查找对应的媒体信息
	medias, _, err := mediaService.List(&mediaModel.Query{
		BusinessIds: ids,
		Businesses:  []string{BusinessName},
	})

	if err != nil {
		return nil, 0, err
	}

	// 查询角色信息
	roles, err := GetRolesByUserIds(ids)
	if err != nil {
		return nil, 0, err
	}

	for _, v := range list {
		for _, vm := range medias {
			if v.ID == vm.BusinessId {
				v.Media = append(v.Media, vm)
				v.Avatar = vm.Url
			}
		}
		for _, vm := range roles {
			if v.ID == vm.UserId {
				v.Roles = append(v.Roles, vm)
			}
		}
	}

	return list, count, nil

}

func (self *GAdminUser) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	sqlStr := countSql(whereSql...)
	if strings.TrimSpace(query.OrgId) != "" {
		sqlStr = countSqlByOrgId(whereSql...)
	}
	stmt, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return 0, err
	}
	var count int64
	err = stmt.Get(&count, query)
	log.Println(stmt.QueryString, query)
	return count, err
}

type UpdateByIDQuery struct {
	UserId string `json:"userId"`

	ID         string `json:"id" db:"id"`
	Password   string `json:"password" db:"password"`
	Username   string `json:"username" db:"username"`
	ContactWay string `json:"contactWay" db:"contact_way"`
	Post       string `json:"post" db:"post"`
	Sort       int    `json:"sort" db:"sort"`

	Roles []*UserAndRoleModel `json:"roles" db:"-"`
	Media []*mediaModel.Media `json:"media" db:"-"`

	Disabled bool `json:"disabled" db:"disabled"`

	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GAdminUser) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareNamed(updateSql(query))
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	query.Updatetime = util.GetCurrentMS()

	// 如果password有更新的话
	if strings.TrimSpace(query.Password) != "" {
		if !util.ValidatePwd(query.Password) {
			return fmt.Errorf("密码太短")
		}
		user, err := self.Get(&GAdminUser{BaseColumns: database.BaseColumns{ID: query.ID}})
		if err != nil {
			return err
		}
		query.Password = util.SignPwd(query.Password, user.Salt)
	}
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}

	// 更新操作直接把之前的角色信息删除，再重新插入
	stmt, err = tx.PrepareNamed(roleDelSql())
	log.Println(stmt.QueryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(&GAdminUserRole{
		UserId: query.ID,
	})
	if err != nil {
		return err
	}

	err = insertRoles(query.Roles, tx)
	if err != nil {
		return err
	}
	// 先把媒体文件插入数据库
	err = mediaService.Del(&mediaModel.DeleteQuery{
		Businesses:  []string{BusinessName},
		BusinessIds: []string{query.ID},
	}, tx)
	if err != nil {
		return err
	}
	medias := mediaService.InitMedias(query.Media, BusinessName, query.ID, query.UserId)
	err = mediaService.MultiCreate(medias, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type DeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
}

// 删除，批量删除
func (self *GAdminUser) Delete(query *DeleteQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if len(query.IDs) == 0 {
		return errors.New("操作条件错误")
	}
	for _, v := range query.IDs {
		if strings.TrimSpace(v) == "" {
			return errors.New("操作条件错误")
		}
	}

	db := pgsql.Open()
	stmt, err := db.PrepareNamed(delSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}

type DisabledQuery struct {
	Disabled bool   `db:"disabled"`
	ID       string `db:"id"`
}

// 启用禁用店铺
func (self *GAdminUser) ToggleDisabled(query *DisabledQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("操作条件错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed(toggleSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}

// 根据条件获取单个用户
func (self *GAdminUser) Get(query *GAdminUser) (*GAdminUser, error) {
	db := pgsql.Open()

	var selectSql = `
		select * from g_admin_user 
	`
	var whereSql = `
		where 1=1 
	`

	if strings.TrimSpace(query.ID) != "" {
		whereSql = whereSql + " and id=:id"
	}
	if strings.TrimSpace(query.Account) != "" {
		whereSql = whereSql + " and account=:account"
	}

	stmt, err := db.PrepareNamed(selectSql + whereSql)
	if err != nil {
		return nil, err
	}
	var user = new(GAdminUser)
	err = stmt.QueryRow(query).StructScan(user)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(query.Password) != "" {
		// 如果密码传过来了，是登录事件
		signedPwd := util.SignPwd(query.Password, user.Salt)
		if signedPwd != user.Password {
			return nil, errors.New("密码错误")
		}
	}
	return user, nil
}
