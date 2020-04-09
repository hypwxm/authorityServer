package model

import (
	"errors"
	"fmt"
	"github.com/hypwxm/rider/utils/cryptos"
	"github.com/lib/pq"
	"log"
	"strings"
	"worldbar/DB/pgsql"
	"worldbar/util"
	"worldbar/util/database"

	"github.com/jmoiron/sqlx"
)

type WbAdminUser struct {
	database.BaseColumns

	Account  string `json:"account" db:"account"`
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
	Salt     string `json:"salt" db:"salt"`
	Avatar   string `json:"avatar" db:"avatar"`
	Type     int    `json:"type" db:"type"`
	RoleId   string `json:"roleId" db:"role_id"`
}

func (self *WbAdminUser) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Account) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Password) == "" {
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
	WbAdminUser
	Like         bool `json:"like" db:"like"`
	TotalLike    int  `json:"totalLike" db:"total_like"`
	TotalComment int  `json:"totalComment" db:"total_comment"`
}

func (self *WbAdminUser) GetByID(query *GetQuery) (*GetModel, error) {
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
	WbAdminUser
	RoleName string `json:"roleName" db:"role_name"`
}

func (self *WbAdminUser) List(query *Query) ([]*ListModel, int64, error) {
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

func (self *WbAdminUser) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	ID       string `db:"id"`
	Avatar   string `db:"avatar"`
	Username string `db:"username"`
	RoleId   string `db:"role_id"`
	Password string `db:"password"`

	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *WbAdminUser) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()
	stmt, err := db.PrepareNamed(updateSql(query))
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	query.Updatetime = util.GetCurrentMS()
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

type DeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
}

// 删除，批量删除
func (self *WbAdminUser) Delete(query *DeleteQuery) error {
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
func (self *WbAdminUser) ToggleDisabled(query *DisabledQuery) error {
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

type UpdateSortQuery struct {
	Sort1 int    `db:"sort1"`
	Sort2 int    `db:"sort2"`
	Id1   string `db:"id1"`
	Id2   string `db:"id2"`
}

// 根据两个枚举的排序
func (self *WbAdminUser) UpdateSort(query *UpdateSortQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if query.Sort1 == 0 || query.Sort2 == 0 || query.Id1 == "" || query.Id2 == "" {
		return errors.New("参数错误")
	}
	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareNamed(`update wb_news_dynamics set sort=:sort1 where id=:id2 and isdelete=false`)
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	stmt, err = tx.PrepareNamed("update wb_news_dynamics set sort=:sort2 where id=:id1 and isdelete=false")
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

type UpdateStatusQuery struct {
	Id           string `db:"id"`
	Status       int    `db:"status"`
	StatusReason string `db:"status_reason"`
	publishTime  int64  `db:"publish_time"`
}

// 更新状态
func (self *WbAdminUser) UpdateStatus(query *UpdateStatusQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if query.Id == "" && query.Status == 0 {
		return errors.New("参数错误")
	}
	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var sqlStr = "update wb_news_dynamics set status=:status, status_reason=:status_reason where id=:id and isdelete=false"
	if query.Status == 1 {
		// 记录发布时间
		query.publishTime = util.GetCurrentMS()
		sqlStr = "update wb_news_dynamics set status=:status, status_reason=:status_reason, publish_time=:publish_time where id=:id and isdelete=false"
	}
	stmt, err := tx.PrepareNamed(sqlStr)
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// 根据条件获取单个用户
func (self *WbAdminUser) Get(query *WbAdminUser) (*WbAdminUser, error) {
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
	var user = new(WbAdminUser)
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
