package controller

import (
	"OpenTeens/model"
	"OpenTeens/services"
	"OpenTeens/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 代码命名格式：
// XXXHandler是会被Gin调用的函数，用于调用业务逻辑和返回数据
// XXXService是会被Handler调用的函数，用于实现业务逻辑
// FuncXXX是业务逻辑辅助模块，用于辅助实现业务逻辑
// 有的业务可能只有Handler函数，没有Service函数，但有的业务可能有很多Func函数

// IndexHandler 处理访问Root的请求
func IndexHandler(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "Hello, OpenTeens!", "data": gin.H{"token": utils.GenerateToken(), "verify": utils.GenerateRandomNumber()}})
}

// UserSendEmailHandler 发送邮箱验证码的处理
func UserSendEmailHandler(c *gin.Context) {
	// 逻辑链：
	// 调用Service函数，传入初始信息
	//		-> 检查邮箱合法和未注册 -> ToolUserEmailCheck（会调用数据库）
	//		-> 检查IP上限 -> FuncUserEmailServiceCheck(会调用数据库)
	// 		-> 发送邮件 -> FuncUserEmailSend（会调用数据库）
	// -> 返回错误信息

	err := services.UserSendEmailService(c.PostForm("email"), c.ClientIP())
	// 如果返回值是Error，说明发送失败，返回400，否则返回200
	if err != true {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	c.JSON(200, gin.H{"msg": "Code send successfully!"})
}

// UserVerifyEmailHandler 验证邮箱验证码的Handler
func UserVerifyEmailHandler(c *gin.Context) {
	// 从PostForm中获取code和email
	code, email := c.PostForm("code"), c.PostForm("email")

	if code == "" || email == "" {
		c.JSON(400, gin.H{"msg": "Code or Email is empty."})
		return
	}

	respond := services.UserVerifyEmailService(code, email)
	if respond == true {
		c.JSON(200, gin.H{"msg": "Email Verify Successfully!"})
	} else {
		c.JSON(400, gin.H{"msg": respond})
	}

}

// UserEmailExistHandler 检查邮箱是否已经注册
func UserEmailExistHandler(c *gin.Context) {
	//	获取get拼接的参数
	email := c.Query("email")
	if email == "" {
		c.JSON(400, gin.H{"msg": "Email is empty."})
		return
	}
	if !utils.ToolValidateValue(email, utils.Email) {
		c.JSON(400, gin.H{"msg": "Email is not valid."})
		return
	}

	have := model.DBIfEmailExistByEmail(email)
	if have {
		c.JSON(200, gin.H{"msg": "Email is exist.", "data": true})
	} else {
		c.JSON(200, gin.H{"msg": "Email is not exist.", "data": false})
	}
}

// UserRegisterHandler 注册用户
func UserRegisterHandler(c *gin.Context) {
	//	用户会提交一个表单（POST）
	// 	包含用户名、邮箱、验证码和密码，先验证验证码是否正确，然后验证邮箱是否存在，用户名是否唯一
	username, email, code, password := c.PostForm("username"), c.PostForm("email"), c.PostForm("code"), c.PostForm("password")
	if code == "" || email == "" || password == "" {
		c.JSON(400, gin.H{"msg": "Code or Email or Password is empty."})
		return
	}

	if utils.ToolValidateValue(email, utils.Email) == false {
		c.JSON(400, gin.H{"msg": "Email is not valid."})
		return
	}

	if utils.ToolValidateValue(username, utils.Username) == false {
		c.JSON(400, gin.H{"msg": "Username is not valid."})
		return
	}

	if model.DBIfEmailExistByEmail(email) {
		c.JSON(400, gin.H{"msg": "Email is exist."})
		return
	}

	if !model.DBGetEmailVerifyByCodeAndEmail(code, email) {
		c.JSON(400, gin.H{"msg": "Code is not valid."})
		return
	}

	if model.DBIfEmailExistByUsername(username) {
		c.JSON(400, gin.H{"msg": "Username is exist."})
		return
	}

	var UserAccount model.UserAccount
	UserAccount.Email = email
	UserAccount.Username = username
	UserAccount.Password, _ = utils.HashPassword(password)
	UserAccount.Status = "activate"
	UserAccount.Nickname = utils.ToolGenerateNickname()

	var UserInfo model.UserInfo
	UserInfo.RealName = ""
	UserInfo.Credit = 0
	UserInfo.RegisterIP = c.ClientIP()
	UserInfo.LastLoginIP = c.ClientIP()

	if model.DBUserAdd(&UserAccount) && model.DBUserAddInfo(&UserInfo) {
		c.JSON(200, gin.H{"msg": "Register Successfully!"})
	} else {
		c.JSON(400, gin.H{"msg": "Register Failed."}) // won't be triggered
	}
}

// UserLoginHandler 用户登录的处理
func UserLoginHandler(c *gin.Context) {
	//	获取用户名和密码
	username, password := c.PostForm("username"), c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(200, gin.H{"msg": "Username or Password is empty."})
		return
	}
	if !utils.ToolValidateValue(username, utils.Username) {
		c.JSON(400, gin.H{"msg": "Username is not valid."})
		return
	}
	success, token := services.UserLoginService(username, password)
	if success {
		c.JSON(200, gin.H{"msg": "Login Successfully!", "data": token})
	} else {
		c.JSON(400, gin.H{"msg": "Login Failed.", "data": false})
	}
}

