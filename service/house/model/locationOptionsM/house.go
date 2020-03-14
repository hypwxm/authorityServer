package locationOptionsM

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

type WbHouseOption struct {
	ID         string        `json:"id" db:"id"`
	Createtime int64         `json:"createtime" db:"createtime"`
	Updatetime sql.NullInt64 `json:"updatetime" db:"updatetime"`
	Deletetime sql.NullInt64 `json:"deletetime" db:"deletetime"`
	Isdelete   bool          `json:"isdelete" db:"isdelete"`
	Disabled   bool          `json:"disabled" db:"disabled"`

	Name         string `json:"name" db:"name"` // 枚举名称，比如 区，房，室等
	Note         string `json:"note" db:"note"` // 房号
	HouseEnumsId string `json:"houseEnumsId" db:"house_enums_id"`
}

func (self *WbHouseOption) Insert() (string, error) {
	var err error
	// 插入时间
	self.Createtime = util.GetCurrentMS()

	// 必须有登录账号
	if strings.TrimSpace(self.Name) == "" {
		return "", errors.New(fmt.Sprintf("名称不能为空"))
	}

	if strings.TrimSpace(self.HouseEnumsId) == "" {
		return "", errors.New(fmt.Sprintf("参数错误"))
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
	stmt, err := tx.PrepareNamed("insert into wb_house_option (createtime, isdelete, disabled, id, name, note, house_enums_id) select :createtime, :isdelete, :disabled, :id, :name, :note, :house_enums_id returning id")

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

func (self *WbHouseOption) GetByID(query *GetQuery) (*WbHouseOption, error) {
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("select * from wb_house_option where id=:id")
	if err != nil {
		return nil, err
	}
	var entity = new(WbHouseOption)
	err = stmt.Get(entity, query)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

type Query struct {
	pgsql.BaseQuery
	Name         string `db:"name"`
	HouseEnumsId string `db:"house_enums_id"`
}

func (self *WbHouseOption) List(query *Query) ([]*WbHouseOption, int64, error) {
	if query == nil {
		query = new(Query)
	}

	log.Printf("%v111", query)

	db := pgsql.Open()
	var selectSql = `SELECT * FROM wb_house_option WHERE 1=1 `
	var whereSql = ""
	whereSql = pgsql.BaseWhere(query.BaseQuery)

	if strings.TrimSpace(query.Name) != "" {
		whereSql = whereSql + " and name like '%" + query.Name + "%'"
	}
	if query.HouseEnumsId != "" {
		whereSql = whereSql + " and house_enums_id='" + query.HouseEnumsId + "'"
	}
	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}
	query.OrderBy = "createtime asc"
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

	var list = make([]*WbHouseOption, 0)
	for rows.Next() {
		var item = new(WbHouseOption)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}

	return list, count, nil

}

func (self *WbHouseOption) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	stmt, err := db.PrepareNamed("select count(*) from wb_house_option where 1=1 " + strings.Join(whereSql, " "))
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
func (self *WbHouseOption) Update(query *UpdateByIDQuery) error {
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

	stmt, err := db.PrepareNamed("update wb_house_option set updatetime=:updatetime " + updateSql + " where id=:id and isdelete=false")
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
func (self *WbHouseOption) Delete(query *DeleteQuery) error {
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
	stmt, err := db.PrepareNamed("update wb_house_option set isdelete=true where id=any(:ids)")
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
func (self *WbHouseOption) ToggleDisabled(query *DisabledQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("操作条件错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("update wb_house_option set disabled=:disabled where id=:id and isdelete=false")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}

type GetListByIdsQuery struct {
	IDs pq.StringArray `db:"ids"`
}

type OptionsWithEnums struct {
	EnumsId    string `json:"enumsId" db:"enums_id"`
	EnumsName  string `json:"enumsName" db:"enums_name"`
	OptionId   string `json:"optionId" db:"option_id"`
	OptionName string `json:"optionName" db:"option_name"`
}

// 根据ids获取详情
func (self *WbHouseOption) GetListByIds(query *GetListByIdsQuery) ([]*OptionsWithEnums, error) {
	if query == nil {
		return nil, errors.New("无操作条件")
	}
	if len(query.IDs) == 0 {
		return nil, errors.New("操作条件错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed("select a.id as option_id, a.name as option_name, b.id as enums_id, b.name as enums_name from wb_house_option a left join wb_house_enums b on a.house_enums_id=b.id where a.id=any(:ids) and a.isdelete=false")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list = make([]*OptionsWithEnums, 0)
	for rows.Next() {
		var item = new(OptionsWithEnums)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	return list, nil
}
