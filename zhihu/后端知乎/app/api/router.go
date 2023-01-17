package api

import (
	"github.com/gin-gonic/gin"
	"main/app/api/middleware"
)

func InitRouter() error {
	r := gin.Default()
	r.Use(middleware.CORS())
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/checkCaptcha", CheckCaptcha)
	//r.POST(fmt.Sprintf("/home/topic/%s/article/%s/post", TopicID, ArticleID), WriteArticle)
	r.GET("/findPassword", FindPassword)
	r.GET("/getCaptcha", CreateCaptcha)
	UserRouter := r.Group("/user")
	{
		UserRouter.Use(middleware.JWTAuthMiddleware())
		UserRouter.GET("/get", getUsernameFromToken)
		UserRouter.GET("/:uid", GetUser)
		UserRouter.POST("/:uid/changePassword", ChangePassword)
		UserRouter.POST("/:uid/changeNickname", ChangeNickname)
		UserRouter.POST("/:uid/changeAvatar", ChangeAvatar)
		UserRouter.POST("/:uid/changeIntroduction", ChangeIntroduction)
		UserRouter.POST("/:uid/postQuestion", PostQuestion)
		UserRouter.POST("/:uid/question/:topicID/postArticle", PostArticle)
		UserRouter.POST("/:uid/writeArticle", WriteArticle)
		UserRouter.GET("/:uid/getQuestion", GetQuestion)
		UserRouter.GET("/:uid/question/:topicID", GetDirectQuestion)
		UserRouter.GET("/:uid/question/:topicID/article", GetArticle)
		UserRouter.POST("/:uid/question/:topicID/article/:articleID/postComment", AddComment)
		UserRouter.GET("/:uid/question/:topicID/article/:articleID/comment", GetComment)
		UserRouter.POST("/:uid/question/:topicID/addLike", AddCollection)
		UserRouter.POST("/:uid/question/:topicID/deleteLike", DeleteCollection)
		UserRouter.GET("/:uid/getCollection", GetCollection)
		//以下为未实装的接口
		UserRouter.POST("/:uid/search", Search)
		UserRouter.POST("/:uid/question/:topicID/article/like", AddLike)
		UserRouter.POST("/:uid/question/:topicID/article/like", DeleteLike)
	}
	r.POST("/uploadCaptcha")
	Home := r.Group("/home")
	{
		Home.Use(middleware.JWTAuthMiddleware())
		Home.GET("/get", getUsernameFromToken)
		//Home.GET(fmt.Sprintf("/user/%s", id), getUser)
	}
	err := r.Run(":8088")
	if err != nil {
		return err
	} else {
		return nil
	}

}
