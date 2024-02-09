package controller

import (
	"OpenTeens/model"
	"fmt"
)

// FuncUserEmailSend 发送邮件验证码，需要计入数据库
func FuncUserEmailSend(email string, code string, ipaddr string) {
	// 首先写入EmailVerification表中
	model.DBUserAddEmailCode(email, code, ipaddr)
	// 打印日志
	fmt.Println("==================> Send Email Successfully to", email)
	fmt.Println("==================> Code is: ", code)
	return
}

// FuncUserEmailServiceCheck 主要检测当前IP在60秒内是否已经发送了2条验证码或一个邮箱发送了超过一条
func FuncUserEmailServiceCheck(ipaddr string, email string) bool {
	number := model.DBUserGetEmailCodeByIp(ipaddr, 3)
	//	查看是否大于2条
	if number >= 2 {
		return false
	}

	amount := model.DBUserGetEmailCodeByEmail(email, 1)
	if amount >= 1 {
		return false
	}

	return true
}
