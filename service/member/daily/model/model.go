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

const BusinessName = "g_growning"

type GDaily struct {
	database.BaseColumns

	// 今日份体重
	Weight float64 `json:"weight" db:"weight"`
	// 今日份身高
	Height float64 `json:"height" db:"height"`

	// 今日份记录
	Diary string `json:"diary" db:"diary"`

	UserId string `json:"userId" db:"user_id"`
	BabyId string `json:"babyId" db:"baby_id"`

	Weather     string  `json:"weather" db:"weather"`
	Mood        string  `json:"mood" db:"mood"`
	Health      string  `json:"health" db:"health"`
	Temperature float64 `json:"temperature" db:"temperature"`

	Sort int `json:"sort" db:"sort"`

	Medias []*mediaModel.Media `json:"medias"`
}

func (self *GDaily) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.UserId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.BabyId) == "" {
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

	medias := mediaService.InitMedias(self.Medias, BusinessName, lastId, self.UserId)
	err = mediaService.MultiCreate(medias)
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
	GDaily
}

func (self *GDaily) GetByID(query *GetQuery) (*GetModel, error) {
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
	UserId   string `db:"user_id"`
}

type ListModel struct {
	GDaily
}

func (self *GDaily) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	if query.UserId == "" {
		return nil, 0, fmt.Errorf("参数错误")
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
		ids = append(ids, item.ID)
		list = append(list, item)
	}

	// 查找对应的媒体信息
	mediaService.ListWithMedia(ids, BusinessName, list, "")

	return list, count, nil

}

func (self *GDaily) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	ID      string `db:"id"`
	Title   string `db:"title"`
	Intro   string `db:"intro"`
	Content string `db:"content"`
	Surface string `db:"surface"`

	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GDaily) Update(query *UpdateByIDQuery) error {
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
func (self *GDaily) Delete(query *DeleteQuery) error {
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
