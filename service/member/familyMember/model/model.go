package model

import (
	"babygrow/DB/pgsql"
	"babygrow/event"
	"babygrow/util"
	"babygrow/util/database"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type GFamilyMembers struct {
	database.BaseColumns
	// 会员id
	MemberId string `json:"memberId" db:"member_id"`
	// 家园id
	FamilyId string `json:"familyId" db:"family_id"`
	// 创建者，默认家园的创建者才能拉成员，
	Creator string `json:"creator" db:"creator"`
	// 能否拉人
	CanInvite bool `json:"canInvite" db:"can_invite"`
	// 能否删除人
	CanRemove bool `json:"canRemove" db:"can_remove"`
	// 能否对家园信息进行编辑
	CanEdit bool `json:"canEdit" db:"can_edit"`
	// 在家园的昵称
	Nickname string `json:"nickname" db:"nickname"`
	// 在家庭中的角色
	RoleName string `json:"roleName" db:"role_name"`
	// 角色类型，  1: 群主、管理员（默认创建者为管理员）,2: 管理员，由群主分配，3：成员
	RoleType int `json:"roleType" db:"role_type"`
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
	tx, ok := ctx.Value("tx").(*sqlx.Tx)
	if !ok {
		db := pgsql.Open()
		tx, err = db.Beginx()
		if err != nil {
			return "", err
		}
	}

	if !ok {
		defer tx.Rollback()
	}
	// 插入判断用户登录账号是否已经存在
	stmt, err := tx.PrepareNamed(insertSql())
	if err != nil {
		return "", err
	}
	log.Println(stmt.QueryString)
	var lastId string
	self.BaseColumns.Init()
	err = stmt.Get(&lastId, self)
	if err != nil {
		return "", err
	}

	if !ok {
		err = tx.Commit()
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
	pgsql.BaseQuery
	UserId   string `db:"user_id"`
	FamilyId string `db:"family_id"`
	Creator  string `db:"creator"`
	Keywords string `db:"keywords"`
}

type ListModel struct {
	GFamilyMembers
}

func (self *GFamilyMembers) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	if query.UserId != "" && query.FamilyId == "" {
		return nil, 0, fmt.Errorf("参数错误")
	}
	db := pgsql.Open()
	whereSql, fullSql := listSql(query)
	// 以上部分为查询条件，接下来是分页和排序
	count, err := self.GetCount(db, query, whereSql)
	if err != nil {
		return nil, 0, err
	}
	stmt, err := db.PrepareNamed(fullSql)
	if err != nil {
		return nil, 0, err
	}
	log.Println(stmt.QueryString)

	rows, err := stmt.Queryx(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users = make([]*ListModel, 0)
	for rows.Next() {
		var user = new(ListModel)
		err = rows.StructScan(&user)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, count, nil

}

func (self *GFamilyMembers) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
	err := db.Select(familyMembersInfo, fmt.Sprintf("select * from %s where id=any(:ids)", table_name), query)
	if err != nil {
		return err
	}

	for _, v := range familyMembersInfo {
		if v.RoleType == 1 {
			// 如果要移除的人中有群主（虽然可以赋予其他人踢人的权利，但是不能删群主）
			// 查询一下该家园中还有没有其他人
			var count int
			err = db.Select(&count, fmt.Sprintf("select count(*) from %s where family_id=%s", table_name, v.FamilyId))
			if err != nil {
				return err
			}
			// 如果群中除了群主外还有其他成员，这种情况下，群主还不能被移除
			if count > 1 {
				return fmt.Errorf("无法移除创建者")
			}
			// 只剩群主一个人了，群主要离开群了，（把群也删了）
			// 发布一个删除家园的事件
			var ch = make(chan int, 1)
			event.Ebus.Publish("memberSv:familyDelete", []string{v.FamilyId}, ch)
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
