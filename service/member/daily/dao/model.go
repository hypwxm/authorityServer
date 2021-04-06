package model

import (
	"babygrow/DB/appGorm"
	mediaModel "babygrow/service/media/model"
	mediaService "babygrow/service/media/service"
	dailyCommentModel "babygrow/service/member/dailyComment/model"
	dailyCommentService "babygrow/service/member/dailyComment/service"
	"babygrow/util"
	"babygrow/util/interfaces"

	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity interfaces.ModelMap) (string, error) {
	var err error
	err = db.Create(&entity).Error
	if err != nil {
		return "", err
	}
	return entity.GetID(), nil
}

type GetQuery struct {
	ID string `db:"id"`
}

func Get(db *gorm.DB, query interfaces.ModelMap) (*gDaily, error) {
	var entity = new(gDaily)
	tx := db.Model(&gDaily{})
	if query.GetID() != "" {
		tx.Where("id=?", query.GetID())
	}
	err := tx.Find(&entity).Error
	return entity, err
}

type Query struct {
	appGorm.BaseQuery
	Keywords string `db:"keywords"`
	Status   int    `db:"status"`
	UserId   string `db:"user_id"`
	BabyId   string `db:"baby_id"`
}

type ListModel struct {
	GDaily
	RoleName string                         `json:"userRoleName" db:"user_role_name" gorm:"column:user_role_name"`
	Account  string                         `json:"userAccount" db:"user_account" gorm:"user_account"`
	RealName string                         `json:"userRealName" db:"user_realname" gorm:"user_realname"`
	Nickname string                         `json:"userNickname" db:"user_nickname" gorm:"user_nickname"`
	Phone    string                         `json:"userPhone" db:"user_phone" gorm:"user_phone"`
	Comments []*dailyCommentModel.ListModel `json:"comments" gorm:"-"`

	MemberMedia []*mediaModel.Media `json:"xxx" gorm:"-"`
	Avatar      string              `json:"avatar" gorm:"-"`
}

func (self *GDaily) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	if query.UserId == "" {
		return nil, 0, fmt.Errorf("参数错误")
	}
	db := appGorm.Open()
	tx := db.Model(&GDaily{}).Select(`
				g_member_baby_grow.*,
				COALESCE(g_member_baby_relation.role_name, '') as user_role_name,
				COALESCE(g_member.realname, '') as user_realname,
				COALESCE(g_member.account, '') as user_account,
				COALESCE(g_member.phone, '') as user_phone,
				COALESCE(g_member.nickname, '') as user_nickname
				`)
	tx.Joins("left join g_member_baby_relation on g_member_baby_relation.baby_id=g_member_baby_grow.baby_id and g_member_baby_relation.user_id=g_member_baby_grow.user_id")
	tx.Joins("left join g_member on g_member_baby_grow.user_id=g_member.id")
	// tx.Where("g_member_baby_grow.user_id=?", query.UserId)
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
	err = tx.Scopes(appGorm.Paginate(query.BaseQuery)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	var ids []string = make([]string, 0)
	var userIds []string = make([]string, 0)
	for _, v := range list {
		ids = append(ids, v.ID)
		userIds = append(userIds, v.UserId)
	}
	userIds = util.ArrayStringDuplicateRemoval(userIds)
	// 查找对应的媒体信息
	mediaService.ListWithMedia(ids, BusinessName, list, "", "ID")
	err = mediaService.ListWithMedia(userIds, "member", list, "MemberMedia", "UserId")
	if err != nil {
		return nil, 0, err
	}

	for _, v := range list {
		if len(v.MemberMedia) > 0 {
			v.Avatar = v.MemberMedia[0].Url
		}
	}

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

	db := appGorm.Open()
	err := db.Model(&GDaily{}).Select("date", "weight", "height", "diary", "weather", "mood", "health", "temperature").Where("id=?", query.ID).Updates(map[string]interface{}{
		"date":        query.Date,
		"weight":      query.Weight,
		"height":      query.Height,
		"diary":       query.Diary,
		"weather":     query.Weather,
		"mood":        query.Mood,
		"health":      query.Health,
		"temperature": query.Temperature,
	}).Error
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

	db := appGorm.Open()
	return db.Where("id=any(?)", query.IDs).Delete(GDaily{}).Error
}
