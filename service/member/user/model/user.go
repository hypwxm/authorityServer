package model

import (
	"babygrowing/DB/pgsql"
	"babygrowing/util/database"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hypwxm/rider/utils/cryptos"
	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type GMember struct {
	database.BaseColumns

	Nickname  string `json:"nickname" db:"nickname"`
	RealName  string `json:"realName" db:"realname"`
	FirstName string `json:"firstName" db:"firstname"`
	LastName  string `json:"lastName" db:"lastname"`

	Avatar   string `json:"avatar" db:"avatar"`
	Phone    string `json:"phone" db:"phone"`
	Account  string `json:"account" db:"account"`
	Password string `json:"password" db:"password"`
	Salt     string `json:"-" db:"salt"`
}

/*
 * Insert
 * 创建会员，会员注册
 */
func (gm *GMember) Insert() (string, error) {
	var err error

	// 必须有登录账号
	if strings.TrimSpace(gm.Account) == "" {
		return "", fmt.Errorf("新用户账号不能为空")
	}
	// 必须有登录密码
	if strings.TrimSpace(gm.Password) == "" {
		return "", fmt.Errorf("新用户密码不能为空")
	}

	gm.BaseColumns.Init()

	// 为新用户创建唯一盐
	gm.Salt = cryptos.RandString()
	gm.Password = SignPwd(gm.Password, gm.Salt)

	db := pgsql.Open()
	// 插入判断用户登录账号是否已经存在
	stmt, err := db.PrepareNamed("insert into g_member (createtime, isdelete, disabled, nickname, realname, firstname, lastname, account, password, salt, id) select :createtime, :isdelete, :disabled, :nickname, :realname, :firstname, :lastname, :account, :password, :salt, :id where not exists(select 1 from g_member where account = :account and isdelete='false') returning id")

	if err != nil {
		return "", err
	}
	log.Println(stmt.QueryString)
	var lastID string
	err = stmt.Get(&lastID, gm)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("该账号已被注册")
		}
		return "", err
	}
	return gm.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

func (self *GMember) GetByID(db *sqlx.DB, query *GetQuery) (*GMember, error) {
	stmt, err := db.PrepareNamed("select * from g_member where id=:id")
	if err != nil {
		return nil, err
	}
	var user = new(GMember)
	err = stmt.Get(user, query)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type Query struct {
	pgsql.BaseQuery
}

func (self *GMember) List(query *Query) ([]*GMember, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	var selectSql = `SELECT * FROM g_member WHERE 1=1 `
	var whereSql = ""
	whereSql = pgsql.BaseWhere(query.BaseQuery)

	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}
	optionSql := pgsql.BaseOption(query.BaseQuery)
	stmt, err := db.PrepareNamed(selectSql + whereSql + optionSql)
	if err != nil {
		return nil, 0, err
	}
	log.Println(stmt.QueryString)

	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users = make([]*GMember, 0)
	for rows.Next() {
		var user = new(GMember)
		err = rows.StructScan(&user)
		if err != nil {
			return nil, 0, err
		}
		user.Password = ""
		user.Salt = ""
		users = append(users, user)
	}

	return users, count, nil

}

func (self *GMember) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from g_member where 1=1 " + strings.Join(whereSql, " "))
	if err != nil {
		return 0, err
	}
	var count int64
	err = stmt.Get(&count, query)
	log.Println(stmt.QueryString, query)
	return count, err
}

type UpdateByIDQuery struct {
	ID        string `db:"id"`
	Nickname  string `db:"nickname"`
	RealName  string `db:"realname"`
	FirstName string `db:"firstname"`
	LastName  string `db:"lastname"`
	Phone     string `db:"phone"`
	Avatar    string `db:"avatar"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GMember) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	var updateSql = ""
	// updateSql = updateSql + " ,nickname=:nickname"
	updateSql = updateSql + " ,realname=:realname"
	updateSql = updateSql + " ,firstname=:firstname"
	updateSql = updateSql + " ,lastname=:lastname"
	updateSql = updateSql + " ,phone=:phone"

	db := pgsql.Open()

	stmt, err := db.PrepareNamed("update g_member set id=:id  " + updateSql + " where id=:id and isdelete=false and disabled=false")
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (self *GMember) UpdateNickname(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	var updateSql = ""
	// updateSql = updateSql + " ,nickname=:nickname"
	updateSql = updateSql + " ,nickname=:nickname"

	db := pgsql.Open()

	stmt, err := db.PrepareNamed("update g_member set id=:id  " + updateSql + " where id=:id and isdelete=false and disabled=false")
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (self *GMember) UpdateAvatar(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	var updateSql = ""
	// updateSql = updateSql + " ,nickname=:nickname"
	updateSql = updateSql + " ,avatar=:avatar"

	db := pgsql.Open()

	stmt, err := db.PrepareNamed("update g_member set id=:id  " + updateSql + " where id=:id and isdelete=false and disabled=false")
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
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
func (self *GMember) Delete(db *sqlx.DB, query *DeleteQuery) error {
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

	stmt, err := db.PrepareNamed("update g_member set isdelete=true where id=any(:ids)")
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
func (self *GMember) ToggleDisabled(query *DisabledQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("操作条件错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("update g_member set disabled=:disabled where id=:id and isdelete=false")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}

// 根据条件获取单个用户
func (self *GMember) Get(query *GMember) (*GMember, error) {
	db := pgsql.Open()

	var selectSql = `
		select * from g_member 
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
	var user = new(GMember)
	err = stmt.QueryRow(query).StructScan(user)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(query.Password) != "" {
		// 如果密码传过来了，是登录事件
		signedPwd := SignPwd(query.Password, user.Salt)
		if signedPwd != user.Password {
			return nil, errors.New("密码错误")
		}
	}
	return user, nil
}
