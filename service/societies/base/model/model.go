/**
小区社团服务
*/
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

/*
	社团基本信息
*/
type WbSocieties struct {
	database.BaseColumns

	// 社团名称
	Name string `json:"name" db:"name"`
	// 社团介绍
	Intro string `json:"intro" db:"intro"`
	// 社团头像
	Surface string `json:"surface" db:"surface"`

	// 社团创建人
	Creator string `json:"creator" db:"creator"`
	// 社团创建人id，只有本小区的人允许创建社团
	CreatorId string `json:"creatorId" db:"creator_id"`

	// 社团最大限制人数, 0表示无限制
	PeopleMax int `json:"peopleMax" db:"people_max"`

	/*
		社团类型，需要维护一个社团类型枚举
	*/
	TypeId   string `json:"typeId" db:"type_id"`
	TypeName string `json:"typeName" db:"type_name"`

	/*
		社团状态，1:创建成功，2:提交审核，审核中，3:撤回，4:驳回，5:通过，6:发布，7:关闭，8:维护中
		在任何情况下如果社团的基本信息被修改，社团的状态都会回到1或2状态，附加属性不影响（如标签）
	*/
	Status int `json:"status" db:"status"`
	/*
		驳回通过原因
	*/
	StatusReason string `json:"statusReason" db:"status_reason"`
	/*
		社团最新发布时间
	*/
	PublishTime int64 `json:"publishTime" db:"publish_time"`
}

func (self *WbSocieties) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Name) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.CreatorId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.TypeId) == "" {
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
	WbSocieties
}

func (self *WbSocieties) GetByID(query *GetQuery) (*GetModel, error) {
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
	WbSocieties
}

func (self *WbSocieties) List(query *Query) ([]*ListModel, int64, error) {
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

func (self *WbSocieties) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	ID        string `db:"id"`
	Name      string `db:"name"`
	Intro     string `db:"intro"`
	PeopleMax int    `db:"people_max"`
	TypeId    int    `db:"type_id"`
	TypeName  int    `db:"type_name"`
	Status    int    `db:"status"`

	Updatetime int64 `db:"updatetime"`
}

/**
	可修改：创建成功，维护中，审核通过但未发布
	提交审核的，可撤回提交再进行修改

 */
func (self *WbSocieties) Update(query *UpdateByIDQuery) error {
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
func (self *WbSocieties) Delete(query *DeleteQuery) error {
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
func (self *WbSocieties) ToggleDisabled(query *DisabledQuery) error {
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
func (self *WbSocieties) UpdateSort(query *UpdateSortQuery) error {
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
func (self *WbSocieties) UpdateStatus(query *UpdateStatusQuery) error {
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
		if query.PublishTime == 0 {
			query.PublishTime = util.GetCurrentMS()
		}
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
