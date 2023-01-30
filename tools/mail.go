package tools

import (
	"log"
	"net/smtp"

	"github.com/jordan-wright/email"
)

func SendEm(zhuti string, data []byte) {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "1820029262@qq.com"
	// 设置接收方的邮箱
	e.To = []string{"1820029262@qq.com", "2647565400@qq.com"}
	//设置主题
	e.Subject = zhuti
	//设置文件发送的内容
	e.Text = data
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "1820029262@qq.com", "binirjcvttoedhah", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}
