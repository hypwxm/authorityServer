package model

import (
	"babygrow/DB/pgsql"
	mediaModel "babygrow/service/media/model"
	mediaService "babygrow/service/media/service"
	dailyCommentModel "babygrow/service/member/dailyComment/model"
	dailyCommentService "babygrow/service/member/dailyComment/service"

	"babygrow/util"
	"babygrow/util/database"

	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

const BusinessName = "g_member_baby_grow"

type GDaily struct {
	appGorm.BaseColumns

	// 今日份体重
	Weight float64 `json:"weight" db:"weight" gorm:"column:weight;not null;default 0"`
	// 今日份身高
	Height float64 `json:"height" db:"height" gorm:"column:height;not null;default 0"`

	// 今日份记录
	Diary string `json:"diary" db:"diary" gorm:"column:diary;type:text;not null;default ''"`

	UserId string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '')"`
	BabyId string `json:"babyId" db:"baby_id" gorm:"column:baby_id;type:varchar(128);not null;check(baby_id <> '')"`

	Weather     string  `json:"weather" db:"weather" gorm:"column:weather;type:varchar(50);not null;default ''"`
	Mood        string  `json:"mood" db:"mood" gorm:"column:mood;type:varchar(40);not null;default ''"`
	Health      string  `json:"health" db:"health" gorm:"column:mood;type:varchar(40);not null;default ''"`
	Temperature float64 `json:"temperature" db:"temperature" gorm:"column:temperature;not null;default 0"`

	Date string `json:"date" db:"date" gorm:"column:date;type:varchar(40);not null;default ''"`

	Sort int `json:"sort" db:"sort" gorm:"column:sort;not null;default 0"`

	Medias []*mediaModel.Media `json:"medias" gorm:"-"`
}

func (self *GDaily) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.UserId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.BabyId) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}

	db := appGorm.Open()
	tx := db.Begin()
	if err == tx.Error; err != nil {
		return "", err
	}
	defer tx.Rollback()
	err = tx.Create(&self).Error
	if err != nil {
		return "", err
	}

	medias := mediaService.InitMedias(self.Medias, BusinessName, self.ID, self.UserId)
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
	db := appGorm.Open()
	var entity = new(GetModel)
	err := db.Model(&GDaily{}).Where("id=?", query.ID).Find(&entity).Error
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
	BabyId   string `db:"baby_id"`
}

type ListModel struct {
	GDaily
	RoleName string                         `json:"userRoleName" db:"user_role_name"`
	Account  string                         `json:"userAccount" db:"user_account"`
	RealName string                         `json:"userRealName" db:"user_realname"`
	Nickname string                         `json:"userNickname" db:"user_nickname"`
	Phone    string                         `json:"userPhone" db:"user_phone"`
	Comments []*dailyCommentModel.ListModel `json:"comments"`
}

func (self *GDaily) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	if query.UserId == "" {
		return nil, 0, fmt.Errorf("参数错误")
	}
	db := appGorm.Open()
	tx := db.Select(`
				SELECT 
				g_member_baby_grow.*,
				COALESCE(g_member_baby_relation.role_name, '') as user_role_name,
				COALESCE(g_member.realname, '') as user_realname,
				COALESCE(g_member.account, '') as user_account,
				COALESCE(g_member.phone, '') as user_phone,
				COALESCE(g_member.nickname, '') as user_nickname
				`)
	tx.Joins("left join g_member_baby_relation on g_member_baby_relation.baby_id=g_member_baby_grow.baby_id and g_member_baby_relation.user_id=g_member_baby_grow.user_id")
	tx.Joins("left join g_member on g_member_baby_grow.user_id=g_member.id")
	tx.Where("g_member_baby_grow.user_id=?", query.UserId)
	if query.BabyId != "" {
		tx.Where("g_member_baby_grow.baby_id=?", query.BabyId)
	}
	tx.Scopes(appGorm.BaseWhere(query.BaseQuery))
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	var list = make([]*ListModel, 0)
	err = tx.Scopes(appGorm.Pagination()).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	var ids []string = make([]string, 0)
	var userIds []string = make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.ID)
		userIds = append(userIds, v.UserId)
	}

	// 查找对应的媒体信息
	mediaService.ListWithMedia(ids, BusinessName, list, "")

	// 获取评价内容
	if comments, _, err := dailyCommentService.List(&dailyCommentModel.Query{
		DiaryIds: ids,
	}); err != nil {
		return nil, 0, err
	} else {
		for _, v := range list {
			for _, vm := range comments {
				if v.ID == vm.DiaryId {
					v.Comments = append(v.Comments, vm)
				}
			}
		}
	}

	return list, count, nil

}

type UpdateByIDQuery struct {
	ID     string  `db:"id"`
	Date   string  `db:"date"`
	Weight float64 `json:"weight" db:"weight"`
	// 今日份身高
	Height float64 `json:"height" db:"height"`

	// 今日份记录
	Diary string `json:"diary" db:"diary"`

	Weather     string  `json:"weather" db:"weather"`
	Mood        string  `json:"mood" db:"mood"`
	Health      string  `json:"health" db:"health"`
	Temperature float64 `json:"temperature" db:"temperature"`

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
