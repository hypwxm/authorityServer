package model

import (
	"babygrow/DB/appGorm"
	"babygrow/util/database"
	"context"
	"log"
)

const BusinessName = "g_message"

type GMessage struct {
	database.BaseColumns

	BusinessName string `json:"businessName" db:"business_name" gorm:"column:business_name;size:50"`
	BusinessId   string `json:"businessId" db:"business_id" gorm:"column:business_id;size:128"`

	Title   string `json:"title" db:"title" gorm:"column:title;size:100"`
	Content string `json:"content" db:"content" gorm:"column:content"`

	// 消息发送时间，消息不一定是实时发送
	Sendtime int64 `json:"sendtime" db:"sendtime" gorm:"column:sendtime"`
	// 是否已读
	IsRead bool `json:"isRead" db:"is_read" gorm:"column:is_read"`
	// 阅读时间
	ReadDuration float64 `json:"readDuration" db:"read_duration" gorm:"column:read_duration"`
	// 读到哪里了
	ReadPercent float64 `json:"readPercent" db:"read_percent" gorm:"column:read_percent"`

	// 发送人信息
	SenderId   string `json:"senderId" db:"sender_id" gorm:"column:sender_id;size:128"`
	SenderName string `json:"senderName" db:"sender_name" gorm:"column:sender_name;size:200"`

	// 接受人信息
	ReceiverId   string `json:"receiverId" db:"receiver_id" gorm:"column:receiver_id;size:128"`
	ReceiverName string `json:"receiverName" db:"receiver_name" gorm:"column:receiver_name;size:200"`
}

func (s *GMessage) GetID() string {
	return s.ID
}

func (s *GMessage) Insert(ctx context.Context) (string, error) {
	return insert(ctx, s)
}

func insert(ctx context.Context, entity *GMessage) (string, error) {
	db := appGorm.Open()
	if err := db.Model(&GMessage{}).Create(entity).Error; err != nil {
		log.Println(err)
		return "", err
	} else {
		return entity.GetID(), nil
	}
}
