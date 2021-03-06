package model

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/DB/pgsql"
	"github.com/hypwxm/authorityServer/util/interfaces"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Media struct {
	appGorm.BaseColumns
	Url        string `db:"url" json:"url" gorm:"column:url;type:varchar(500);not null;default ''"`
	UserID     string `db:"user_id" json:"userId" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '')"`
	Business   string `db:"business" json:"business" gorm:"column:business;type:varchar(50);not null;default ''"`
	BusinessId string `db:"business_id" json:"businessId" gorm:"column:business_id;type:varchar(128);not null;check(business_id <> '')"`
	Size       int    `db:"size" json:"size" gorm:"column:size;not null;default 0"`
}

func (self *Media) Insert(tx ...*sqlx.Tx) (string, error) {
	var err error

	if strings.TrimSpace(self.UserID) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Business) == "" {
		return "", errors.New(fmt.Sprintf("未知业务"))
	}
	if strings.TrimSpace(self.BusinessId) == "" {
		return "", errors.New(fmt.Sprintf("未知业务"))
	}

	var ltx *sqlx.Tx
	var stmt *sqlx.NamedStmt
	if len(tx) > 0 {
		ltx = tx[0]
	}
	if ltx != nil {
		stmt, err = ltx.PrepareNamed(insertSql())
	} else {
		db := pgsql.Open()
		stmt, err = db.PrepareNamed(insertSql())
	}

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
	return self.ID, nil
}

/**
 * 初始化媒体信息
 * 前端回传过来url和size信息，
 * business和businessId由后台对应的业务生成，
 * creator为当前的登陆信息中取
 *
 * 初始化的时候需要考虑 list参数中的其中对象中是否已经存在business，id等信息，
 * 因为对业务进行编辑的时候，涉及到媒体时，会把原来的媒体全部删除，再重新插入，那么有以上信息的话，就不对这个媒体信息重新设置数据，而是直接插入
 */
func InitMedias(list []*Media, businessName string, businessId string, creator string) []*Media {
	medias := make([]*Media, 0)
	for _, v := range list {
		if v == nil {
			return nil
		}
		// 有一种情况是这个媒体是之前的信息，就不进行重新保存了
		if !(v.Business != "" && v.BusinessId != "" && v.UserID != "" && v.ID != "") {
			v.Init()
			v.Business = businessName
			v.BusinessId = businessId
			v.UserID = creator
		}
		medias = append(medias, v)
	}
	return medias
}

func StoreMedias(list []*Media, tx ...*sqlx.Tx) error {
	var err error
	for _, v := range list {
		if v == nil {
			return fmt.Errorf("文件错误")
		}
		_, err = v.Insert(tx...)
		if err != nil {
			return err
		}
	}
	return nil
}

type Query struct {
	pgsql.BaseQuery
	BusinessIds pq.StringArray `json:"businessIds" db:"business_ids"`
	Businesses  pq.StringArray `json:"businesses" db:"businesses"`
}

func (self *Media) List(query *Query) ([]*Media, int, error) {
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

	var list = make([]*Media, 0)
	for rows.Next() {
		var item = new(Media)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
	}
	return list, count, nil
}

func (self *Media) ListWithMedia(query *Query, olist interface{}, mediaName string, keyName string) error {
	medias, _, err := self.List(query)

	if err != nil {
		return err
	}
	if mediaName == "" {
		mediaName = "Medias"
	}
	if keyName == "" {
		keyName = "ID"
	}

	// log.Println(reflect.ValueOf(olist).Kind())
	// log.Println(reflect.ValueOf(olist))
	// log.Println(reflect.TypeOf(olist))
	// log.Println(reflect.ValueOf(olist).NumField())

	v := reflect.ValueOf(olist)
	if v.Kind() == reflect.Slice {
		len := v.Len()
		ret := make([]interface{}, len)
		for i := 0; i < len; i++ {
			ret[i] = v.Index(i).Interface()
		}
		for _, v := range ret {
			// log.Println(reflect.ValueOf(v).Kind())
			nv := reflect.ValueOf(v)
			if nv.Kind() == reflect.Ptr {
				// log.Println(nv.Elem(), reflect.TypeOf(v).Elem())
				t := nv.Elem()
				// fmt.Println("Number of fields", v.NumField())
				// v = v.Elem()
				m := t.FieldByName(mediaName)
				if m.Kind() == reflect.Invalid {
					return fmt.Errorf("mediaName定义错误")
				}
				// log.Println(m.Interface().([]*Media))
				if ll, ok := m.Interface().([]*Media); ok {
					for _, v := range medias {
						if v.BusinessId == t.FieldByName(keyName).String() {
							ll = append(ll, v)
							m.Set(reflect.ValueOf(ll))
						}
					}
				}
			}

		}
	}
	return nil
}

func (self *Media) ListMapWithMedia(query *Query, olist interfaces.ModelMapSlice, mediaName string, keyName string) error {
	medias, _, err := self.List(query)

	if err != nil {
		return err
	}
	if mediaName == "" {
		mediaName = "medias"
	}
	if keyName == "" {
		keyName = "id"
	}

	for _, v := range olist {
		v.Set(mediaName, make([]*Media, 0))
		for _, vm := range medias {
			if v.GetValue(keyName) == vm.BusinessId {
				v.Set(mediaName, append(v.GetValue(mediaName).([]*Media), vm))
			}
		}
	}
	return nil
}

func (self *Media) ListMapWithMediaFirst(query *Query, olist interfaces.ModelMapSlice, keyName string, valueKey string) error {
	medias, _, err := self.List(query)
	if err != nil {
		return err
	}
	for _, v := range olist {
		for _, vm := range medias {
			if v.GetValue(keyName) == vm.BusinessId {
				v.Set(valueKey, vm.Url)
				break
			}
		}
	}
	return nil
}

func (self *Media) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int, error) {
	if query == nil {
		query = new(Query)
	}
	sqlStr := countSql(whereSql...)
	stmt, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return 0, err
	}
	var count int
	err = stmt.Get(&count, query)
	log.Println(stmt.QueryString, query)
	return count, err
}

type DeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
	// 存储媒体的业务名称，一般为对应的业务的主表的名称
	Businesses pq.StringArray `json:"businesses" db:"businesses"`
	// 业务主表存储的对应的主键id
	BusinessIds pq.StringArray `json:"businessIds" db:"business_ids"`
}

// 删除，批量删除
func (self *Media) Delete(query *DeleteQuery, tx ...*sqlx.Tx) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if len(query.IDs) == 0 && len(query.BusinessIds) == 0 {
		return errors.New("操作条件错误")
	}

	var ltx *sqlx.Tx
	if len(tx) > 0 {
		ltx = tx[0]
	}
	var stmt *sqlx.NamedStmt
	var err error
	if ltx != nil {
		// 通过事务来删除
		stmt, err = ltx.PrepareNamed(delSql(query))
	} else {
		db := pgsql.Open()
		stmt, err = db.PrepareNamed(delSql(query))
	}
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}
