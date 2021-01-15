package model

import (
	"babygrowing/DB/pgsql"
	mediaModel "babygrowing/service/media/model"
	mediaService "babygrowing/service/media/service"

	"babygrowing/util"
	"babygrowing/util/database"

	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const BusinessName = "g_admin_org"

type GOrg struct {
	database.BaseColumns

	Name     string `json:"name" db:"name"`
	Sort     int    `json:"sort" db:"sort"`
	ParentId string `json:"parentId" db:"parent_id"`

	UserId string              `json:"userId" db:"-"`
	Media  []*mediaModel.Media `json:"media" db:"-"`
}

func (self *GOrg) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Name) == "" {
		return "", fmt.Errorf("操作错误")
	}

	db := pgsql.Open()
	stmt, err := db.PrepareNamed(insertSql())
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

	medias := mediaService.InitMedias(self.Media, BusinessName, lastId, self.UserId)
	err = mediaService.MultiCreate(medias)
	if err != nil {
		return "", err
	}

	return self.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

type GetModel struct {
	GOrg
}

func (self *GOrg) GetByID(query *GetQuery) (*GetModel, error) {
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
	GOrg
	Avatar string `json:"avatar"`
}

func (self *GOrg) List(query *Query) ([]*ListModel, int64, error) {
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

	var list = make([]*ListModel, 0)

	var ids []string = make([]string, 0)

	for rows.Next() {
		var item = new(ListModel)
		err = rows.StructScan(&item)
		if err != nil {
			return nil, 0, err
		}
		list = append(list, item)
		ids = append(ids, item.ID)
	}

	// 查找对应的媒体信息
	medias, _, err := mediaService.List(&mediaModel.Query{
		BusinessIds: ids,
		Businesses:  []string{BusinessName},
	})

	if err != nil {
		return nil, 0, err
	}

	for _, v := range list {
		for _, vm := range medias {
			if v.ID == vm.BusinessId {
				v.Media = append(v.Media, vm)
				v.Avatar = vm.Url
			}
		}
	}

	return list, count, nil

}

func (self *GOrg) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	ID    string              `db:"id"`
	Name  string              `db:"name"`
	Media []*mediaModel.Media `json:"media" db:"-"`

	Disabled   bool  `json:"disabled" db:"disabled"`
	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GOrg) Update(query *UpdateByIDQuery) error {
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
func (self *GOrg) Delete(query *DeleteQuery) error {
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

func (self *GOrg) ToggleDisabled(query *DisabledQuery) error {
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
