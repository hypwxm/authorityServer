package model

import (
	"babygrow/DB/appGorm"
	"babygrow/DB/pgsql"
	"babygrow/util"
	"babygrow/util/database"

	memberModel "babygrow/service/member/user/model"
	memberService "babygrow/service/member/user/service"

	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type GFamilyMembers struct {
	database.BaseColumns
	// 会员id
	MemberId string `json:"memberId" db:"member_id" gorm:"column:member_id;type:varchar(128);check:member_id<>'';not null;uniqueIndex:index_family_member"`
	// 家园id
	FamilyId string `json:"familyId" db:"family_id" gorm:"column:family_id;type:varchar(128);check:family_id<>'';not null;uniqueIndex:index_family_member"`
	// 创建者，默认家园的创建者才能拉成员，
	Creator string `json:"creator" db:"creator" gorm:"column:creator;type:varchar(128);check:creator<>'';not null"`
	// 能否拉人
	CanInvite bool `json:"canInvite" db:"can_invite" gorm:"column:can_invite;type:bool;default false;not null"`
	// 能否删除人
	CanRemove bool `json:"canRemove" db:"can_remove" gorm:"column:can_remove;type:bool;default false;not null"`
	// 能否对家园信息进行编辑
	CanEdit bool `json:"canEdit" db:"can_edit" gorm:"column:can_edit;type:bool;default:false;not null"`
	// 在家园的昵称
	Nickname string `json:"nickname" db:"nickname" gorm:"column:nickname;type:varchar(20);default '';not null;"`
	// 在家庭中的角色
	RoleName string `json:"roleName" db:"role_name" gorm:"column:role_name;type:varchar(50);default '';not null"`
	// 角色类型，  1: 群主、管理员（默认创建者为管理员）,2: 管理员，由群主分配，3：成员
	RoleType int `json:"roleType" db:"role_type" gorm:"column:role_type;default 3;not null"`
}

func (self *GFamilyMembers) Insert(ctx context.Context) (string, error) {
	var err error

	if strings.TrimSpace(self.MemberId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.FamilyId) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.Creator) == "" {
		return "", fmt.Errorf("操作错误")
	}
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if !ok {
		db := appGorm.Open()
		tx = db.Begin()
		if err := tx.Error; err != nil {
			return "", err
		}
	}

	if !ok {
		defer tx.Rollback()
	}

	err = tx.Create(self).Error
	if err != nil {
		return "", err
	}

	if !ok {
		err = tx.Commit().Error
		if err != nil {
			return "", err
		}
	}

	return self.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

type GetModel struct {
	GFamilyMembers
}

func (self *GFamilyMembers) GetByID(query *GetQuery) (*GetModel, error) {
	db := pgsql.Open()
	stmt, err := db.PrepareNamed(getByIdSql())
	if err != nil {
		return nil, err
	}
	var entity = new(GetModel)
	err = stmt.Get(entity, query)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

type Query struct {
	appGorm.BaseQuery
	UserId   string `db:"user_id"`
	FamilyId string `db:"family_id"`
	Creator  string `db:"creator"`
	Keywords string `db:"keywords"`
}

type ListModel struct {
	GFamilyMembers
	MemberAccount  string `json:"memberAccount" gorm:"-"`
	MemberAvatar   string `json:"memberAvatar" gorm:"-"`
	MemberNickname string `json:"memberNickname" gorm:"-"`
}

func (self *GFamilyMembers) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	if query.UserId != "" && query.FamilyId == "" {
		return nil, 0, fmt.Errorf("参数错误")
	}
	db := appGorm.Open()
	// SELECT
	// g_member_family_member.*
	// FROM g_member_family_member WHERE 1=1  and g_member_family_member.isdelete='false'
	// and g_member_family_member.family_id=$1  order by g_member_family_member.createtime desc
	var count int64
	db = db.Model(&GFamilyMembers{}).Scopes(appGorm.BaseWhere(query.BaseQuery)).Where("family_id=?", query.FamilyId)
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	var list = make([]*ListModel, 0)
	err = db.Scopes(appGorm.Paginate(query.BaseQuery)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	ids := make([]string, len(list))
	for _, v := range list {
		ids = append(ids, v.MemberId)
	}
	// 获取成员的基本信息
	users, _, err := memberService.List(&memberModel.Query{
		BaseQuery: pgsql.BaseQuery{
			IDs: ids,
		},
	})
	if err != nil {
		return nil, 0, err
	}
	for _, v := range list {
		for _, vm := range users {
			if vm.ID == v.MemberId {
				v.MemberAccount = vm.Account
				v.MemberAvatar = vm.Avatar
				v.MemberNickname = vm.Nickname
				break
			}
		}
	}
	return list, count, nil

}

type UpdateByIDQuery struct {
	ID string `db:"id"`
	// 姓名
	Name string `json:"name" db:"name"`

	Updatetime int64 `db:"updatetime"`
}

// 更新,根据用户id和数据id进行更新
// 部分字段不允许更新，userID, id
func (self *GFamilyMembers) Update(query *UpdateByIDQuery) error {
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
func (self *GFamilyMembers) Delete(ctx context.Context, query *DeleteQuery) error {
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

	// 查询要删除的这些人的是不是家园的管理员，管理员要退出家园得先解散家园或者家园中只有自己一个人了
	var familyMembersInfo = make([]*GFamilyMembers, 0)
	err := db.Select(&familyMembersInfo, fmt.Sprintf("select * from %s where id=any($1) and isdelete=false", table_name), query.IDs)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, v := range familyMembersInfo {
		if v.RoleType == 1 {
			// 如果要移除的人中有群主（虽然可以赋予其他人踢人的权利，但是不能删群主）
			// 查询一下该家园中还有没有其他人
			var count int
			err = db.Get(&count, fmt.Sprintf("select count(*) from %s where family_id=$1 and isdelete=false", table_name), v.FamilyId)
			if err != nil {
				log.Println(err)
				return err
			}
			// 如果群中除了群主外还有其他成员，这种情况下，群主还不能被移除
			if count > 1 {
				return fmt.Errorf("无法移除创建者")
			}
		}
	}

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
func (self *GFamilyMembers) ToggleDisabled(query *DisabledQuery) error {
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
