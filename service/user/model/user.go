package model

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hypwxm/rider/utils/cryptos"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
	"worldbar/DB/pgsql"
	"worldbar/util"

	"github.com/jmoiron/sqlx"
)

type WbUser struct {
	ID         string `json:"id" db:"id"`
	Createtime int64  `json:"createtime" db:"createtime"`
	Updatetime int64  `json:"updatetime" db:"updatetime"`
	Deletetime int64  `json:"deletetime" db:"deletetime"`
	Isdelete   bool   `json:"isdelete" db:"isdelete"`
	Disabled   bool   `json:"disabled" db:"disabled"`

	Nickname  string `json:"nickname" db:"nickname"`
	RealName  string `json:"realName" db:"realname"`
	FirstName string `json:"firstName" db:"firstname"`
	LastName  string `json:"lastName" db:"lastname"`

	Avatar string `json:"avatar" db:"avatar"`

	Account  string `json:"account" db:"account"`
	Password string `json:"password" db:"password"`
	Salt     string `json:"-" db:"salt"`

	Type sql.NullString `db:"type"`

	House string `json:"house" db:"house"`
}

func (self *WbUser) Insert() (string, error) {
	var err error
	// 插入时间
	self.Createtime = util.GetCurrentMS()

	// 必须有登录账号
	if strings.TrimSpace(self.Account) == "" {
		return "", errors.New(fmt.Sprintf("新用户账号不能为空"))
	}
	// 必须有登录密码
	if strings.TrimSpace(self.Password) == "" {
		return "", errors.New(fmt.Sprintf("新用户密码不能为空"))
	}
	// 必须有邮箱
	if self.Type.String == "1" && strings.TrimSpace(self.House) == "" {
		return "", errors.New(fmt.Sprintf("新用户房号不能为空"))
	}
	// 默认添加直接启用
	self.Disabled = false
	self.Isdelete = false

	// 为新用户创建唯一盐
	self.Salt = cryptos.RandString()
	self.Password = SignPwd(self.Password, self.Salt)

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	// 插入判断用户登录账号是否已经存在
	stmt, err := tx.PrepareNamed("insert into wb_user (createtime, isdelete, disabled, nickname, realname, firstname, lastname, account, password, salt, type, id, house) select :createtime, :isdelete, :disabled, :nickname, :realname, :firstname, :lastname, :account, :password, :salt, :type, :id, :house where not exists(select 1 from wb_user where account = :account and isdelete='false') returning id")

	if err != nil {
		return "", err
	}
	log.Println(stmt.QueryString)
	var lastId string
	self.ID = uuid.NewV4().String()
	err = stmt.Get(&lastId, self)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("该手机号已被注册")
		}
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

func (self *WbUser) GetByID(db *sqlx.DB, query *GetQuery) (*WbUser, error) {
	stmt, err := db.PrepareNamed("select * from wb_user where id=:id")
	if err != nil {
		return nil, err
	}
	var user = new(WbUser)
	err = stmt.Get(user, query)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type Query struct {
	pgsql.BaseQuery
}

func (self *WbUser) List(query *Query) ([]*WbUser, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	var selectSql = `SELECT * FROM wb_user WHERE 1=1 `
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

	var users = make([]*WbUser, 0)
	for rows.Next() {
		var user = new(WbUser)
		err = rows.StructScan(&user)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, count, nil

}

func (self *WbUser) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from wb_user where 1=1 " + strings.Join(whereSql, " "))
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
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *WbUser) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	var updateSql = ""
	updateSql = updateSql + " ,nickname=:nickname"
	updateSql = updateSql + " ,realname=:realname"
	updateSql = updateSql + " ,firstname=:firstname"
	updateSql = updateSql + " ,lastname=:lastname"
	updateSql = updateSql + " ,avatar=:avatar"

	db := pgsql.Open()

	stmt, err := db.PrepareNamed("update wb_user set id=:id  " + updateSql + " where id=:id and isdelete=false and disabled=false")
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
func (self *WbUser) Delete(db *sqlx.DB, query *DeleteQuery) error {
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

	stmt, err := db.PrepareNamed("update wb_user set isdelete=true where id=any(:ids)")
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
func (self *WbUser) ToggleDisabled(db *sqlx.DB, query *DisabledQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("操作条件错误")
	}
	stmt, err := db.PrepareNamed("update wb_user set disabled=:disabled where id=:id and isdelete=false")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}

// 根据条件获取单个用户
func (self *WbUser) Get(query *WbUser) (*WbUser, error) {
	db := pgsql.Open()

	var selectSql = `
		select * from wb_user 
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
	var user = new(WbUser)
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
