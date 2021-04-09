package dao

import (
	"babygrow/DB/appGorm"
	"babygrow/service/media/dbModel"
	"babygrow/util/interfaces"

	"errors"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func Insert(db *gorm.DB, entity *dbModel.Media) (string, error) {
	err := db.Create(&entity).Error
	return entity.ID, err
}

func Get(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelInterface, error) {
	var entity = make(map[string]interface{})
	tx := db.Model(&dbModel.Media{})
	if query.GetID() != "" {
		tx.Where("id=?", query.GetID())
	}
	err := tx.Find(&entity).Error
	mMap := interfaces.NewModelMapFromMap(entity)
	return mMap.ToCamelKey(), err
}

func List(db *gorm.DB, query interfaces.QueryInterface) (interfaces.ModelMapSlice, int64, error) {
	tx := db.Model(&dbModel.Media{})
	// tx.Where("g_member_baby_grow.user_id=?", query.UserId)

	if businesses, ok := query.GetValue("businesses").(pq.StringArray); ok {
		tx.Where("business=any(?)", businesses)
	}
	if businessIds, ok := query.GetValue("businessIds").(pq.StringArray); ok {
		tx.Where("business_id=any(?)", businessIds)
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
	err := db.Model(&dbModel.Media{}).Select("size", "url", "user_id").Where("id=?", query.GetID()).Updates(map[string]interface{}{
		"size":    query.GetValueWithDefault("size", 0),
		"url":     query.GetValueWithDefault("url", ""),
		"user_id": query.GetValueWithDefault("userId", ""),
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
	return db.Where("id=any(?)", query.GetIDs()).Delete(&dbModel.Media{}).Error
}
