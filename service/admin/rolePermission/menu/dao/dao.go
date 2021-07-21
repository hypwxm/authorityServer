package dao

import (
	"authorityServer/DB/appGorm"
	"authorityServer/service/admin/rolePermission/menu/dbModel"
	"authorityServer/util/interfaces"
	"fmt"

	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity *dbModel.GRoleMenu) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

func MultiInsert(db *gorm.DB, entity []*dbModel.GRoleMenu) ([]*dbModel.GRoleMenu, error) {
	err := db.Create(&entity).Error
	return entity, err
}

func Get(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	var entity = make(map[string]interface{})
	tx := db.Model(&dbModel.GRoleMenu{})
	if query.GetID() != "" {
		tx.Where("id=?", query.GetID())
	}
	err := tx.Find(&entity).Error
	mMap := interfaces.NewModelMapFromMap(entity)
	return mMap.ToCamelKey(), err
}

func List(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelMapSlice, error) {
	tx := db.Model(&dbModel.GRoleMenu{})
	tx.Select(`
			g_authority_role_menu.createtime,
			g_authority_role_menu.updatetime,
			g_authority_role_menu.role_id,
			g_authority_role_menu.menu_id,
			g_authority_menu.parent_id,
			g_authority_menu.name,
			g_authority_menu.path,
			g_authority_menu.icon
`)
	tx.Joins("inner join g_authority_menu on g_authority_role_menu.menu_id=g_authority_role_menu.id")
	tx.Where("g_authority_role_menu.role_id=any(?)", query.ToStringArray("roleIds"))

	tx.Scopes(appGorm.BaseWhere2(query, ""))
	var list = make([]map[string]interface{}, 0)
	err := tx.Scopes(appGorm.Paginate2(query, "")).Find(&list).Error
	nlist := interfaces.NewModelMapSliceFromMapSlice(list)
	return nlist.ToCamelKey(), err
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func Update(db *gorm.DB, query interfaces.QueryInterface) error {
	err := db.Model(&dbModel.GRoleMenu{}).Select("name", "path", "icon").Where("id=?", query.GetID()).Updates(map[string]interface{}{
		"name": query.GetValueWithDefault("name", ""),
		"path": query.GetValueWithDefault("path", ""),
		"icon": query.GetValueWithDefault("icon", ""),
	}).Error
	return err
}

// 删除，批量删除
func Delete(db *gorm.DB, query interfaces.QueryInterface) error {
	if query.GetStringValue("roleId") == "" {
		return fmt.Errorf("操作错误")
	}
	return db.Where("role_id=? and menu_id=any(?)", query.GetStringValue("roleId"), query.ToStringArray("menuIds")).Delete(&dbModel.GRoleMenu{}).Error
}
