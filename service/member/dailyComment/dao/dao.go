package dao

import (
	"babygrow/DB/appGorm"
	"babygrow/service/member/dailyComment/dbModel"
	"babygrow/util/interfaces"
	"fmt"

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
	if query.GetStringValue("userId") != "" {
		tx.Where("g_member_baby_grow_comment.user_id=?", query.GetStringValue("userId"))
	}
	if query.GetStringValue("diaryId") != "" {
		tx.Where("g_member_baby_grow_comment.diary_id=?", query.GetStringValue("diaryId"))
	}
	if len(query.ToStringArray("diaryIds")) > 0 {
		tx.Where("g_member_baby_grow_comment.diary_id=any(?)", query.ToStringArray("diaryIds"))
	}
	if query.GetStringValue("babyId") != "" {
		tx.Where("g_member_baby_grow_comment.baby_id=?", query.GetStringValue("babyId"))
	}
	tx.Scopes(appGorm.BaseWhere2(query, ""))
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var list = make([]map[string]interface{}, 0)
	err = tx.Scopes(appGorm.Paginate2(query, "g_member_baby_grow_comment")).Find(&list).Error
	nlist := interfaces.NewModelMapSliceFromMapSlice(list)
	return nlist.ToCamelKey(), count, err
}

func Count(db *gorm.DB, query interfaces.QueryInterface) (map[string]int64, error) {
	tx := db.Model(&dbModel.GDailyComment{})
	m := make(map[string]interface{})
	if diaryId := query.GetStringValue("diaryId"); diaryId != "" {
		if strings.Contains(diaryId, " ") {
			return nil, fmt.Errorf("非法操作")
		}
		tx.Select("count(*) as a"+diaryId).Where("g_member_baby_grow_comment.diary_id=?", diaryId)
	} else if dids := query.ToStringArray("diaryIds"); len(dids) > 0 {
		for _, v := range dids {
			if strings.Contains(v, " ") {
				return nil, fmt.Errorf("非法操作")
			}
			tx.Select("(?) as a"+v, db.Model(&dbModel.GDailyComment{}).Select("count(*)").Where("diary_id=?", v))
		}
	} else if commentIds := query.ToStringArray("commentIds"); len(commentIds) > 0 {
		sqlStr := ""
		sqlRaws := make([]interface{}, len(commentIds))
		for k, v := range commentIds {
			if strings.Contains(v, " ") {
				return nil, fmt.Errorf("非法操作")
			}
			sqlStr = sqlStr + ",(?) as a" + v
			// 计算数量时需要减掉自己
			sqlRaws[k] = db.Raw("select (case when count(*)>0 then count(*)-1 else 0 end) from (with recursive t as (select * from g_member_baby_grow_comment g where g.delete_at is null and g.id=? union all select k.* from g_member_baby_grow_comment k,t where t.id = k.comment_id) select * from t) as at", v)
			// postgres=# with recursive t as(select id,name,parentid from tb9 where id=1 union all select k.id,k.name,k.parentid from tb9 k,t where t.id=k.parentid) select * from t;
		}
		tx.Select(sqlStr[1:], sqlRaws...)
	}
	err := tx.Find(&m).Error
	if err != nil {
		return nil, err
	}
	mi := make(map[string]int64)
	for k, v := range m {
		mi[k[1:]] = v.(int64)
	}
	return mi, err
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func Update(db *gorm.DB, query interfaces.QueryInterface) error {
	err := db.Model(&dbModel.GDailyComment{}).Select("content").Where("id=?", query.GetID()).Updates(map[string]interface{}{
		"content": query.GetValueWithDefault("content", ""),
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
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.GDailyComment{}).Error
}
