package locationEnumsM

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"strings"
	"worldbar/DB/pgsql"
	"worldbar/util"

	"github.com/jmoiron/sqlx"
)

type WbHouseEnums struct {
	ID         string        `json:"id" db:"id"`
	Createtime int64         `json:"createtime" db:"createtime"`
	Updatetime sql.NullInt64 `json:"updatetime" db:"updatetime"`
	Deletetime sql.NullInt64 `json:"deletetime" db:"deletetime"`
	Isdelete   bool          `json:"isdelete" db:"isdelete"`
	Disabled   bool          `json:"disabled" db:"disabled"`

	Name string `json:"name" db:"name"` // 枚举名称，比如 区，房，室等
	Note string `json:"note" db:"note"` // 房号
	Sort int    `json:"sort" db:"sort"`
}

func (self *WbHouseEnums) Insert() (string, error) {
	var err error
	// 插入时间
	self.Createtime = util.GetCurrentMS()

	// 必须有登录账号
	if strings.TrimSpace(self.Name) == "" {
		return "", errors.New(fmt.Sprintf("名称不能为空"))
	}

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
	stmt, err := tx.PrepareNamed("insert into wb_house_enums (createtime, isdelete, disabled, id, name, note) select :createtime, :isdelete, :disabled, :id, :name, :note returning id")

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

func (self *WbHouseEnums) GetByID(query *GetQuery) (*WbHouseEnums, error) {
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("select * from wb_house_enums where id=:id")
	if err != nil {
		return nil, err
	}
	var entity = new(WbHouseEnums)
	err = stmt.Get(entity, query)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

type Query struct {
	pgsql.BaseQuery
	Name string `db:"name"`
}

func (self *WbHouseEnums) List(query *Query) ([]*WbHouseEnums, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	var selectSql = `SELECT * FROM wb_house_enums WHERE 1=1 `
	var whereSql = ""
	query.OrderBy = ""
	whereSql = pgsql.BaseWhere(query.BaseQuery)

	if strings.TrimSpace(query.Name) != "" {
		whereSql = whereSql + " and name like '%" + query.Name + "%'"
	}

	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}
	query.OrderBy = "sort asc"
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

	var users = make([]*WbHouseEnums, 0)
	for rows.Next() {
		var user = new(WbHouseEnums)
		err = rows.StructScan(&user)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, count, nil

}

func (self *WbHouseEnums) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from wb_house_enums where 1=1 " + strings.Join(whereSql, " "))
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
	Name       string `db:"name"`
	Note       string `db:"note"`
	Updatetime int64  `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *WbHouseEnums) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()
	var updateSql = ""
	updateSql = updateSql + " ,name=:name"
	updateSql = updateSql + " ,note=:note"

	stmt, err := db.PrepareNamed("update wb_house_enums set updatetime=:updatetime " + updateSql + " where id=:id and isdelete=false")
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
func (self *WbHouseEnums) Delete(query *DeleteQuery) error {
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
	stmt, err := db.PrepareNamed("update wb_house_enums set isdelete=true where id=any(:ids)")
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
func (self *WbHouseEnums) ToggleDisabled(query *DisabledQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("操作条件错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("update wb_house_enums set disabled=:disabled where id=:id and isdelete=false")
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
func (self *WbHouseEnums) UpdateSort(query *UpdateSortQuery) error {
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
	stmt, err := tx.PrepareNamed(`update wb_house_enums set sort=:sort1 where id=:id2 and isdelete=false`)
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	stmt, err = tx.PrepareNamed("update wb_house_enums set sort=:sort2 where id=:id1 and isdelete=false")
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	_, err = stmt.Exec(query)
	err = tx.Commit()
	return err
}
