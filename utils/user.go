package utils

// ToolGenerateNickname Nickname生成器
func ToolGenerateNickname() string {
	return "高中用户_" + ToolGenerateRandomString("0123456789", 6)
}
