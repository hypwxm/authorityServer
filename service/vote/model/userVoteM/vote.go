package vote

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
	"time"
)

type WbUserVote struct {
	ID         string        `db:"id"`
	Createtime int64         `db:"createtime"`
	Updatetime sql.NullInt64 `db:"updatetime"`
	Deletetime sql.NullInt64 `db:"deletetime"`
	Isdelete   bool          `db:"isdelete"`
	Disabled   bool          `db:"disabled"`

	UserID   string `db:"user_id"`
	VoteID   string `db:"vote_id"`
	OptionID string `db:"option_id"` // 投票备注
}

func (self *WbUserVote) Insert(db *sqlx.DB) (int64, error) {
	// 插入时间
	self.Createtime = time.Now().UnixNano()

	// 默认添加直接启用
	self.Disabled = false
	self.Isdelete = false

	// 插入判断店铺是否已经存在
	stmt, err := db.PrepareNamed("insert into wb_user_vote (id, createtime, isdelete, disabled, user_id, vote_id, option_id) select :id, :createtime, :isdelete, :disabled, :user_id, :vote_id, :option_id returning id")

	log.Println(stmt.QueryString)

	if err != nil {
		return 0, err
	}
	var lastId int64
	self.ID = uuid.NewV4().String()
	err = stmt.Get(&lastId, self)
	if err != nil {
		return 0, err
	}
	return lastId, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

func (self *WbUserVote) GetByID(db *sqlx.DB, query *GetQuery) (*WbUserVote, error) {
	stmt, err := db.PrepareNamed("select * from wb_user_vote where id=:id")
	if err != nil {
		return nil, err
	}
	var email = new(WbUserVote)
	err = stmt.Get(email, query)
	if err != nil {
		return nil, err
	}
	return email, nil
}

type Query struct {
	IDs       pq.StringArray `db:"ids"`
	Current   int            `db:"current"`
	PageSize  int            `db:"pagesize"`
	Offset    int            `db:"offset"`
	Starttime int64          `db:"starttime"`
	Endtime   int64          `db:"endtime"`
	OrderBy   string         `db:"order_by"`
	SortFlag  string         `db:"sort_flag"`
}

func (self *WbUserVote) List(db *sqlx.DB, query *Query) ([]*WbUserVote, int64, error) {
	if query == nil {
		query = new(Query)
	}
	var selectSql = `SELECT * FROM wb_user_vote WHERE 1=1 `
	var whereSql = ""
	if query.IDs != nil {
		whereSql = whereSql + ` and id=any(:ids)`
	}
	if query.Starttime > 0 {
		whereSql = whereSql + ` and createtime>=:starttime`
	}
	if query.Endtime > 0 {
		whereSql = whereSql + ` and createtime<=:endtime`
	}

	whereSql = whereSql + ` and isdelete='false'`

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

	var votes = make([]*WbUserVote, 0)
	for rows.Next() {
		var vote = new(WbUserVote)
		err = rows.StructScan(&vote)
		if err != nil {
			return nil, 0, err
		}
		votes = append(votes, vote)
	}

	return votes, count, nil

}

func (self *WbUserVote) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from wb_user_vote where 1=1 " + strings.Join(whereSql, " "))
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
	OptionID   string `db:"option_id"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *WbUserVote) Update(db *sqlx.DB, query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	var updateSql = ""
	if strings.TrimSpace(query.OptionID) != "" {
		updateSql = updateSql + " ,option_id=:option_id"
	}

	var whereSql = " where isdelete=false and id=:id"

	stmt, err := db.PrepareNamed("update wb_user_vote set updatetime=:updatetime  " + updateSql + whereSql)
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
func (self *WbUserVote) Delete(db *sqlx.DB, query *DeleteQuery) error {
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

	stmt, err := db.PrepareNamed("update wb_user_vote set isdelete=true, deletetime=:deletetime " + whereSql)
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
func (self *WbUserVote) ToggleDisabled(db *sqlx.DB, query *DisabledQuery) error {
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

	stmt, err := db.PrepareNamed("update wb_user_vote set disabled=:disabled, updatetime=:updatetime " + whereSql)
	if err != nil {
		return err
	}
	query.Updatetime = time.Now().UnixNano()
	_, err = stmt.Exec(query)
	return err
}
