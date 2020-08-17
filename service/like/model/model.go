package model

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"
	"babygrowing/DB/pgsql"
	"babygrowing/util"
	"babygrowing/util/database"

	"github.com/jmoiron/sqlx"
)

type WbLike struct {
	database.BaseColumns

	UserId     string `json:"userId" db:"user_id"`
	SourceType int    `json:"sourceType" db:"source_type"`
	SourceId   string `json:"sourceId" db:"source_id"`
}

const (
	_ = iota
	SourceTypeNews
	SourceTypeUser
	SourceTypeComment
	SourceTypeMatter
)

func (self *WbLike) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.UserId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if self.SourceType == 0 {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.SourceId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	// 插入判断用户登录账号是否已经存在
	stmt, err := tx.PrepareNamed(insertSql(self))
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

func (self *WbLike) GetByID(query *GetQuery) (*WbLike, error) {
	db := pgsql.Open()
	stmt, err := db.PrepareNamed(getByIdSql())
	if err != nil {
		return nil, err
	}
	var entity = new(WbLike)
	err = stmt.Get(entity, query)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

type Query struct {
	pgsql.BaseQuery
	SourceType int `db:"source_type"`
}

type ListModel struct {
	WbLike
	Avatar       string `json:"avatar" db:"avatar"`
	Nickname     string `json:"nickname" db:"nickname"`
	LikeAvatar   string `json:"likeAvatar" db:"like_avatar"`
	LikeNickname string `json:"likeNickname" db:"like_nickname"`
	NewsTitle    string `json:"newsTitle" db:"news_title"`
	NewsSurface  string `json:"newsSurface" db:"news_surface"`
}

func (self *WbLike) List(query *Query) ([]*ListModel, int64, error) {
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

func (self *WbLike) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
func (self *WbLike) Update(query *UpdateByIDQuery) error {
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
func (self *WbLike) Delete(query *DeleteQuery) error {
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
func (self *WbLike) ToggleDisabled(query *DisabledQuery) error {
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
func (self *WbLike) UpdateSort(query *UpdateSortQuery) error {
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
func (self *WbLike) UpdateStatus(query *UpdateStatusQuery) error {
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
