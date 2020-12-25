package model

import (
	"babygrowing/DB/pgsql"
	orgModel "babygrowing/service/admin/org/model"
	roleModel "babygrowing/service/admin/role/model"

	"babygrowing/util"
	"babygrowing/util/database"

	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hypwxm/rider/utils/cryptos"
	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type GAdminUserRole struct {
	UserId string `json:"userId" db:"user_id"`
	RoleId string `json:"roleId" db:"role_id"`
}

type GAdminUserOrg struct {
	UserId string `json:"userId" db:"user_id"`
	OrgId  string `json:"orgId" db:"org_id"`
}

type GAdminUser struct {
	database.BaseColumns

	Account    string `json:"account" db:"account"`
	Password   string `json:"password" db:"password"`
	Username   string `json:"username" db:"username"`
	ContactWay string `json:"contactWay" db:"contact_way"`
	Post       string `json:"post" db:"post"`
	Salt       string `json:"salt" db:"salt"`
	Avatar     string `json:"avatar" db:"avatar"`
	Sort       int    `json:"sort" db:"sort"`

	Orgs  []*GAdminUserOrg  `json:"orgs" db:"orgs"`
	Roles []*GAdminUserRole `json:"roles" db:"roles"`
}

func insertOrgs(orgs []*GAdminUserOrg, tx *sqlx.Tx) error {
	for _, v := range orgs {
		stmt, err := tx.PrepareNamed(orgInsertSql())
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

func insertRoles(orgs []*GAdminUserRole, tx *sqlx.Tx) error {
	for _, v := range orgs {
		stmt, err := tx.PrepareNamed(roleInsertSql())
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
根据用户ids拿到对应的组织列表
*/
type UserAndOrgModel struct {
	orgModel.GOrg
}

func GetOrgsByUserIds(ids []string) ([]*UserAndOrgModel, error) {
	db := pgsql.Open()
	whereSQL := orgListSql()
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

	var list = make([]*UserAndOrgModel, 0)
	for rows.Next() {
		var item = new(UserAndOrgModel)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}

	return list, nil
}

/**
根据用户ids拿到对应角色列表
*/
type UserAndRoleModel struct {
	roleModel.GAdminRole
}

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
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Password) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Username) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if len(self.Orgs) == 0 {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if len(self.Roles) == 0 {
		return "", errors.New(fmt.Sprintf("操作错误"))
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

	//
	err = insertOrgs(self.Orgs, tx)
	if err != nil {
		return "", err
	}

	err = insertRoles(self.Roles, tx)
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
	RoleName string `json:"roleName" db:"role_name"`
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
	return entity, nil
}

type Query struct {
	pgsql.BaseQuery
	Keywords string `db:"keywords"`
	Status   int    `db:"status"`
}

type ListModel struct {
	GAdminUser
	RoleName string `json:"roleName" db:"role_name"`
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

	var users = make([]*ListModel, 0)
	for rows.Next() {
		var user = new(ListModel)
		err = rows.StructScan(&user)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, count, nil

}

func (self *GAdminUser) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	sqlStr := countSql(whereSql...)
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
	ID         string `db:"id"`
	Password   string `json:"password" db:"password"`
	Username   string `json:"username" db:"username"`
	ContactWay string `json:"contactWay" db:"contact_way"`
	Post       string `json:"post" db:"post"`
	Avatar     string `json:"avatar" db:"avatar"`
	Sort       int    `json:"sort" db:"sort"`

	Orgs  []*GAdminUserOrg  `json:"orgs" db:"orgs"`
	Roles []*GAdminUserRole `json:"roles" db:"roles"`

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
		if util.ValidatePwd(query.Password) {
			return fmt.Errorf("密码太短")
		}
		user, err := self.Get(&GAdminUser{BaseColumns: database.BaseColumns{ID: self.ID}})
		if err != nil {
			return err
		}
		query.Password = util.SignPwd(query.Password, user.Salt)
	}
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}

	// 更新操作直接把之前的部门信息删除，再重新插入
	stmt, err = tx.PrepareNamed(orgDelSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(map[string]string{
		"user_id": self.ID,
	})
	if err != nil {
		return err
	}

	err = insertOrgs(self.Orgs, tx)
	if err != nil {
		return err
	}

	// 更新操作直接把之前的角色信息删除，再重新插入
	stmt, err = tx.PrepareNamed(roleDelSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(map[string]string{
		"user_id": self.ID,
	})
	if err != nil {
		return err
	}

	err = insertRoles(self.Roles, tx)
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
		select * from wb_admin_user 
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
