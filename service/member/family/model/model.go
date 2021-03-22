package model

import (
	"babygrow/DB/appGorm"
	"babygrow/event"
	familyMemberModel "babygrow/service/member/familyMember/model"
	familyMemberService "babygrow/service/member/familyMember/service"
	"log"

	mediaModel "babygrow/service/media/model"
	mediaService "babygrow/service/media/service"

	memberModel "babygrow/service/member/user/model"
	memberService "babygrow/service/member/user/service"
	messageModel "babygrow/service/message/model"
	"context"

	"babygrow/util"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type GFamily struct {
	appGorm.BaseColumns
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

	// 插入封面图
	medias := mediaService.InitMedias(self.Medias, BusinessName, self.ID, self.Creator)
	err = mediaService.MultiCreate(medias)
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
	CreatorInfo memberModel.GetByIdModel `json:"creatorInfo" gorm:"-"`
}

/**
 * 根据家园id获取家园信息
 *
 * 家园的创建者信息
 */
func (self *GFamily) GetByID(query *GetQuery) (*GetModel, error) {
	db := appGorm.Open()
	var entity = new(GetModel)
	err := db.Where("id=?", query.ID).Find(&entity).Error
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

	// 查询媒体信息
	err = mediaService.ListWithMedia([]string{query.ID}, BusinessName, []*GFamily{&entity.GFamily}, "Medias")
	if err != nil {
		log.Println(err)
	}
	return entity, nil
}

type Query struct {
	appGorm.BaseQuery
	UserId string `db:"user_id"`

	Creator  string `db:"creator"`
	Keywords string `db:"keywords"`
}

type ListModel struct {
	familyMemberModel.GFamilyMembers
	FamilyName        string              `json:"familyName" db:"family_name"`
	FamilyCreator     string              `json:"familyCreator" db:"family_creator"`
	FamilyCreatorName string              `json:"familyCreatorName" db:"family_creator_name"`
	FamilyLabel       string              `json:"familyLabel" db:"family_label"`
	FamilyIntro       string              `json:"familyIntro" db:"family_intro"`
	FamilyCreatetime  int                 `json:"familyCreatetime" db:"family_createtime"`
	Medias            []*mediaModel.Media `json:"medias" gorm:"-"`
}

func (self *GFamily) List(query *Query) ([]*ListModel, int64, error) {
	if query == nil {
		query = new(Query)
	}
	db := appGorm.Open()
	tx := db.Table("g_member_family").Select(`
				COALESCE(g_member_family_member.member_id, '') as member_id ,
				COALESCE(g_member_family_member.creator, '') as creator,
				COALESCE(g_member_family_member.can_invite, false) as can_invite,
				COALESCE(g_member_family_member.can_remove, false) as can_remove,
				COALESCE(g_member_family_member.can_edit, false) as can_edit,
				COALESCE(g_member_family_member.nickname, '') as nickname,
				COALESCE(g_member_family_member.role_name, '') as role_name,
				COALESCE(g_member_family_member.role_type, 0) as role_type,
				COALESCE(g_member_family_member.createtime, 0) as createtime,
				g_member_family.name as family_name,
				g_member_family.creator as family_creator,
				g_member_family_member.nickname as family_creator_name,
				g_member_family.label as family_label,
				g_member_family.intro as family_intro,
				g_member_family.createtime as family_createtime,
				g_member_family.id as id
	`).Joins("left join g_member_family_member on g_member_family.id=g_member_family_member.family_id").Joins("left join g_member on g_member_family.creator=g_member.id ")
	tx.Scopes(appGorm.BaseWhere(query.BaseQuery))
	tx.Where("(g_member_family.creator=? or g_member_family_member.member_id=?)", query.UserId, query.UserId)
	var count int64
	err := tx.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	var list = make([]*ListModel, 0)
	err = tx.Scopes(appGorm.Paginate(query.BaseQuery)).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取家园封面
	var ids = make([]string, len(list))
	for _, v := range list {
		ids = append(ids, v.ID)
	}
	err = mediaService.ListWithMedia(ids, BusinessName, list, "Medias")
	if err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

type UpdateByIDQuery struct {
	ID string `db:"id"`
	// 姓名
	Name string `json:"name" db:"name"`
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

	db := appGorm.Open()
	err := db.Model(&GFamily{}).Where("id=?", query.ID).Updates(map[string]interface{}{
		"name": query.Name,
	}).Error
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
	db := appGorm.Open()
	return db.Where("id=any(?)", query.IDs).Delete(GFamily{}).Error
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
	db := appGorm.Open()
	return db.Model(&GFamily{}).Where("id=?", query.ID).Update("disabled=?", query.Disabled).Error
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
