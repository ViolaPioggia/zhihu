package model

import "time"

type Article struct {
	ID         int64     `form:"id" json:"id"`
	UserID     int64     `form:"userID" json:"userID"`
	TopicID    int64     `form:"topicID" json:"topicID"`
	Context    string    `form:"context" json:"context" binding:"required"`
	VIP        string    `form:"vip" json:"vip" `
	CreateTime time.Time `form:"createTime" json:"createTime"`
	UpdateTime time.Time `form:"updateTime" json:"updateTime"`
}
