package voteOptionM

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
	"time"
	"worldbar/DB/pgsql"
)

type WbVoteOption struct {
	ID         string        `json:"id" db:"id"`
	Createtime int64         `json:"createtime" db:"createtime"`
	Updatetime sql.NullInt64 `json:"updatetime" db:"updatetime"`
	Deletetime sql.NullInt64 `json:"deletetime" db:"deletetime"`
	Isdelete   bool          `json:"isdelete" db:"isdelete"`
	Disabled   bool          `json:"disabled" db:"disabled"`

	VoteID string `json:"voteId" db:"vote_id"` // 投票外键

	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
}

func (self *WbVoteOption) Insert() (string, error) {

	if strings.TrimSpace(self.VoteID) == "" {
		return "", errors.New("参数错误")
	}

	if strings.TrimSpace(self.Title) == "" {
		return "", errors.New("选项标题不能为空")
	}

	// 插入时间
	self.Createtime = time.Now().UnixNano()

	// 默认添加直接启用
	self.Disabled = false
	self.Isdelete = false
	db := pgsql.Open()
	// 插入判断店铺是否已经存在
	stmt, err := db.PrepareNamed("insert into wb_vote_option (id, createtime, isdelete, disabled, vote_id, title, comment) select :id, :createtime, :isdelete, :disabled, :vote_id, :title, :comment returning id")

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
	ID string `db:"id"`
}

func (self *WbVoteOption) GetByID(db *sqlx.DB, query *GetQuery) (*WbVoteOption, error) {
	stmt, err := db.PrepareNamed("select * from wb_vote_option where id=:id")
	if err != nil {
		return nil, err
	}
	var email = new(WbVoteOption)
	err = stmt.Get(email, query)
	if err != nil {
		return nil, err
	}
	return email, nil
}

type Query struct {
	pgsql.BaseQuery
	VoteID string `json:"voteId" db:"vote_id"` // 投票外键
}

func (self *WbVoteOption) List(query *Query) ([]*WbVoteOption, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	var selectSql = `SELECT * FROM wb_vote_option WHERE 1=1 `
	var whereSql = ""
	whereSql = pgsql.BaseWhere(query.BaseQuery)

	if strings.TrimSpace(self.VoteID) != "" {
		whereSql = whereSql + " and vote_id=:vote_id"
	}

	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}

	var optionSql = ""
	if strings.TrimSpace(query.OrderBy) != "" {
		optionSql = optionSql + ` order by ` + query.OrderBy
	} else {
		optionSql = optionSql + ` order by id`
	}
	if strings.TrimSpace(query.SortFlag) != "" {
		optionSql = optionSql + ` ` + query.SortFlag
	} else {
		optionSql = optionSql + ` desc`
	}
	if query.Current > 0 {
		optionSql = optionSql + ` limit :pagesize`
	}
	query.Offset = (query.Current - 1) * query.PageSize

	if query.Offset > 0 {
		optionSql = optionSql + ` offset :offset`
	}
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

	var voteOptions = make([]*WbVoteOption, 0)
	for rows.Next() {
		var vo = new(WbVoteOption)
		err = rows.StructScan(&vo)
		if err != nil {
			return nil, 0, err
		}
		voteOptions = append(voteOptions, vo)
	}

	return voteOptions, count, nil

}

func (self *WbVoteOption) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from wb_vote_option where 1=1 " + strings.Join(whereSql, " "))
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
func (self *WbVoteOption) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()

	var updateSql = ""
	if strings.TrimSpace(query.Title) != "" {
		updateSql = updateSql + " ,title=:title"
	}
	if strings.TrimSpace(query.Comment) != "" {
		updateSql = updateSql + " ,comment=:comment"
	}

	var whereSql = " where isdelete=false and id=:id"
	// 如果条件有用户id，非总后台操作的会在权限部分解析出用户id，并且传过来

	stmt, err := db.PrepareNamed("update wb_vote_option set updatetime=:updatetime  " + updateSql + whereSql)
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
func (self *WbVoteOption) Delete(query *DeleteQuery) error {
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
	var whereSql = "where id=any(:ids)"
	stmt, err := db.PrepareNamed("update wb_vote_option set isdelete=true, deletetime=:deletetime " + whereSql)
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
func (self *WbVoteOption) ToggleDisabled(db *sqlx.DB, query *DisabledQuery) error {
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

	stmt, err := db.PrepareNamed("update wb_vote_option set disabled=:disabled, updatetime=:updatetime " + whereSql)
	if err != nil {
		return err
	}
	query.Updatetime = time.Now().UnixNano()
	_, err = stmt.Exec(query)
	return err
}
