package model

import (
	"time"
)

// TableName sets Message's table name to be `messages`
func (Message) TableName() string {
	return "messages"
}

type Message struct {
	MessagePublic
	MessagePrivate
}

type MessagePublic struct {
	Message         string   `json:"message" binding:"required"`
	Subject         *string  `json:"subject"`
	RecipientId     *string  `json:"recipientId"`
	Recipient       *User    `gorm:"-" json:"recipient"`
	Parent          *Message `gorm:"foreignkey:ParentId;association_foreignkey:ID" json:"parent"`
	ParentId        *uint64  `json:"parentId"`
	DeleteAfterRead *bool    `json:"deleteAfterRead"`
}

type MessagePrivate struct {
	ID                   uint64     `gorm:"primary_key" json:"id"`
	SenderId             *string    `json:"senderId"`
	Sender               *User      `gorm:"-" json:"sender"`
	Edited               *bool      `json:"edited"`
	DeletedForSender     *bool      `json:"-"`
	DeletedForRecipient  *bool      `json:"-"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	IsSenderRead         bool       `json:"isSenderRead"`
	IsRecipientRead      bool       `json:"isRecipientRead"`
	Children             []*Message `gorm:"foreignkey:ParentId;association_foreignkey:ID" json:"children"`
	IsRecipientIncoming  bool       `json:"isRecipientIncoming"`
	LastMessageCreatedAt time.Time  `gorm:"-" json:"lastMessageCreatedAt"`
}

type MessageEditable struct {
	Message     string  `json:"message"`
	RecipientId *string `json:"recipientId"`
}
