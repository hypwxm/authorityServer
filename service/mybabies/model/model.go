package model

import (
	"babygrowing/DB/pgsql"
	"babygrowing/util"
	"babygrowing/util/database"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type GMyBabies struct {
	database.BaseColumns

	// 姓名
	Name string `json:"name" db:"name"`
	// 生日,(公历生日)
	Birthday string `json:"birthday" db:"birthday"`
	// 性别
	Gender string `json:"gender" db:"gender"`
	// 照片
	Avatar string `json:"avatar" db:"avatar"`
	// 身份证号
	IdCard string `json:"idCard" db:"id_card"`
	// 兴趣
	Hobby string `json:"hobby" db:"hobby"`
	// 特长
	GoodAt string `json:"goodAt" db:"good_at"`
	// 喜欢的食物
	FavoriteFood string `json:"favoriteFood" db:"favorite_food"`
	// 喜欢的颜色
	FavoriteColor string `json:"favoriteColor" db:"favorite_color"`
	// 志向
	Ambition string `json:"ambition" db:"ambition"`

	UserID string `json:"userId" db:"user_id"`

	Weight float64 `json:"weight" db:"weight"`
	Height float64 `json:"height" db:"height"`
}

func (self *GMyBabies) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Name) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Birthday) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Gender) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.UserID) == "" {
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
	GMyBabies
}

func (self *GMyBabies) GetByID(query *GetQuery) (*GetModel, error) {
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
}

type ListModel struct {
	GMyBabies
}

func (self *GMyBabies) List(query *Query) ([]*ListModel, int64, error) {
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

func (self *GMyBabies) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	ID string `db:"id"`
	// 姓名
	Name string `json:"name" db:"name"`
	// 生日,(公历生日)
	Birthday string `json:"birthday" db:"birthday"`
	// 性别
	Gender string `json:"gender" db:"gender"`
	// 照片
	Avatar string `json:"avatar" db:"avatar"`
	// 身份证号
	IdCard string `json:"idCard" db:"id_card"`
	// 兴趣
	Hobby string `json:"hobby" db:"hobby"`
	// 特长
	GoodAt string `json:"goodAt" db:"good_at"`
	// 喜欢的食物
	FavoriteFood string `json:"favoriteFood" db:"favorite_food"`
	// 喜欢的颜色
	FavoriteColor string `json:"favoriteColor" db:"favorite_color"`
	// 志向
	Ambition string `json:"ambition" db:"ambition"`

	Updatetime int64 `db:"updatetime"`

	Weight float64 `db:"weight"`

	Height float64 `db:"height"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GMyBabies) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()
	stmt, err := db.PrepareNamed(updateSql())
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
func (self *GMyBabies) Delete(query *DeleteQuery) error {
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
func (self *GMyBabies) ToggleDisabled(query *DisabledQuery) error {
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
