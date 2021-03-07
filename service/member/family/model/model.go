package model

import (
	"babygrow/DB/pgsql"
	mediaModel "babygrow/service/media/model"
	familyMemberModel "babygrow/service/member/familyMember/model"
	familyMemberService "babygrow/service/member/familyMember/service"

	memberModel "babygrow/service/member/user/model"
	memberService "babygrow/service/member/user/service"

	"context"

	"babygrow/util"
	"babygrow/util/database"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type GFamily struct {
	database.BaseColumns
	Name string `json:"name" db:"name"`
	// 存储头像
	Medias []*mediaModel.Media `json:"medias"`
	// 家族相册
	Album []*mediaModel.Media `json:"album"`

	// 家庭标签，直接存字符串  逗号隔开
	Label string `json:"label" db:"label"`

	// 一个简单的描述
	Intro string `json:"intro" db:"intro"`

	Creator string `json:"creator" db:"creator"`
}

func (self *GFamily) Insert(ctx context.Context) (string, error) {
	var err error

	if strings.TrimSpace(self.Creator) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.Name) == "" {
		return "", fmt.Errorf("操作错误")
	}

	db := pgsql.Open()
	tx, err := db.Beginx()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
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

	ctxTx := context.WithValue(ctx, "tx", tx)

	// 创建家园要先把创建者加入到家园中，角色为超管
	_, err = familyMemberService.Create(ctxTx, &familyMemberModel.GFamilyMembers{
		MemberId:  self.Creator,
		FamilyId:  lastId,
		Creator:   self.Creator,
		CanInvite: true,
		CanRemove: true,
		CanEdit:   true,
		RoleType:  1,
	})
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return self.ID, nil
}

type GetQuery struct {
	ID string `db:"id"`
}

type GetModel struct {
	GFamily
	CreatorInfo memberModel.GetByIdModel `json:"creatorInfo"`
}

/**
 * 根据家园id获取家园信息
 *
 * 家园的创建者信息
 */
func (self *GFamily) GetByID(query *GetQuery) (*GetModel, error) {
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
	// 获取创建者
	creatorInfo, err := memberService.GetUserById(context.Background(), &memberModel.GetQuery{
		ID: entity.Creator,
	})
	if err != nil {
		return nil, err
	}
	entity.CreatorInfo = *creatorInfo
	return entity, nil
}

type Query struct {
	pgsql.BaseQuery
	UserId string `db:"user_id"`

	Creator  string `db:"creator"`
	Keywords string `db:"keywords"`
}

type ListModel struct {
	familyMemberModel.GFamilyMembers
	FamilyName        string `json:"familyName" db:"family_name"`
	FamilyCreator     string `json:"familyCreator" db:"family_creator"`
	FamilyCreatorName string `json:"familyCreatorName" db:"family_creator_name"`
	FamilyLabel       string `json:"familyLabel" db:"family_label"`
	FamilyIntro       string `json:"familyIntro" db:"family_intro"`
	FamilyCreatetime  int    `json:"familyCreatetime" db:"family_createtime"`
}

func (self *GFamily) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := pgsql.Open()
	_, fullSql := listSql(query)
	// 以上部分为查询条件，接下来是分页和排序
	// count, err := self.GetCount(db, query, whereSql)
	// if err != nil {
	// 	return nil, 0, err
	// }
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

	return users, 0, nil

}

func (self *GFamily) GetCount(db *sqlx.DB, query *Query, whereSql ...string) (int64, error) {
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
func (self *GFamily) Update(query *UpdateByIDQuery) error {
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
func (self *GFamily) Delete(query *DeleteQuery) error {
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
func (self *GFamily) ToggleDisabled(query *DisabledQuery) error {
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
