package dao

import (
	"babygrow/DB/appGorm"
	"babygrow/service/member/dailyComment/dbModel"
	"babygrow/util/interfaces"

	"errors"
	"strings"

	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity *dbModel.GDailyComment) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

func Get(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	var entity = make(map[string]interface{})
	tx := db.Model(&dbModel.GDailyComment{})
	if query.GetID() != "" {
		tx.Where("g_member_baby_grow_comment.id=?", query.GetID())
	}
	err := tx.Find(&entity).Error
	mMap := interfaces.NewModelMapFromMap(entity)
	return mMap.ToCamelKey(), err
}

func List(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	tx := db.Model(&dbModel.GDailyComment{}).Select(`
	g_member_baby_grow_comment.*,
	COALESCE(g_member_baby_relation.role_name, '') as user_role_name,
	COALESCE(g_member.realname, '') as user_realname,
	COALESCE(g_member.account, '') as user_account,
	COALESCE(g_member.phone, '') as user_phone,
	COALESCE(g_member.nickname, '') as user_nickname`)
	tx.Joins("left join g_member_baby_relation on g_member_baby_relation.baby_id=g_member_baby_grow_comment.baby_id and g_member_baby_relation.user_id=g_member_baby_grow_comment.user_id")
	tx.Joins("left join g_member on g_member_baby_grow_comment.user_id=g_member.id")
	if query.GetValue("userId") != "" {
		tx.Where("g_member_baby_grow_comment.user_id=?", query.UserId)
	}
	if query.GetValue("diaryId") != "" {
		tx.Where("g_member_baby_grow_comment.diary_id=?", query.DiaryId)
	}
	if len(query.GetValue("diaryIds")) > 0 {
		tx.Where("g_member_baby_grow_comment.diary_id=any(?)", query.DiaryIds)
	}
	if query.GetValue("babyId") != "" {
		tx.Where("g_member_baby_grow_comment.baby_id=?", query.BabyId)
	}
	tx.Scopes(appGorm.BaseWhere2(query, ""))
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var list = make([]map[string]interface{}, 0)
	err = tx.Scopes(appGorm.Paginate2(query, "g_member_baby_grow")).Find(&list).Error
	nlist := interfaces.NewModelMapSliceFromMapSlice(list)
	return nlist.ToCamelKey(), count, err
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func Update(db *gorm.DB, query interfaces.QueryInterface) error {
	err := db.Model(&dbModel.GDaily{}).Select("date", "weight", "height", "diary", "weather", "mood", "health", "temperature").Where("id=?", query.GetID()).Updates(map[string]interface{}{
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
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.GDaily{}).Error
}
