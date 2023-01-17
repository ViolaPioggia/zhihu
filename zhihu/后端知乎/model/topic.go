package model

import "time"

type Topic struct {
	ID         int64     `form:"id" json:"id"`
	UserID     int64     `form:"userID" json:"userID"`
	ArticleID  int64     `form:"articleID" json:"articleID"`
	Title      string    `form:"title" json:"title" binding:"required"`
	Context    string    `form:"context" json:"context" binding:"required"`
	Tags       string    `form:"tags" json:"tags"`
	CreateTime time.Time `form:"createTime" json:"createTime"`
	UpdateTime time.Time `form:"updateTime" json:"updateTime"`
}
