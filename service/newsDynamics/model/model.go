package model

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"
	"worldbar/DB/pgsql"
	"worldbar/util"
	"worldbar/util/database"

	"github.com/jmoiron/sqlx"
)

type WbNewsDynamics struct {
	database.BaseColumns

	Title   string `json:"title" db:"title"`
	Intro   string `json:"intro" db:"intro"`
	Surface string `json:"surface" db:"surface"`
	Content string `json:"content" db:"content"`

	Publisher string `json:"publisher" db:"publisher"`
	Type      int    `json:"type" db:"type"`

	Sort int `json:"sort" db:"sort"`

	Status       int    `json:"status" db:"status"`
	StatusReason string `json:"statusReason" db:"status_reason"`
	PublishTime  int64  `json:"publishTime" db:"publish_time"`
}

func (self *WbNewsDynamics) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Title) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Surface) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Content) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Publisher) == "" {
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
	WbNewsDynamics
	Like         bool `json:"like" db:"like"`
	TotalLike    int  `json:"totalLike" db:"total_like"`
	TotalComment int  `json:"totalComment" db:"total_comment"`
}

func (self *WbNewsDynamics) GetByID(query *GetQuery) (*GetModel, error) {
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
	WbNewsDynamics
	Avatar   string `json:"avatar" db:"avatar"`
	Nickname string `json:"nickname" db:"nickname"`
	Like     bool   `json:"like" db:"like"`
}

func (self *WbNewsDynamics) List(query *Query) ([]*ListModel, int64, error) {
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

func (self *WbNewsDynamics) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	ID      string `db:"id"`
	Title   string `db:"title"`
	Intro   string `db:"intro"`
	Content string `db:"content"`
	Surface string `db:"surface"`

	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *WbNewsDynamics) Update(query *UpdateByIDQuery) error {
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
func (self *WbNewsDynamics) Delete(query *DeleteQuery) error {
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
func (self *WbNewsDynamics) ToggleDisabled(query *DisabledQuery) error {
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
func (self *WbNewsDynamics) UpdateSort(query *UpdateSortQuery) error {
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
	PublishTime  int64  `db:"publish_time"`
}

// 更新状态
func (self *WbNewsDynamics) UpdateStatus(query *UpdateStatusQuery) error {
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
		query.PublishTime = util.GetCurrentMS()
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
