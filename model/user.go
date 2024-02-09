package model

import (
	"OpenTeens/dao"
	"OpenTeens/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// UserAccount 用户账号表
type UserAccount struct {
	gorm.Model
	Username string
	Nickname string
	Password string // 加密后的密码
	Email    string `gorm:"unique"`
	Status   string // 账户状态
}

// UserInfo 用户信息表
type UserInfo struct {
	gorm.Model
	RealName    string
	Credit      int
	RegisterIP  string
	LastLoginIP string
}

// EmailVerification 邮箱验证码表
type EmailVerification struct {
	gorm.Model
	Email     string
	Code      string
	ExpiresAt time.Time // 验证码到期时间
	Used      bool      // 验证码是否已使用
	IpAddr    string    // ip地址
}

// UserToken 用户Token表
type UserToken struct {
	gorm.Model
	UserAccountID uint
	Token         string
	ExpiresAt     time.Time
	LoginMethod   string
	LastUseAt     time.Time
}

// SiteMessage 站内信表
type SiteMessage struct {
	gorm.Model
	SenderID uint      // 发送者ID
	Content  string    // 站内信内容
	SentAt   time.Time // 发送时间
}

// SiteMessageRecipient 站内信接收者表
type SiteMessageRecipient struct {
	gorm.Model
	MessageID  uint // 站内信ID，外键，关联到SiteMessage
	ReceiverID uint // 接收者ID，外键，关联到UserAccount
	IsRead     bool // 是否已读
}

func DBUserAddEmailCode(email string, code string, ipaddr string) {
	//	创建一条记录，设置Email和Code的值，设置ExpiresAt为30分钟后，Used为false
	dao.DB.Create(&EmailVerification{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Minute * 30),
		Used:      false,
		IpAddr:    ipaddr,
	})
}

// DBUserGetEmailCodeByIp 获取邮箱验证码N分钟内ipaddr的发送记录
func DBUserGetEmailCodeByIp(ipaddr string, minutes int) int {
	var count int64

	// 计算从当前时间起过去N分钟的时间
	pastTime := time.Now().Add(-time.Duration(minutes) * time.Minute)

	// 在 EmailVerification 表中查找 CreatedAt 在过去N分钟内的记录
	dao.DB.Model(&EmailVerification{}).
		Where("created_at > ?", pastTime).
		Where("ip_addr = ?", ipaddr).
		Count(&count)

	// 返回记录数
	return int(count)
}

// DBUserGetEmailCodeByEmail 获取邮箱验证码N分钟内email的发送记录
func DBUserGetEmailCodeByEmail(email string, minutes int) int {
	var count int64

	// 计算从当前时间起过去N分钟的时间
	pastTime := time.Now().Add(-time.Duration(minutes) * time.Minute)

	// 在 EmailVerification 表中查找 CreatedAt 在过去N分钟内的记录
	dao.DB.Model(&EmailVerification{}).
		Where("created_at > ?", pastTime).
		Where("email = ?", email).
		Count(&count)

	// 返回记录数
	return int(count)
}

// DBGetEmailVerifyByCodeAndEmail 根据验证码和邮箱获取邮箱验证码
func DBGetEmailVerifyByCodeAndEmail(code string, email string) bool {
	var emailVerifications []EmailVerification

	// 使用Find而不是First
	dao.DB.Where("code = ? AND email = ? AND expires_at > ? AND used = ?", code, email, time.Now(), false).Find(&emailVerifications)

	// 检查是否有记录
	if len(emailVerifications) == 0 {
		// 没有找到记录的处理
		fmt.Println("No record found")
		return false
	}

	// 找到记录，更新Used字段
	dao.DB.Model(&emailVerifications[0]).Update("used", true)
	return true
}

// DBIfEmailExistByEmail 根据邮箱判断是否存在
func DBIfEmailExistByEmail(email string) bool {
	var userAccount []UserAccount

	// 在 UserAccount 表中查找匹配的 email 记录
	dao.DB.Where("email = ?", email).Find(&userAccount)

	// 如果找到记录，返回 true
	if len(userAccount) > 0 {
		return true
	}

	// 如果未找到记录或发生其他错误，返回 false
	return false
}

// DBIfEmailExistByUsername 根据邮箱判断是否存在
func DBIfEmailExistByUsername(username string) bool {
	var userAccount []UserAccount

	// 在 UserAccount 表中查找匹配的 email 记录
	dao.DB.Where("email = ?", username).Find(&userAccount)

	// 如果找到记录，返回 true
	if len(userAccount) > 0 {
		return true
	}

	// 如果未找到记录或发生其他错误，返回 false
	return false
}

func DBUserAdd(userAccount *UserAccount) bool {
	// 保存用户账号信息
	dao.DB.Create(&userAccount)
	return true
}

func DBUserAddInfo(userInfo *UserInfo) bool {
	// 保存用户账号信息
	dao.DB.Create(&userInfo)
	return true
}

