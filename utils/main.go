package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

// ToolGenerateRandomString 生成随机字符串
func ToolGenerateRandomString(charset string, length int) string {
	rand.Seed(time.Now().UnixNano())
	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		randomString[i] = charset[randomIndex]
	}
	return string(randomString)
}

// GenerateToken 随机生成32位的Token
func GenerateToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	return ToolGenerateRandomString(charset, 32)
}

// GenerateRandomNumber 随机生成6位数字
func GenerateRandomNumber() string {
	const charset = "0123456789"
	return ToolGenerateRandomString(charset, 6)
}

// ValueType 用于定义值类型
type ValueType int

const (
	Email ValueType = iota
	Username
	Nickname
)

// ToolValidateValue 根据不同类型验证值
func ToolValidateValue(value string, valueType ValueType) int {
	switch valueType {
	case Email:
		return validateEmail(value)
	case Username:
		return validateUsername(value)
	case Nickname:
		return validateNickname(value)
	default:
		return 1
	}
}

// validateEmail 验证邮箱格式
func validateEmail(email string) int {
	// 验证邮箱格式：name@domain
	splited := strings.Split(email, "@")
	if len(splited) != 2 {
		return -1
	}

	// extract name and domain
	emailName, emailDomain := splited[0], splited[1]

	// check1: name
	emailNamePattern := `^[a-zA-Z0-9._%+-]+$`
	success1, _ := regexp.MatchString(emailNamePattern, emailName)
	if !success1 {
		return -1
	}

	// check2: domain
	trustedEmailDomains := []string{"openteens.org", "gmail.com"}
	for _, domain := range trustedEmailDomains {
		if emailDomain == domain {
			return 0
		}
	}

	// domain not supported
	return -2
}

// validateUsername 验证用户名
func validateUsername(username string) int {
	usernamePattern := `^[a-zA-Z0-9]{6,18}$`
	success, _ := regexp.MatchString(usernamePattern, username)
	if success {
		return 0
	}
	return -1
}

// validateNickname 验证昵称
func validateNickname(nickname string) int {
	// 验证长度 (16个字符以内，一个汉字算作两个字符)
	if utf8.RuneCountInString(nickname) <= 16 {
		return 0
	}
	return -1
}

// ToolUserEmailCheck 检查邮箱是否符合格式或已经注册

const salt = "T1H2I3S4I5S6A7V8E0R1Y2L3O4N5G6A7N8D9U0NI.QU3E3H1A1S1H2S2A2L2T2"

// HashPassword 创建带有固定盐值的哈希密码
func HashPassword(password string) (string, error) {
	// 将盐值附加到密码末尾
	passwordWithSalt := password + salt

	// 使用 SHA-256 哈希函数
	hash := sha256.Sum256([]byte(passwordWithSalt))
	return hex.EncodeToString(hash[:]), nil
}

// CheckPasswordHash 验证密码与哈希是否匹配
func CheckPasswordHash(password, hash string) bool {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return false
	}
	return hashedPassword == hash
}