// UserGetInfoHandler 获取用户信息
func UserGetInfoHandler(c *gin.Context) {
	userid, _ := c.Get("user")
	account, info, ok := model.DBUserGetDetailsByID(userid.(uint))
	fmt.Println(userid, account, info, ok)
	if ok {
		c.JSON(200, gin.H{"msg": "Get User Info Successfully!", "data": gin.H{"account": account, "info": info}})
		return
	}
	c.JSON(400, gin.H{"msg": "Get User Info Failed.", "data": false})
}

// UserEditInfoHandler 编辑用户信息
func UserEditInfoHandler(c *gin.Context) {
	userid, _ := c.Get("user")

	// 解析 JSON 数据
	var updateData map[string]interface{}
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(400, gin.H{"msg": "Invalid data format", "data": false})
		return
	}

	// 调用 Service 层处理业务逻辑
	if ok := UserServiceEditInfo(userid.(uint), updateData); ok {
		c.JSON(200, gin.H{"msg": "User info updated successfully", "data": true})
	} else {
		c.JSON(500, gin.H{"msg": "Failed to update user info", "data": false})
	}
}

// UserServiceEditInfo 用户信息编辑业务逻辑
func UserServiceEditInfo(userID uint, updateData map[string]interface{}) bool {
	// 调用数据库函数来更新信息
	return model.DBUserEditInfo(userID, updateData)
}

// UserSiteMessageSendHandler 发送站内信
func UserSiteMessageSendHandler(c *gin.Context) {
	// 这里仅示例，具体实现需要根据前端数据格式进行调整
	var messageData struct {
		SenderID  uint
		Content   string
		Recipient []uint // 接收者列表
	}
	// SenderId从c.Get("user")获取
	//SenderID, _ := c.Get("user")
	SenderID := c.PostForm("sender_id")
	// 转换为SenderID 为 uint
	NewSenderID, _ := strconv.Atoi(SenderID)

	// Content从c.PostForm("content")获取
	Content := c.PostForm("content")
	// TODO: 验证消息合法性
	// Recipient从c.PostForm("recipient")获取，用英文逗号分割字符串
	RecipientRaw := c.PostForm("recipient")
	// TODO: 验证接收者合法性
	RecipientSlice := strings.Split(RecipientRaw, ",")

	messageData.SenderID = uint(NewSenderID)
	messageData.Content = Content
	messageData.Recipient = make([]uint, 0)
	for _, recipient := range RecipientSlice {
		recipientID, _ := strconv.Atoi(recipient)
		messageData.Recipient = append(messageData.Recipient, uint(recipientID))
	}

	success := services.UserSiteMessageSendService(messageData.SenderID, messageData.Content, messageData.Recipient)
	if success {
		c.JSON(200, gin.H{"msg": "Message sent successfully!", "data": true})
	} else {
		c.JSON(400, gin.H{"msg": "Failed to send message", "data": false})
	}
}

// UserSiteMessageGetAllHandler 获取当前用户所有站内信
func UserSiteMessageGetAllHandler(c *gin.Context) {
	userid, _ := c.Get("user") // 假设用户ID已通过认证中间件添加到上下文
	messages, ok := UserSiteMessageGetAllService(userid.(uint))
	if ok {
		c.JSON(200, gin.H{"msg": "Get Messages Successfully!", "data": messages})
		return
	}
	c.JSON(400, gin.H{"msg": "Failed to get messages.", "data": nil})
}

// UserSiteMessageGetAllService 获取所有站内信业务逻辑
func UserSiteMessageGetAllService(userID uint) ([]model.SiteMessageDetail, bool) {
	return model.DBSiteMessageGetAllByUserID(userID)
}
