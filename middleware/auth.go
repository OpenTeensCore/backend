package middleware

import (
	"OpenTeens/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization") // 从请求头中获取 Token

		// 调用 getUserFromToken 来验证 Token 并获取用户信息
		userid, isValid := MidGetUserFromToken(token)
		if !isValid {
			// Token 验证失败，返回错误信息并终止请求
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "unauthorized"})
			c.Abort()
			return
		}

		// 将用户信息添加到请求的上下文中
		c.Set("user", userid)

		// 继续处理请求
		c.Next()
	}
}

func MidGetUserFromToken(token string) (uint, bool) {
	// 从数据库中获取用户信息 DBUserGetAccountIDFromToken
	return model.DBUserGetAccountIDFromToken(token)
}

func MidCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
