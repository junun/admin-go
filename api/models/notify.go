package models

import (
	"api/pkg/logging"
	"time"
)

// 通知表
type Notify struct {
	Model
	Title        string
	Type      	 int
	Source       int
	Unread  	 int
	Content      string
	Link		 string
	CreateTime   time.Time
}

// 通用生成通知记录方法
func MakeNotify(mytype, source int, title, content, link string)  {
	var data Notify
	data.Type 		= mytype
	data.Source		= source
	data.Unread 	= 1
	data.Title		= title
	data.Content	= content
	data.CreateTime = time.Now().AddDate(0,0,0)
	data.Link		= link

	e := DB.Save(&data).Error

	if e != nil {
		logging.Error("Add Notify Failed.  Error :" + e.Error())
	}
}