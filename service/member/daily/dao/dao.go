package dao

import (
	"babygrow/DB/appGorm"
	"babygrow/service/member/daily/dbModel"
	"babygrow/util/interfaces"

	"errors"
	"strings"

	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity *dbModel.GDaily) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

func Get(db *gorm.DB, query interfaces.QueryMap) (interfaces.ModelMap, error) {
	var entity = make(interfaces.ModelMap)
	tx := db.Model(&dbModel.GDaily{})
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

func List(db *gorm.DB, query interfaces.QueryMap) ([]interfaces.ModelMap, int64, error) {
	tx := db.Model(&dbModel.GDaily{}).Select(`
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
	if babyId, ok := query["babyId"].(string); ok {
		tx.Where("g_member_baby_grow.baby_id=?", babyId)
	}
	tx.Scopes(query.BaseWhere())
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	var list = make([]interfaces.ModelMap, 0)
	err = tx.Scopes(query.Paginate()).Find(&list).Error
	return list, count, err
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func Update(db *gorm.DB, query interfaces.QueryMap) error {
	err := db.Model(&dbModel.GDaily{}).Select("date", "weight", "height", "diary", "weather", "mood", "health", "temperature").Where("id=?", query.GetID()).Updates(map[string]interface{}{
		"date":        query["date"],
		"weight":      query["weight"],
		"height":      query["height"],
		"diary":       query["diary"],
		"weather":     query["weather"],
		"mood":        query["mood"],
		"health":      query["health"],
		"temperature": query["temperature"],
	}).Error
	return err
}

// 删除，批量删除
func Delete(db *gorm.DB, query interfaces.QueryMap) error {
	if len(query.GetIDs()) == 0 {
		return errors.New("操作条件错误")
	}
	for _, v := range query.GetIDs() {
		if strings.TrimSpace(v) == "" {
			return errors.New("操作条件错误")
		}
	}
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.GDaily{}).Error
}
