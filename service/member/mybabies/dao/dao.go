package dao

import (
	"babygrow/DB/appGorm"
	"babygrow/service/member/mybabies/dbModel"
	"babygrow/util/interfaces"

	"errors"
	"strings"

	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity *dbModel.GMyBabies) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

func Get(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	var entity = make(map[string]interface{})
	tx := db.Model(&dbModel.GMyBabies{})
	if id := query.GetID(); id != "" {
		tx.Where("id=?", id)
	}
	err := tx.Find(&entity).Error
	mMap := interfaces.NewModelMapFromMap(entity)
	return mMap.ToCamelKey(), err
}

func List(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	tx := db.Model(&dbModel.GMyBabies{})
	tx.Select(`
		g_member_baby_relation.role_name,
		g_member_baby_relation.user_id,
		g_member_baby.id,
		g_member_baby.createtime,
		g_member_baby.updatetime,
		g_member_baby.name, 
		g_member_baby.birthday, 
		g_member_baby.gender, 
		g_member_baby.avatar,
		g_member_baby.id_card, 
		g_member_baby.hobby,
		g_member_baby.good_at,
		g_member_baby.favorite_food, 
		g_member_baby.favorite_color, 
		g_member_baby.ambition
	`)
	tx.Joins("left join g_member_baby_relation on g_member_baby.id=g_member_baby_relation.baby_id")
	tx.Where("g_member_baby_relation.delete_at is null")
	if userId := query.GetStringValue("userId"); userId != "" {
		tx.Where("g_member_baby_relation.user_id=?", userId)
	}
	tx.Scopes(appGorm.BaseWhere2(query, ""))
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var list = make([]map[string]interface{}, 0)
	err = tx.Find(&list).Error
	nlist := interfaces.NewModelMapSliceFromMapSlice(list)
	return nlist.ToCamelKey(), count, err
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func Update(db *gorm.DB, query interfaces.QueryInterface) error {
	err := db.Model(&dbModel.GMyBabies{}).Select("date", "weight", "height", "diary", "weather", "mood", "health", "temperature").Where("id=?", query.GetID()).Updates(map[string]interface{}{
		"date":        query.GetValueWithDefault("date", ""),
		"weight":      query.GetValueWithDefault("weight", 0),
		"height":      query.GetValueWithDefault("height", 0),
		"diary":       query.GetValue("diary"),
		"weather":     query.GetValue("weather"),
		"mood":        query.GetValue("mood"),
		"health":      query.GetValue("health"),
		"temperature": query.GetValueWithDefault("temperature", 0),
	}).Error
	return err
}

// 删除，批量删除
func Delete(db *gorm.DB, query interfaces.QueryInterface) error {
	if len(query.GetIDs()) == 0 {
		return errors.New("操作条件错误")
	}
	for _, v := range query.GetIDs() {
		if strings.TrimSpace(v) == "" {
			return errors.New("操作条件错误")
		}
	}
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.GMyBabies{}).Error
}
