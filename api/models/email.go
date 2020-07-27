package models

import (
	"gopkg.in/gomail.v2"
	"strconv"
)

// 自定义发送邮箱
func InitDialer(host, user, pass string , port int) *gomail.Dialer {
	return gomail.NewDialer(host, port, user, pass)
}

func SendEmail(mailinfo map[string]string, msg *gomail.Message) error {
	port, _ := strconv.Atoi(mailinfo["port"])
	gd		:= InitDialer(mailinfo["server"], mailinfo["username"], mailinfo["password"], port)
	if e 	:= gd.DialAndSend(msg); e != nil {
		return e
	}

	return nil
}

// 生成消息体
func CreateMsg(mailFrom string, mailTo []string, subject string, body string) *gomail.Message{
	m := gomail.NewMessage()
	m.SetHeader("From","Monitor" + "<" + mailFrom + ">")
	m.SetHeader("To", mailTo...)  //发送给多个用户
	m.SetHeader("Subject", subject)  //设置邮件主题
	m.SetBody("text/html", body)  //设置邮件正文

	return m
}

// 生成带附件的消息体， 不支持非实时发送。
func CreateMsgWithAnnex(mailFrom string, mailTo []string,subject string, body string, annex string) *gomail.Message{
	m := gomail.NewMessage()
	m.SetHeader("From","Monitor" + "<" + mailFrom + ">")
	m.SetHeader("To", mailTo...)  //发送给多个用户
	m.SetHeader("Subject", subject)  //设置邮件主题
	m.SetBody("text/html", body)     //设置邮件正文
	m.Attach(annex)							// 设置附件

	return m
}
