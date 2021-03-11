package model

import (
	"babygrow/DB/pgsql"
	mediaModel "babygrow/service/media/model"
	mediaService "babygrow/service/media/service"

	"babygrow/util"
	"babygrow/util/database"

	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const BusinessName = "g_message"

type GMessage struct {
	database.BaseColumns

	BusinessName string `json:"businessName" db:"business_name"`
	BusinessId   string `json:"businessId" db:"business_id"`

	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`

	// 消息发送时间，消息不一定是实时发送
	Sendtime int `json:"sendtime" db:"sendtime"`
	// 是否已读
	IsRead bool `json:"isRead" db:"is_read"`
	// 阅读时间
	ReadDuration float64 `json:"readDuration" db:"read_duration"`
	// 读到哪里了
	ReadPercent float64 `json:"readPercent" db:"read_percent"`

	// 发送人信息
	SenderId   string `json:"sendId" db:"send_id"`
	SenderName string `json:"senderName" db:"sender_name"`

	// 接受人信息
	ReceiverId   string `json:"receiverId" db:"receiver_id"`
	ReceiverName string `json:"receiverName" db:"receiver_name"`
}

func (self *GMessage) Insert() (string, error) {
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
	GMessage
}

func (self *GMessage) GetByID(query *GetQuery) (*GetModel, error) {
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
	GMessage
	Avatar string `json:"avatar"`
}

func (self *GMessage) List(query *Query) ([]*ListModel, int64, error) {
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

func (self *GMessage) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	ID     string `json:"id" db:"id"`
	UserId string
	Name   string              `db:"name"`
	Media  []*mediaModel.Media `json:"media" db:"-"`

	Disabled   bool  `json:"disabled" db:"disabled"`
	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GMessage) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareNamed(updateSql())
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	query.Updatetime = util.GetCurrentMS()
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}

	// 先把媒体文件插入数据库
	err = mediaService.Del(&mediaModel.DeleteQuery{
		Businesses:  []string{BusinessName},
		BusinessIds: []string{query.ID},
	}, tx)
	if err != nil {
		return err
	}
	medias := mediaService.InitMedias(query.Media, BusinessName, query.ID, query.UserId)
	err = mediaService.MultiCreate(medias, tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

type DeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
}

// 删除，批量删除
func (self *GMessage) Delete(query *DeleteQuery) error {
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

func (self *GMessage) ToggleDisabled(query *DisabledQuery) error {
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
