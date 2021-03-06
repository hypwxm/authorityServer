package dao

import (
	"github.com/hypwxm/authorityServer/DB/appGorm"
	"github.com/hypwxm/authorityServer/service/admin/user/dbModel"
	"github.com/hypwxm/authorityServer/util/interfaces"

	"errors"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

const NeedPwdWords = "yes i need pwd"

func Insert(db *gorm.DB, entity *dbModel.GAdminUser) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

// 根据条件获取单个用户
func Get(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	var entity = make(map[string]interface{})
	tx := db.Model(&dbModel.GAdminUser{})
	if query.GetStringValue("needPwd") != NeedPwdWords {
		tx.Omit("password,salt")
	}
	if query.GetID() != "" {
		tx.Where("id=?", query.GetID())
	}
	if account := query.GetStringValue("account"); account != "" {
		tx.Where("account=?", account)
	}
	err := tx.Find(&entity).Error
	mMap := interfaces.NewModelMapFromMap(entity)
	return mMap.ToCamelKey(), err
}

func List(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	var tx *gorm.DB

	// 如果是以组织或者角色维度进行查询，要以"用户角色"表作为主表
	if strings.TrimSpace(query.GetStringValue("orgId")) != "" || len(query.ToStringArray("orgIds")) > 0 || len(query.ToStringArray("roleIds")) > 0 {
		tx = db.Model(&dbModel.GUserRole{})
		tx.Select(`
					g_authority_user_role.org_id,
					g_authority_user.id,
					g_authority_user.createtime,
					g_authority_user.updatetime,
					g_authority_user.account,
					g_authority_user.username,
					g_authority_user.post,
					g_authority_user.disabled,
					g_authority_user.sort,
					g_authority_user.creator_id,
					g_authority_user.creator,
					g_authority_user.contact_way
		`)
		tx.Joins("inner join g_authority_user on g_authority_user.id=g_authority_user_role.user_id")
	} else {
		tx = db.Model(&dbModel.GAdminUser{})
		tx.Select(`
		g_authority_user.id,
		g_authority_user.createtime,
		g_authority_user.updatetime,
		g_authority_user.account,
		g_authority_user.username,
		g_authority_user.post,
		g_authority_user.disabled,
		g_authority_user.sort,
		g_authority_user.creator_id,
		g_authority_user.creator,
		g_authority_user.contact_way
`)
	}
	orgIds := query.ToStringArray("orgIds")
	roleIds := query.ToStringArray("roleIds")
	if keywords := query.GetStringValue("keywords"); keywords != "" {
		tx.Where("g_authority_user.account::text ilike '%" + keywords + "%' or g_authority_user.username::text ilike '%" + keywords + "%'")
	}
	if orgId := query.GetStringValue("orgId"); orgId != "" && len(orgIds) == 0 {
		tx.Where("g_authority_user_role.org_id=?", orgId)
	}
	if len(orgIds) > 0 {
		tx.Where("g_authority_user_role.org_id=any(?)", orgIds)
	}
	if len(roleIds) > 0 {
		tx.Where("g_authority_user_role.org_id=any(?)", roleIds)
	}

	tx.Scopes(appGorm.BaseWhere2(query, ""))
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
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.GetStringValue("id")) == "" {
		return errors.New("更新条件错误")
	}
	if err := db.Model(&dbModel.GAdminUser{}).Select("username", "post", "sort", "contact_way", "disabled").Where("id=?", query.GetID()).Updates(map[string]interface{}{
		"username":    query.GetValueWithDefault("username", ""),
		"post":        query.GetValueWithDefault("post", ""),
		"sort":        query.GetValueWithDefault("sort", 0),
		"contact_way": query.GetValueWithDefault("contact_way", ""),
		"disabled":    query.GetValueWithDefault("disabled", ""),
	}).Error; err != nil {
		return err
	}
	if pwd := query.GetStringValue("password"); pwd != "" {
		if err := db.Model(&dbModel.GAdminUser{}).Select("password").Where("id=?", query.GetID()).Updates(map[string]interface{}{
			"password": pwd,
		}).Error; err != nil {
			return err
		}
	}
	return nil
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
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.GAdminUser{}).Error
}

func RolesInsert(db *gorm.DB, roles []*dbModel.GUserRole) error {
	return db.Create(&roles).Error
}

/**
根据用户ids拿到对应角色列表
*/

func GetRolesByUserIds(db *gorm.DB, ids pq.StringArray) (interfaces.ModelMapSlice, error) {
	tx := db.Model(&dbModel.GUserRole{})
	tx.Select("g_authority_user_role.*,g_authority_role.*")
	tx.Joins("inner join g_authority_role on g_authority_user_role.role_id=g_authority_role.id and g_authority_role.delete_at is null")
	tx.Where("g_authority_user_role.user_id=any(?)", ids)
	var list = make([]map[string]interface{}, 0)
	err := tx.Find(&list).Error
	nlist := interfaces.NewModelMapSliceFromMapSlice(list)
	return nlist.ToCamelKey(), err
}

func DeleteRoles(db *gorm.DB, query interfaces.QueryInterface) error {
	return db.Where("user_id=?", query.GetID()).Unscoped().Delete(&dbModel.GUserRole{}).Error
}
