package model

import (
	"babygrow/DB/appGorm"
	mediaDBModel "babygrow/service/media/dbModel"
	mediaService "babygrow/service/media/service2"

	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

const BusinessName = "g_member_baby_grow_comment"

type GDailyComment struct {
	appGorm.BaseColumns

	Content string `json:"content" db:"content" gorm:"column:content;type:text;default '';not null"`

	UserId    string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id <> '')"`
	BabyId    string `json:"babyId" db:"baby_id" gorm:"column:baby_id;type:varchar(128);not null;check(baby_id <> '')"`
	DiaryId   string `json:"diaryId" db:"diary_id" gorm:"column:diary_id;type:varchar(128);not null;check(diary_id <> '')"`
	CommentId string `json:"commentId" db:"comment_id" gorm:"column:comment_id;type:varchar(128);not null;default ''"`

	Sort int `json:"sort" db:"sort" gorm:"column:sort;not null;default 0"`

	Medias []*mediaDBModel.Media `json:"medias" gorm:"-"`
}

func (self *GDailyComment) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.UserId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.BabyId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.DiaryId) == "" {
		return "", fmt.Errorf("操作错误")
	}

	db := appGorm.Open()

	// 插入判断用户登录账号是否已经存在
	err = db.Create(&self).Error
	if err != nil {
		return "", err
	}

	medias := mediaService.InitMedias(self.Medias, BusinessName, self.ID, self.UserId)
	err = mediaService.MultiCreate(medias)
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return self.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

type GetModel struct {
	GDailyComment
}

func (self *GDailyComment) GetByID(query *GetQuery) (*GetModel, error) {
	db := appGorm.Open()
	var entity = new(GetModel)
	err := db.Where("id=?", query.ID).Find(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

type Query struct {
	appGorm.BaseQuery
	Keywords string         `db:"keywords"`
	Status   int            `db:"status"`
	UserId   string         `db:"user_id"`
	BabyId   string         `db:"baby_id"`
	DiaryId  string         `json:"diaryId" db:"diary_id"`
	DiaryIds pq.StringArray `json:"diaryIds" db:"diary_ids"`
}

type ListModel struct {
	GDailyComment
	RoleName string `json:"userRoleName" db:"user_role_name" gorm:"column:user_role_name"`
	Account  string `json:"userAccount" db:"user_account" gorm:"column:user_account"`
	RealName string `json:"userRealName" db:"user_realname" gorm:"column:user_realname"`
	Nickname string `json:"userNickname" db:"user_nickname" gorm:"column:user_nickname"`
	Phone    string `json:"userPhone" db:"user_phone" gorm:"column:user_phone"`

	MemberMedia []*mediaDBModel.Media `json:"-" gorm:"-"`
	Avatar      string                `json:"avatar" gorm:"-"`
}

func (self *GDailyComment) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}

	db := appGorm.Open()
	tx := db.Model(&GDailyComment{}).Select(`
	g_member_baby_grow_comment.*,
	COALESCE(g_member_baby_relation.role_name, '') as user_role_name,
	COALESCE(g_member.realname, '') as user_realname,
	COALESCE(g_member.account, '') as user_account,
	COALESCE(g_member.phone, '') as user_phone,
	COALESCE(g_member.nickname, '') as user_nickname`)
	tx.Joins("left join g_member_baby_relation on g_member_baby_relation.baby_id=g_member_baby_grow_comment.baby_id and g_member_baby_relation.user_id=g_member_baby_grow_comment.user_id")
	tx.Joins("left join g_member on g_member_baby_grow_comment.user_id=g_member.id")
	if query.UserId != "" {
		tx.Where("g_member_baby_grow_comment.user_id=?", query.UserId)
	}
	if query.DiaryId != "" {
		tx.Where("g_member_baby_grow_comment.diary_id=?", query.DiaryId)
	}
	if len(query.DiaryIds) > 0 {
		tx.Where("g_member_baby_grow_comment.diary_id=any(?)", query.DiaryIds)
	}
	if query.BabyId != "" {
		tx.Where("g_member_baby_grow_comment.baby_id=?", query.BabyId)
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

	// // 查找对应的媒体信息
	// err = mediaService.ListWithMedia(ids, BusinessName, list, "", "ID")
	// if err != nil {
	// 	log.Println("error=========", err)
	// }
	// err = mediaService.ListWithMedia(userIds, "member", list, "MemberMedia", "UserId")
	// if err != nil {
	// 	log.Println("error=========", err)
	// }
	// for _, v := range list {
	// 	if len(v.MemberMedia) > 0 {
	// 		v.Avatar = v.MemberMedia[0].Url
	// 	}
	// }

	return list, count, nil

}

type UpdateByIDQuery struct {
	ID      string `db:"id"`
	Content string `json:"diary" db:"diary"`

	Medias []*mediaDBModel.Media `json:"medias"`
	UserId string                `json:"userId"`
}
type DeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
}

// 删除，批量删除
func (self *GDailyComment) Delete(query *DeleteQuery) error {
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
	return db.Where("id=any(?)", query.IDs).Delete(&GDailyComment{}).Error
}