func DBUserLoginCreateToken(userAccount UserAccount, method string) (bool, string) {
	var existingToken UserToken

	// 查找与此账户关联的、未过期的 Token
	result := dao.DB.Where("user_account_id = ? AND expires_at > ?", userAccount.ID, time.Now()).First(&existingToken)

	// 检查是否找到 Token
	if result.Error == nil {
		// 找到了 Token，更新 LastUseAt 并返回
		dao.DB.Model(&existingToken).Update("last_use_at", time.Now())
		return true, existingToken.Token
	}

	// 生成新的 Token
	token := utils.GenerateToken()
	// 设置 Token 过期时间
	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	// 创建新的 UserToken 记录
	dao.DB.Create(&UserToken{
		Token:         token,
		ExpiresAt:     expiresAt,
		UserAccountID: userAccount.ID,
		LoginMethod:   method,
		LastUseAt:     time.Now(),
	})

	return true, token
}

// LoginMethod 定义登录方式的枚举类型
type LoginMethod int

const (
	LoginByUsername LoginMethod = iota
	LoginByEmail
	LoginByUserID
)

// DBUserCheckAccount 通过不同方式和密码判断用户是否存在
func DBUserCheckAccount(loginValue string, password string, method LoginMethod) (bool, UserAccount) {
	var userAccount []UserAccount
	var field string

	switch method {
	case LoginByUsername:
		field = "username"
	case LoginByEmail:
		field = "email"
	case LoginByUserID:
		field = "user_id"
	}

	query := fmt.Sprintf("%s = ? AND password = ?", field)
	dao.DB.Where(query, loginValue, password).Find(&userAccount)

	if len(userAccount) > 0 {
		fmt.Println("User Account Found")
		return true, userAccount[0]
	} else {
		fmt.Println("User Account Not Found")
		return false, UserAccount{}
	}
}

// DBUserGetAccountIDFromToken 根据 Token 获取用户 ID
func DBUserGetAccountIDFromToken(token string) (uint, bool) {
	var userTokens []UserToken

	dao.DB.Where("token = ? AND expires_at > ?", token, time.Now()).Find(&userTokens)
	if len(userTokens) == 0 {
		return 0, false // Token 无效或找不到
	}

	return userTokens[0].UserAccountID, true
}

// DBUserGetDetailsByID 根据用户 ID 获取用户的详细信息
func DBUserGetDetailsByID(userID uint) (*UserAccount, *UserInfo, bool) {
	var userAccount UserAccount
	var userInfo UserInfo

	// 获取 UserAccount 信息
	accountResult := dao.DB.Where("id = ?", userID).First(&userAccount)
	if accountResult.Error != nil {
		return nil, nil, false
	}

	// 获取 UserInfo 信息
	infoResult := dao.DB.Where("id = ?", userID).First(&userInfo)
	if infoResult.Error != nil {
		return nil, nil, false
	}

	// 清除敏感信息
	userAccount.Password = ""
	// 可以继续清除或修改其他敏感字段

	return &userAccount, &userInfo, true
}

// DBUserEditInfo 更新用户信息
func DBUserEditInfo(userID uint, updateData map[string]interface{}) bool {
	// 更新 UserAccount 和 UserInfo 表
	result1 := dao.DB.Model(&UserAccount{}).Where("id = ?", userID).Updates(updateData)
	result2 := dao.DB.Model(&UserInfo{}).Where("user_account_id = ?", userID).Updates(updateData)

	return result1.Error == nil && result2.Error == nil
}

// DBUserSiteMessageCreate 创建站内信记录
func DBUserSiteMessageCreate(message *SiteMessage) {
	dao.DB.Create(&message)
}

// DBUserSiteMessageRecipientCreate 创建站内信接收者记录
func DBUserSiteMessageRecipientCreate(recipient *SiteMessageRecipient) {
	dao.DB.Create(&recipient)
}

// SiteMessageDetail 包含站内信详细信息的结构体
type SiteMessageDetail struct {
	ID             uint      `json:"id"`
	SenderID       uint      `json:"sender_id"`
	SenderUsername string    `json:"sender_username"` // 发送者用户名
	Content        string    `json:"content"`
	SentAt         time.Time `json:"sent_at"`
	IsRead         bool      `json:"is_read"`
}

// DBSiteMessageGetAllByUserID 根据用户ID获取所有站内信
func DBSiteMessageGetAllByUserID(userID uint) ([]SiteMessageDetail, bool) {
	var messageDetails []SiteMessageDetail
	// 联表查询填充SiteMessageDetail结构体
	dao.DB.Table("site_messages").
		Select("site_messages.id, site_messages.sender_id, user_accounts.username as sender_username, site_messages.content, site_messages.sent_at, site_message_recipients.is_read").
		Joins("join site_message_recipients on site_message_recipients.message_id = site_messages.id").
		Joins("left join user_accounts on site_messages.sender_id = user_accounts.id").
		Where("site_message_recipients.receiver_id = ?", userID).
		Scan(&messageDetails)

	if len(messageDetails) > 0 {
		return messageDetails, true
	} else {
		return nil, false
	}
}
