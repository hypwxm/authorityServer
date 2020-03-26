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

type WbNewsDynamicsComment struct {
	database.BaseColumns

	Content string `json:"content" db:"content"`

	Publisher string `json:"publisher" db:"publisher"`
	NewsId    string `json:"newsId" db:"news_id"`

	PrevPublisher string `json:"prevPublisher" db:"prev_publisher"`
	PrevCommentId string `json:"prevCommentId" db:"prev_comment_id"`
	TopCommentId  string `json:"topCommentId" db:"top_comment_id"`
}

func (self *WbNewsDynamicsComment) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Content) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Publisher) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.NewsId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}

	// 插入时间
	self.Createtime = util.GetCurrentMS()

	// 默认添加直接启用
	self.Disabled = false
	self.Isdelete = false

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	// 插入判断用户登录账号是否已经存在
	stmt, err := tx.PrepareNamed("insert into wb_news_dynamics_comment (createtime, isdelete, disabled, id, content, publisher, news_id) select :createtime, :isdelete, :disabled, :id, :content, :publisher, :news_id returning id")

	if err != nil {
		return "", err
	}
	log.Println(stmt.QueryString)
	var lastId string
	self.ID = util.GetUuid()
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

func (self *WbNewsDynamicsComment) GetByID(query *GetQuery) (*WbNewsDynamicsComment, error) {
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("select * from wb_news_dynamics_comment where id=:id")
	if err != nil {
		return nil, err
	}
	var entity = new(WbNewsDynamicsComment)
	err = stmt.Get(entity, query)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

type Query struct {
	pgsql.BaseQuery
	PrevCommentId string `db:"prev_comment_id"`
	TopCommentId  string `db:"top_comment_id"`
	PrevPublisher string `db:"prev_publisher"`
	NewsId        string `db:"news_id"`
}

type ListModel struct {
	WbNewsDynamicsComment
	SubComments     []WbNewsDynamicsComment
	SubCommentCount int64
	Avatar          string `json:"avatar" db:"avatar"`
	Nickname        string `json:"nickname" db:"nickname"`
	PrevAvatar      string `json:"prevAvatar" db:"prev_avatar"`
	PrevNickname    string `json:"prevNickname" db:"prev_nickname"`
}

func (self *WbNewsDynamicsComment) List(query *Query) ([]*ListModel, int64, error) {
	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()

	list, total, err := self.list(tx, query)
	if err == nil {
		err = tx.Commit()
		if err != nil {
			return nil, 0, err
		}
	}
	return list, total, err
}

func (self *WbNewsDynamicsComment) list(tx *sqlx.Tx, query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	if query.NewsId == "" {
		return nil, 0, errors.New("参数错误")
	}

	whereSql, sqlStr := listSql(query)
	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(tx, query, whereSql)
	if err != nil {
		return nil, 0, err
	}
	stmt, err := tx.PrepareNamed(sqlStr)
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
	for rows.Next() {
		var data = new(ListModel)
		err = rows.StructScan(&data.WbNewsDynamicsComment)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, data)
	}

	// 获取下级评论的条数
	for _, v := range list {
		_query := *query
		_query.TopCommentId = v.ID
		whereSql, _ := listSql(&_query)
		count, err := self.GetCount(tx, &_query, whereSql)
		if err != nil {
			return nil, 0, err
		}
		v.SubCommentCount = count
	}

	return list, count, nil
}

func (self *WbNewsDynamicsComment) GetCount(db *sqlx.Tx, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from wb_news_dynamics_comment where 1=1 " + strings.Join(whereSql, " "))
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
	Content string `db:"content"`

	Publisher string `db:"publisher"`

	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *WbNewsDynamicsComment) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()
	var updateSql = ""
	updateSql = updateSql + " ,content=:content"

	// 评论者自己更新评论
	if query.Publisher != "" {
		updateSql = updateSql + " ,publisher=:publisher"
	}

	stmt, err := db.PrepareNamed("update wb_news_dynamics_comment set updatetime=:updatetime " + updateSql + " where id=:id and isdelete=false")
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
	IDs       pq.StringArray `db:"ids"`
	Publisher string         `db:"publisher"`
}

// 删除，批量删除
func (self *WbNewsDynamicsComment) Delete(query *DeleteQuery) error {
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
	var sqlStr string
	sqlStr = "update wb_news_dynamics_comment set isdelete=true where id=any(:ids)"
	if self.Publisher != "" {
		// 自己删除自己的评论
		sqlStr = "update wb_news_dynamics_comment set isdelete=true where id=any(:ids) and publisher=:publisher"
	}
	stmt, err := db.PrepareNamed(sqlStr)
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

// 启用禁用
func (self *WbNewsDynamicsComment) ToggleDisabled(query *DisabledQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("操作条件错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("update wb_news_dynamics_comment set disabled=:disabled where id=:id and isdelete=false")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}
