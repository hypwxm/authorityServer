package voteM

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
	"time"
	"babygrowing/DB/pgsql"
)

type WbVote struct {
	ID         string        `db:"id" json:"id"`
	Createtime int64         `db:"createtime" json:"createtimecreatetime"`
	Updatetime sql.NullInt64 `db:"updatetime" json:"updatetime"`
	Deletetime sql.NullInt64 `db:"deletetime" json:"deletetime"`
	Isdelete   bool          `db:"isdelete" json:"isdelete"`
	Disabled   bool          `db:"disabled" json:"disabled"`

	Title   string `db:"title" json:"title"`     // 投票标题
	Comment string `db:"comment" json:"comment"` // 投票备注
}

func (self *WbVote) Insert() (string, error) {

	if strings.TrimSpace(self.Title) == "" {
		return "", errors.New("投票主题不能为空")
	}

	// 插入时间
	self.Createtime = time.Now().UnixNano()

	// 默认添加直接启用
	self.Disabled = false
	self.Isdelete = false

	db := pgsql.Open()

	// 插入判断店铺是否已经存在
	stmt, err := db.PrepareNamed("insert into wb_vote (id, createtime, isdelete, disabled, title, comment) select :id, :createtime, :isdelete, :disabled, :title, :comment returning id")

	log.Println(stmt.QueryString)

	if err != nil {
		return "", err
	}
	var lastId string
	self.ID = uuid.NewV4().String()
	err = stmt.Get(&lastId, self)
	if err != nil {
		return "", err
	}
	return lastId, nil
}

type GetQuery struct {
	ID string `db:"id" json:"id"`
}

func (self *WbVote) GetByID(query *GetQuery) (*WbVote, error) {
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("select * from wb_vote where id=:id")
	if err != nil {
		return nil, err
	}
	var email = new(WbVote)
	err = stmt.Get(email, query)
	if err != nil {
		return nil, err
	}
	return email, nil
}

type Query struct {
	pgsql.BaseQuery

	Title   string `db:"title"`
	Comment string `db:"comment"`
}

func (self *WbVote) List(query *Query) ([]*WbVote, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	var selectSql = `SELECT * FROM wb_vote WHERE 1=1 `
	whereSql := pgsql.BaseWhere(query.BaseQuery)

	if strings.TrimSpace(query.Title) != "" {
		whereSql = whereSql + " and title like '%" + query.Title + "%'"
	}

	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}

	var optionSql = pgsql.BaseOption(query.BaseQuery)
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

	var votes = make([]*WbVote, 0)
	for rows.Next() {
		var vote = new(WbVote)
		err = rows.StructScan(&vote)
		if err != nil {
			return nil, 0, err
		}
		votes = append(votes, vote)
	}

	return votes, count, nil

}

func (self *WbVote) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from wb_vote where 1=1 " + strings.Join(whereSql, " "))
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
	Updatetime int64  `db:"updatetime"`
	Title      string `db:"title"`
	Comment    string `db:"comment"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *WbVote) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	var updateSql = ""
	if strings.TrimSpace(query.Title) != "" {
		updateSql = updateSql + " ,title=:title"
	}
	if strings.TrimSpace(query.Comment) != "" {
		updateSql = updateSql + " ,comment=:comment"
	}

	db := pgsql.Open()

	var whereSql = " where isdelete=false and id=:id"

	stmt, err := db.PrepareNamed("update wb_vote set updatetime=:updatetime  " + updateSql + whereSql)
	if err != nil {
		return err
	}
	query.Updatetime = time.Now().UnixNano()
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

type DeleteQuery struct {
	IDs        pq.StringArray `db:"ids"`
	Deletetime int64          `db:"deletetime"`
}

// 删除，批量删除
func (self *WbVote) Delete(db *sqlx.DB, query *DeleteQuery) error {
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
	var whereSql = "where id=any(:ids)"

	stmt, err := db.PrepareNamed("update wb_vote set isdelete=true, deletetime=:deletetime " + whereSql)
	if err != nil {
		return err
	}
	query.Deletetime = time.Now().UnixNano()
	_, err = stmt.Exec(query)
	return err
}

type DisabledQuery struct {
	Disabled   bool           `db:"disabled"`
	IDs        pq.StringArray `db:"ids"`
	Updatetime int64          `db:"updatetime"`
}

// 启用禁用店铺
func (self *WbVote) ToggleDisabled(db *sqlx.DB, query *DisabledQuery) error {
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
	var whereSql = "where id=any(:ids) and isdelete=false"

	stmt, err := db.PrepareNamed("update wb_vote set disabled=:disabled, updatetime=:updatetime " + whereSql)
	if err != nil {
		return err
	}
	query.Updatetime = time.Now().UnixNano()
	_, err = stmt.Exec(query)
	return err
}
