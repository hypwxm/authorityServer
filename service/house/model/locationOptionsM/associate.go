package locationOptionsM

import (
	"errors"
	"log"
	"strings"
	"worldbar/DB/pgsql"
)

type WbHouseOptionAssociate struct {
	SuperOptionId string `db:"super_option_id"`
	SubOptionId   string `db:"sub_option_id"`
}

type AssociateQuery struct {
	List []WbHouseOptionAssociate
}

// 对上下级的属性进行关联
func (self *WbHouseOptionAssociate) Associate(query *AssociateQuery) error {
	if query == nil {
		return errors.New("参数错误")
	}

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, v := range query.List {
		if strings.TrimSpace(v.SuperOptionId) == "" || strings.TrimSpace(v.SubOptionId) == "" || strings.TrimSpace(v.SuperOptionId) == strings.TrimSpace(v.SubOptionId) {
			return errors.New("参数错误")
		}
		stmt, err := db.PrepareNamed("insert into wb_house_option_associate (super_option_id, sub_option_id) select :super_option_id, :sub_option_id where not exists(select 1 from wb_house_option_associate where (super_option_id=:super_option_id and sub_option_id=:sub_option_id) or (super_option_id=:sub_option_id and sub_option_id=:super_option_id)) returning true")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(&v)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

// 对上下级的属性进行关联
func (self *WbHouseOptionAssociate) DeleteAssociates(query *AssociateQuery) error {
	if query == nil {
		return errors.New("参数错误")
	}

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, v := range query.List {
		if strings.TrimSpace(v.SuperOptionId) == "" || strings.TrimSpace(v.SubOptionId) == "" || strings.TrimSpace(v.SuperOptionId) == strings.TrimSpace(v.SubOptionId) {
			return errors.New("参数错误")
		}
		stmt, err := db.PrepareNamed("delete from wb_house_option_associate where super_option_id=:super_option_id and sub_option_id=:sub_option_id")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(&v)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

type AssociateGetQuery struct {
	OptionId          string `db:"option_id"`            // 查询上级或下级其中之一满足条件的，存在option_id的时候不应存在super_option_id和sub_option_id
	SuperOptionId     string `db:"super_option_id"`      // 查询上级为该条件的，可与sub_option_id同时存在，相当于查询一条数据
	SubOptionId       string `db:"sub_option_id"`        // 查询下级为该条件的
	SuperHouseEnumsId string `db:"super_house_enums_id"` // 查询上级的枚举为该条件的
	SubHouseEnumsId   string `db:"sub_house_enums_id"`   // 查询下级的枚举为该条件的
}

type Associate struct {
	SuperOptionId  string `json:"superOptionId" db:"super_option_id"`
	SubOptionId    string `json:"subOptionId" db:"sub_option_id"`
	SuperName      string `json:"superName" db:"super_name"`
	SuperNote      string `json:"superNote" db:"super_note"`
	SubName        string `json:"subName" db:"sub_name"`
	SubNote        string `json:"subNote" db:"sub_note"`
	SuperEnumsId   string `json:"superEnumsId" db:"super_enums_id"`
	SuperEnumsName string `json:"superEnumsName" db:"super_enums_name"`
	SubEnumsId     string `json:"subEnumsId" db:"sub_enums_id"`
	SubEnumsName   string `json:"subEnumsName" db:"sub_enums_name"`
}

func (self *WbHouseOptionAssociate) GetAssociate(query *AssociateGetQuery) ([]*Associate, error) {
	if query == nil {
		return nil, errors.New("参数错误")
	}
	db := pgsql.Open()
	var whereSql = ""
	if query.SuperOptionId != "" {
		whereSql = whereSql + " and a.super_option_id=:super_option_id"
	}
	if query.SubOptionId != "" {
		whereSql = whereSql + " and a.sub_option_id=:sub_option_id"
	}
	if query.OptionId != "" {
		whereSql = whereSql + " and (a.super_option_id=:option_id or a.sub_option_id=:option_id)"
	}
	if query.SuperHouseEnumsId != "" {
		whereSql = whereSql + " and b.house_enums_id=:super_house_enums_id"
	}
	if query.SubHouseEnumsId != "" {
		whereSql = whereSql + " and c.house_enums_id=:sub_house_enums_id"
	}

	stmt, err := db.PrepareNamed("select a.super_option_id, a.sub_option_id, b.name as super_name, b.note as super_note, b.house_enums_id as super_enums_id, c.name as sub_name, c.note as sub_note, c.house_enums_id as sub_enums_id from wb_house_option_associate a inner join wb_house_option b on a.super_option_id=b.id inner join wb_house_option c on a.sub_option_id=c.id where 1=1 and b.isdelete=false and c.isdelete=false " + whereSql)

	if err != nil {
		return nil, err
	}
	log.Println(stmt.QueryString)

	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list = make([]*Associate, 0)
	for rows.Next() {
		var data = new(Associate)
		err = rows.StructScan(&data)
		if err != nil {
			return nil, err
		}
		list = append(list, data)
	}
	return list, nil
}

type AssociateTree struct {
	Name             string           `json:"name"`
	OptionId         string           `json:"optionId"`
	OptionName       string           `json:"optionName"`
	EnumsId          string           `json:"enumsId"`
	EnumsName        string           `json:"enumsName"`
	ParentOptionId   string           `json:"parentOptionId"`
	ParentOptionName string           `json:"parentOptionName"`
	ParentEnumsId    string           `json:"parentEnumsId"`
	ParentEnumsName  string           `json:"parentEnumsName"`
	Children         []*AssociateTree `json:"children"`
}
