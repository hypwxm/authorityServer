package model

import (
	"babygrow/DB/appGorm"
	"babygrow/DB/pgsql"
	"babygrow/util"
	"errors"
	"fmt"
	"log"
	"strings"
	"unsafe"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type GMyBabies struct {
	appGorm.BaseColumns

	// 操作的用户
	UserID string `json:"userId" db:"user_id" gorm:"column:user_id;type:varchar(128);not null;check(user_id<>'');index"`
	// 姓名
	Name string `json:"name" db:"name" gorm:"column:name;type:varchar(50);not null;default ''"`
	// 生日,(公历生日)
	Birthday string `json:"birthday" db:"birthday" gorm:"column:birthday;type:varchar(50);not null;default ''"`
	// 性别
	Gender string `json:"gender" db:"gender" gorm:"column:gender;type:varchar(10);not null;default ''"`
	// 照片
	Avatar string `json:"avatar" db:"avatar" gorm:"column:avatar;type:varchar(500);not null;default ''"`
	// 身份证号
	IdCard string `json:"idCard" db:"id_card" gorm:"column:id_card;type:varchar(50);not null;default ''"`
	// 兴趣
	Hobby string `json:"hobby" db:"hobby" gorm:"column:hobby;type:text;not null;default ''"`
	// 特长
	GoodAt string `json:"goodAt" db:"good_at" gorm:"column:good_at;type:text;not null;default ''"`
	// 喜欢的食物
	FavoriteFood string `json:"favoriteFood" db:"favorite_food" gorm:"column:favorite_food;type:text;not null;default ''"`
	// 喜欢的颜色
	FavoriteColor string `json:"favoriteColor" db:"favorite_color" gorm:"column:favorite_color;type:text;not null;default ''"`
	// 志向
	Ambition string `json:"ambition" db:"ambition" gorm:"column:ambition;type:text;not null;default ''"`

	Weight float64 `json:"weight" db:"weight" gorm:"column:weight;not null;default 0"`
	Height float64 `json:"height" db:"height" gorm:"column:height;not null;default 0"`

	// 该字段新建表时不创建，主要是传给创建关系用的
	RoleName string `json:"roleName" db:"role_name" gorm:"-"`
}

func (self *GMyBabies) Insert() (string, error) {
	var err error

	if strings.TrimSpace(self.Name) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Birthday) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.Gender) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}
	if strings.TrimSpace(self.UserID) == "" {
		return "", errors.New(fmt.Sprintf("操作错误"))
	}

	db := appGorm.Open()
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return "", err
	}
	defer tx.Rollback()
	// 插入判断用户登录账号是否已经存在
	err = tx.Model(&GMyBabies{}).Create(self).Error
	if err != nil {
		return "", err
	}

	memberBaby := new(GMemberBabyRelation)
	memberBaby.RoleName = self.RoleName
	memberBaby.BabyId = self.ID
	memberBaby.UserId = self.UserID
	_, err = memberBaby.Insert(tx)
	if err != nil {
		return "", err
	}

	err = tx.Commit().Error
	if err != nil {
		return "", err
	}

	return self.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

type GetModel struct {
	GMyBabies
	Relations []*GMemberBabyRelation `json:"relations" gorm:"-"`
}

func (self *GMyBabies) Get(query *GetQuery) (*GetModel, error) {
	db := appGorm.Open()
	tx := db.Model(&GMyBabies{})
	var entity = new(GetModel)
	err := tx.Where("id=?", query.ID).Find(&entity).Error
	if err != nil {
		return nil, err
	}
	return entity, nil
}

type Query struct {
	appGorm.BaseQuery

	UserId   string `db:"user_id"`
	Keywords string `db:"keywords"`
}

type ListModel struct {
	GMyBabies
	RoleName string `json:"roleName" db:"role_name" gorm:"column:role_name"`
}

func (self *GMyBabies) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := appGorm.Open()
	// 以上部分为查询条件，接下来是分页和排序
	// count, err := self.GetCount(db, query, whereSql)
	// if err != nil {
	// 	return nil, 0, err
	// }
	tx := db.Model(&GMyBabies{})
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
	if query.UserId != "" {
		tx.Where("g_member_baby_relation.user_id=?", query.UserId)
	}

	var list = make([]*ListModel, 0)
	err := tx.Find(&list).Error
	if err != nil {
		return nil, 0, err
	}
	count := len(list)
	return list, *(*int64)(unsafe.Pointer(&count)), nil

}

func (self *GMyBabies) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
	if query == nil {
		query = new(Query)
	}
	sqlStr := countSql(whereSql...)
	stmt, err := db.PrepareNamed(sqlStr)
	if err != nil {
		return 0, err
	}
	var count int64
	err = stmt.Get(&count, query)
	log.Println(stmt.QueryString, query)
	return count, err
}

type UpdateByIDQuery struct {
	ID string `db:"id"`
	// 姓名
	Name string `json:"name" db:"name"`
	// 生日,(公历生日)
	Birthday string `json:"birthday" db:"birthday"`
	// 性别
	Gender string `json:"gender" db:"gender"`
	// 照片
	Avatar string `json:"avatar" db:"avatar"`
	// 身份证号
	IdCard string `json:"idCard" db:"id_card"`
	// 兴趣
	Hobby string `json:"hobby" db:"hobby"`
	// 特长
	GoodAt string `json:"goodAt" db:"good_at"`
	// 喜欢的食物
	FavoriteFood string `json:"favoriteFood" db:"favorite_food"`
	// 喜欢的颜色
	FavoriteColor string `json:"favoriteColor" db:"favorite_color"`
	// 志向
	Ambition string `json:"ambition" db:"ambition"`

	Updatetime int64 `db:"updatetime"`

	Weight float64 `db:"weight"`

	Height float64 `db:"height"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GMyBabies) Update(query *UpdateByIDQuery) error {
	if query == nil {
		return errors.New("无更新条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("更新条件错误")
	}

	db := pgsql.Open()
	stmt, err := db.PrepareNamed(updateSql())
	if err != nil {
		return err
	}
	log.Println(stmt.QueryString)
	query.Updatetime = util.GetCurrentMS()
	_, err = stmt.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

type DeleteQuery struct {
	IDs pq.StringArray `db:"ids"`
}

// 删除，批量删除
func (self *GMyBabies) Delete(query *DeleteQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if len(query.IDs) == 0 {
		return errors.New("操作条件错误")
	}
	for _, v := range query.IDs {
		if strings.TrimSpace(v) == "" {
			return errors.New("操作条件错误")
		}
	}

	db := pgsql.Open()
	stmt, err := db.PrepareNamed(delSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}

type DisabledQuery struct {
	Disabled bool   `db:"disabled"`
	ID       string `db:"id"`
}

// 启用禁用店铺
func (self *GMyBabies) ToggleDisabled(query *DisabledQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.ID) == "" {
		return errors.New("操作条件错误")
	}
	db := pgsql.Open()
	stmt, err := db.PrepareNamed(toggleSql())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(query)
	return err
}
