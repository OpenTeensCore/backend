package services

import (
	"OpenTeens/controller"
	"OpenTeens/model"
	"OpenTeens/utils"
	"fmt"
	"time"
)

// UserSendEmailService 发送邮箱验证码业务逻辑（包括存入数据库）
func UserSendEmailService(email string, ipaddr string) any {
	emailCheck := utils.ToolValidateValue(email, utils.Email)
	serviceCheck := controller.FuncUserEmailServiceCheck(ipaddr, email) // TODO: Need to be changed

	if emailCheck && serviceCheck {
		code := utils.GenerateRandomNumber()
		controller.FuncUserEmailSend(email, code, ipaddr)
		return true
	} else if !emailCheck {
		fmt.Println("==================> Send Email Failed to", email, "where ip is", ipaddr, "because of email is not valid.")
		return "Email Pattern is not Valid."
	} else if !serviceCheck {
		fmt.Println("==================> Send Email Failed to", email, "where ip is", ipaddr, "because of service is not valid.")
		return "Service is not valid. Usually you send too much email."
	} else {
		return "Unknown Error."
	}
}

func UserVerifyEmailService(code string, email string) any {
	// 从数据库中获取code和email
	success := model.DBGetEmailVerifyByCodeAndEmail(code, email)
	if !success {
		return "Code is not valid."
	}
	return true
}

// UserLoginService 用户登录业务逻辑（包括存入数据库）
func UserLoginService(username string, password string) (bool, string) {
	//	先加密密码
	passwordHashed, _ := utils.HashPassword(password)

	// 然后首先用用户名密码DBUserCheckAccountByUsernameAndPasswordHashed
	ok, user := model.DBUserCheckAccount(username, passwordHashed, model.LoginByUsername)
	if ok {
		_, token := model.DBUserLoginCreateToken(user, "username")
		return true, token
	}

	// 然后用邮箱密码DBUserCheckAccountByEmailAndPasswordHashed
	ok, user = model.DBUserCheckAccount(username, passwordHashed, model.LoginByEmail)
	if ok {
		_, token := model.DBUserLoginCreateToken(user, "email")
		return true, token
	}

	//	然后用uid密码DBUserCheckAccountByUidAndPasswordHashed
	ok, user = model.DBUserCheckAccount(username, passwordHashed, model.LoginByUserID)
	if ok {
		_, token := model.DBUserLoginCreateToken(user, "uid")
		return true, token
	}
	return false, ""
}

// UserSiteMessageSendService 发送站内信业务逻辑
func UserSiteMessageSendService(senderID uint, content string, recipients []uint) bool {
	// 创建站内信
	message := model.SiteMessage{SenderID: senderID, Content: content, SentAt: time.Now()}
	model.DBUserSiteMessageCreate(&message)

	// 为每个接收者创建记录
	for _, receiverID := range recipients {
		recipient := model.SiteMessageRecipient{MessageID: message.ID, ReceiverID: receiverID}
		model.DBUserSiteMessageRecipientCreate(&recipient)
	}

	return true // 或者根据实际情况返回成功与否
}
