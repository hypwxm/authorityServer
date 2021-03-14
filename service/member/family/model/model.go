package model

import (
	"babygrow/DB/appGorm"
	"babygrow/DB/pgsql"
	"babygrow/event"
	mediaModel "babygrow/service/media/model"
	familyMemberModel "babygrow/service/member/familyMember/model"
	familyMemberService "babygrow/service/member/familyMember/service"

	memberModel "babygrow/service/member/user/model"
	memberService "babygrow/service/member/user/service"
	messageModel "babygrow/service/message/model"
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
	Name string `json:"name" db:"name" gorm:"column:name;type:varchar(20);not null"`
	// 存储头像
	Medias []*mediaModel.Media `json:"medias" gorm:"-"`
	// 家族相册
	Album []*mediaModel.Media `json:"album" gorm:"-"`

	// 家庭标签，直接存字符串  逗号隔开
	Label string `json:"label" db:"label" gorm:"column:label;type:varchar(250);not null;default '';"`

	// 一个简单的描述
	Intro string `json:"intro" db:"intro" gorm:"column:intro;type:text;not null;default ''"`

	Creator string `json:"creator" db:"creator" gorm:"column:creator;type:varchar(128);not null;check:creator<>'';index"`
}

const BusinessName = "g_family"

func (self *GFamily) Insert(ctx context.Context) (string, error) {
	var err error

	if strings.TrimSpace(self.Creator) == "" {
		return "", fmt.Errorf("操作错误")
	}
	if strings.TrimSpace(self.Name) == "" {
		return "", fmt.Errorf("操作错误")
	}

	db := appGorm.Open()
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return "", err
	}
	defer tx.Rollback()
	err = tx.Create(self).Error
	if err != nil {
		return "", err
	}

	ctxTx := context.WithValue(ctx, "tx", tx)

	// 创建家园要先把创建者加入到家园中，角色为超管
	_, err = familyMemberService.Create(ctxTx, &familyMemberModel.GFamilyMembers{
		MemberId:  self.Creator,
		FamilyId:  self.ID,
		Creator:   self.Creator,
		CanInvite: true,
		CanRemove: true,
		CanEdit:   true,
		RoleType:  1,
	})
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

type InviteQuery struct {
	FamilyId      string `json:"familyId" db:"family_id"`
	MemberId      string `json:"memberId" db:"member_id"`
	MemberAccount string `json:"memberAccount" db:"member_account"`
	Invitor       string `json:"invitor"`
	InvitorName   string
}

// 发起邀请
// 邀请需要发送邀请信息
// 用户收到信息，同意后才会进入家园
// 注：消息系统还没开发，所以先直接邀请成功
func (self *GFamily) SendInviteMessage(query *InviteQuery) error {
	if query == nil {
		return errors.New("无操作条件")
	}
	if strings.TrimSpace(query.FamilyId) == "" {
		return errors.New("操作条件错误")
	}
	if strings.TrimSpace(query.MemberId) == "" {
		return errors.New("操作条件错误")
	}
	if strings.TrimSpace(query.MemberAccount) == "" {
		return errors.New("操作条件错误")
	}
	if strings.TrimSpace(query.Invitor) == "" {
		return errors.New("操作条件错误")
	}

	// TODO 需要发送邀请信息，用户同意后才进入家园
	// 简单点先直接将成员拉近家园
	_, err := familyMemberService.Create(context.Background(), &familyMemberModel.GFamilyMembers{
		MemberId: query.MemberId,
		FamilyId: query.FamilyId,
		Creator:  query.Invitor,
	})
	if err != nil {
		return err
	}

	event.Ebus.Publish("serve:message", &messageModel.GMessage{
		BusinessId:   query.FamilyId,
		BusinessName: BusinessName,
		SenderId:     query.Invitor,
		SenderName:   query.InvitorName,
		Title:        "家园邀请信息",
		Content:      "邀请信息",
		Sendtime:     util.GetCurrentMS(),
		ReceiverId:   query.MemberId,
		ReceiverName: "",
	})

	return nil
}
