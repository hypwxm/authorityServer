package dao

import (
	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/service/admin/role/dbModel"
	"github.com/hypwxm/authorityServer/util/interfaces"

	"errors"
	"strings"

	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity *dbModel.GRole) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

func Get(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	var entity = make(map[string]interface{})
	tx := db.Model(&dbModel.GRole{})
	if query.GetID() != "" {
		tx.Where("id=?", query.GetID())
	}
	err := tx.Find(&entity).Error
	mMap := interfaces.NewModelMapFromMap(entity)
	return mMap.ToCamelKey(), err
}

func List(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	tx := db.Model(&dbModel.GRole{})
	tx.Scopes(appGorm.BaseWhere2(query, ""))

	if keywords := query.GetStringValue("keywords"); keywords != "" {
		tx.Where("(g_authority_role.name like '%" + keywords + "%' or g_authority_role.intro like '%" + keywords + "%')")
	}
	if orgId := query.GetStringValue("orgId"); orgId != "" {
		tx.Where("g_authority_role.org_id=?", orgId)
	}

	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var list = make([]map[string]interface{}, 0)
	err = tx.Scopes(appGorm.Paginate2(query, "")).Find(&list).Error
	nlist := interfaces.NewModelMapSliceFromMapSlice(list)
	return nlist.ToCamelKey(), count, err
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func Update(db *gorm.DB, query interfaces.QueryInterface) error {
	err := db.Model(&dbModel.GRole{}).Select("name", "intro").Where("id=?", query.GetID()).Updates(map[string]interface{}{
		"name":  query.GetValueWithDefault("name", ""),
		"intro": query.GetValueWithDefault("intro", ""),
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
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.GRole{}).Error
}
