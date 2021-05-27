package daoApply

import (
	"babygrow/DB/appGorm"
	"babygrow/service/member/mybabies/dbModel"
	"babygrow/util/interfaces"

	"errors"
	"strings"

	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity *dbModel.GMemberBabyRelationApply) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

func List(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	tx := db.Model(&dbModel.GMemberBabyRelationApply{}).Select(`g_member_baby_relation_apply.*,
	apply.account as apply_account,
	apply.nickname as apply_nickname,
	apply.realname as apply_realname,
	invite.account as invite_account,
	invite.nickname as invite_nickname,
	invite.realname as invite_realname,
	g_member_baby.name as baby_name
	`)
	tx.Joins("left join g_member apply on g_member_baby_relation_apply.user_id=apply.id")
	tx.Joins("left join g_member invite on g_member_baby_relation_apply.inviter_id=invite.id")
	tx.Joins("left join g_member_baby on g_member_baby_relation_apply.baby_id=g_member_baby.id")

	if userId := query.GetStringValue("userId"); userId != "" {
		tx.Where("g_member_baby_relation_apply.user_id=?", userId)
	}

	if babyId := query.GetStringValue("babyId"); babyId != "" {
		tx.Where("g_member_baby_relation_apply.baby_id=?", babyId)
	}

	if inviterId := query.GetStringValue("inviterId"); inviterId != "" {
		tx.Where("g_member_baby_relation_apply.inviter_id=?", inviterId)
	}
	tx.Scopes(appGorm.BaseWhere2(query, ""))
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var list = make([]map[string]interface{}, 0)
	err = tx.Scopes(appGorm.Paginate2(query, "g_member_baby_relation_apply")).Find(&list).Error
	nlist := interfaces.NewModelMapSliceFromMapSlice(list)
	return nlist.ToCamelKey(), count, err
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
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.GMemberBabyRelationApply{}).Error
}

func UpdateApplyStatus(db *gorm.DB, query interfaces.QueryInterface) error {
	return db.Model(&dbModel.GMemberBabyRelationApply{}).Update("status", query.GetValue("status")).Error
}

func Get(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	var entity = make(map[string]interface{})
	tx := db.Model(&dbModel.GMemberBabyRelationApply{})
	if id := query.GetID(); id != "" {
		tx.Where("id=?", id)
	}
	err := tx.Find(&entity).Error
	mMap := interfaces.NewModelMapFromMap(entity)
	return mMap.ToCamelKey(), err
}
