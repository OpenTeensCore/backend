package router

import (
	"OpenTeens/controller"
	"OpenTeens/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetUpRoute() *gin.Engine {
	r := gin.Default()

	//r.Static("/static", "static")
	//r.LoadHTMLGlob("templates/*")
	r.Use(middleware.MidCors())

	r.GET("/", controller.IndexHandler)
	v1Group := r.Group("v1")
	{
		// 定义用户路由组
		userGroup := v1Group.Group("user")
		{
			// 主路由
			userGroup.GET("/", middleware.AuthMiddleware(), func(context *gin.Context) {
				context.JSON(200, gin.H{"message": "Hello Tsinglan!"})
			})

			// 定义用户账号路由组
			userAccount := userGroup.Group("account")
			{
				// 主路由
				userAccount.GET("/", func(context *gin.Context) {
					context.JSON(200, gin.H{"message": "Hello Tsinglan!"})
				})

				// 发送验证码
				userAccount.POST("/sendEmail", controller.UserSendEmailHandler)
				// 验证邮箱
				userAccount.POST("/verifyEmail", controller.UserVerifyEmailHandler)
				// 判断邮箱是否注册get
				userAccount.GET("/isExistEmail", controller.UserEmailExistHandler)
				// 注册接口Post
				userAccount.POST("/reg", controller.UserRegisterHandler)
				// 登录接口Post
				userAccount.POST("/login", controller.UserLoginHandler)
				// 获取用户信息
				userAccount.GET("/getUserInfo", middleware.AuthMiddleware(), controller.UserGetInfoHandler)
			}

			// SiteMessage
			siteMessage := userGroup.Group("siteMessage")
			{
				// 发送消息
				siteMessage.POST("/send", controller.UserSiteMessageSendHandler)
				// 获取消息
				siteMessage.GET("/get", middleware.AuthMiddleware(), controller.UserSiteMessageGetAllHandler)
			}

			// Attachment
			attachment := userGroup.Group("attachment")
			{
				// 上传附件
				attachment.PUT("/upload", middleware.AuthMiddleware(), controller.AttachmentUploadHandler)
				// 获取附件
				attachment.GET("/get", controller.AttachmentGetHandler)
			}
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "Api Not Fount",
		})
	})

	return r
}
