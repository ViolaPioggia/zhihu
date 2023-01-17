package model

import "time"

type Comment struct {
	ID         int64     `form:"id" json:"id"`
	UserID     int64     `form:"userID" json:"userID"`
	ArticleID  int64     `form:"articleID" json:"articleID"`
	Context    string    `form:"context" json:"context" binding:"required"`
	CreateTime time.Time `form:"createTime" json:"createTime"`
	UpdateTime time.Time `form:"updateTime" json:"updateTime"`
}
