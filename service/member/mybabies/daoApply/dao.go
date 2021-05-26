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
	apply.account as applyAccount,
	apply.nickname as applyNickname,
	apply.realname as applyRealname,
	invite.account as inviteAccount,
	invite.nickname as inviteNickname,
	invite.realname as inviteRealname,
	`)
	tx.Joins("left join g_member apply on g_member_baby_relation_apply.user_id=g_member.id")
	tx.Joins("left join g_member invite on g_member_baby_relation_apply.inviter_id=g_member.id")

	if userId := query.GetStringValue("userId"); userId != "" {
		tx.Where("user_id=?", userId)
	}

	if babyId := query.GetStringValue("babyId"); babyId != "" {
		tx.Where("baby_id=?", babyId)
	}

	if inviterId := query.GetStringValue("inviterId"); inviterId != "" {
		tx.Where("inviter_id=?", inviterId)
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
