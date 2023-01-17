package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"main/app/global"
	"main/dao/mysql"
	"main/dao/redis"
	"main/model"
	"main/utils"
	"net/http"
	"time"
	"unicode/utf8"
)

func GetDirectQuestion(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	articleID := c.Param("topicID")
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"title":   global.Rdb.Get(c, fmt.Sprintf("topic:%s:title", articleID)).Val(),
		"context": global.Rdb.Get(c, fmt.Sprintf("topic:%s:context", articleID)).Val(),
		"time":    global.Rdb.Get(c, fmt.Sprintf("topic:%s:time", articleID)).Val(),
	})
}

func PostArticle(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	topicID := c.Param("topicID")
	var article model.Article
	if err := c.ShouldBind(&article); err != nil {
		fmt.Println(err)
		utils.ResponseFail(c, "verification failed")
		return
	}
	context := article.Context
	if utf8.RuneCountInString(context) > 30000 {
		utils.ResponseFail(c, "描述太长")
		return
	}
	err := redis.Set(c, fmt.Sprintf("article:%d:context", mysql.FindArticleID()), context, 0)
	_ = redis.Set(c, fmt.Sprintf("article:%d:username", mysql.FindArticleID()), fmt.Sprintf("%s", username), 0)
	_ = redis.Set(c, fmt.Sprintf("article:%d:userid", mysql.FindArticleID()), Userid, 0)
	_ = redis.Set(c, fmt.Sprintf("article:%d:topicid", mysql.FindArticleID()), topicID, 0)
	_ = redis.Set(c, fmt.Sprintf("article:%v:time", mysql.FindArticleID()), time.Now().Format("2006-01-02 15:04:05"), 0)
	if err != nil {
		utils.ResponseFail(c, "write into redis failed")
		return
	}
	flag1, msg := mysql.AddArticle(context, Userid, topicID) //写入数据库
	if flag1 {
		utils.ResponseSuccess(c, "post article success")
	} else {
		utils.ResponseFail(c, fmt.Sprintf("post article failed,%s", msg))
	}
}

func GetArticle(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	topicID := c.Param("topicID")
	for i := 1; global.Rdb.Get(c, fmt.Sprintf("article:%d:topicid", i)).Val() != ""; i++ {
		if global.Rdb.Get(c, fmt.Sprintf("article:%d:topicid", i)).Val() == topicID {

			c.JSON(http.StatusOK, gin.H{
				"status":   200,
				"context":  global.Rdb.Get(c, fmt.Sprintf("article:%d:context", i)).Val(),
				"time":     global.Rdb.Get(c, fmt.Sprintf("article:%d:time", i)).Val(),
				"nickname": global.Rdb.Get(c, fmt.Sprintf("%s:nickname", global.Rdb.Get(c, fmt.Sprintf("article:%d:username", i)).Val())).Val(),
				"avatar":   global.Rdb.Get(c, fmt.Sprintf("%s:avatar", global.Rdb.Get(c, fmt.Sprintf("article:%d:username", i)).Val())).Val(),
			})
		}
	}
}
func WriteArticle(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	var article model.Topic
	if err := c.ShouldBind(&article); err != nil {
		fmt.Println(err)
		utils.ResponseFail(c, "verification failed")
		return
	}

	title := article.Title
	context := article.Context

	if utf8.RuneCountInString(title) > 15 {
		utils.ResponseFail(c, "标题太长")
		return
	}
	if utf8.RuneCountInString(context) > 30000 {
		utils.ResponseFail(c, "描述太长")
		return
	}
	err := redis.Set(c, fmt.Sprintf("article:%d:context", mysql.FindArticleID()), context, 0)
	_ = redis.Set(c, fmt.Sprintf("article:%d:username", mysql.FindArticleID()), fmt.Sprintf("%s", username), 0)
	_ = redis.Set(c, fmt.Sprintf("article:%d:userid", mysql.FindArticleID()), Userid, 0)
	_ = redis.Set(c, fmt.Sprintf("article:%d:topicid", mysql.FindArticleID()), "0", 0)
	_ = redis.Set(c, fmt.Sprintf("article:%v:time", mysql.FindArticleID()), time.Now().Format("2006-01-02 15:04:05"), 0)
	if err != nil {
		utils.ResponseFail(c, "write into redis failed")
		return
	}
	flag1, msg := mysql.AddArticle(context, Userid, "0") //写入数据库
	if flag1 {
		utils.ResponseSuccess(c, "post article success")
	} else {
		utils.ResponseFail(c, fmt.Sprintf("post article failed,%s", msg))
	}
}
func AddComment(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	c.Param("topicID")
	articleID := c.Param("articleID")
	var comment model.Comment
	if err := c.ShouldBind(&comment); err != nil {
		fmt.Println(err)
		utils.ResponseFail(c, "verification failed")
		return
	}

	context := comment.Context
	if utf8.RuneCountInString(context) > 1000 {
		utils.ResponseFail(c, "描述太长")
		return
	}
	err := redis.Set(c, fmt.Sprintf("comment:%d:context", mysql.FindCommentID()), context, 0)
	_ = redis.Set(c, fmt.Sprintf("comment:%d:username", mysql.FindCommentID()), fmt.Sprintf("%s", username), 0)
	_ = redis.Set(c, fmt.Sprintf("comment:%d:userid", mysql.FindCommentID()), Userid, 0)
	_ = redis.Set(c, fmt.Sprintf("comment:%d:articleid", mysql.FindCommentID()), articleID, 0)
	_ = redis.Set(c, fmt.Sprintf("comment:%v:time", mysql.FindCommentID()), time.Now().Format("2006-01-02 15:04:05"), 0)
	if err != nil {
		utils.ResponseFail(c, "write into redis failed")
		return
	}
	flag1, msg := mysql.AddComment(context, Userid, articleID, fmt.Sprintf("%s", username)) //写入数据库
	if flag1 {
		utils.ResponseSuccess(c, "post comment success")
	} else {
		utils.ResponseFail(c, fmt.Sprintf("post comment failed,%s", msg))
	}
}

func GetComment(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	c.Param("topicID")
	articleID := c.Param("articleID")
	for i := 1; global.Rdb.Get(c, fmt.Sprintf("comment:%d:articleid", i)).Val() != ""; i++ {
		if global.Rdb.Get(c, fmt.Sprintf("comment:%d:articleid", i)).Val() == articleID {

			c.JSON(http.StatusOK, gin.H{
				"status":   200,
				"context":  global.Rdb.Get(c, fmt.Sprintf("comment:%d:context", i)).Val(),
				"time":     global.Rdb.Get(c, fmt.Sprintf("comment:%d:time", i)).Val(),
				"nickname": global.Rdb.Get(c, fmt.Sprintf("%s:nickname", global.Rdb.Get(c, fmt.Sprintf("comment:%d:username", i)).Val())).Val(),
				"avatar":   global.Rdb.Get(c, fmt.Sprintf("%s:avatar", username)).Val(),
			})
		}
	}
}
func AddLike(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	c.Param("topicID")
	articleID := c.Param("articleID")
	isMember, err := global.Rdb.SIsMember(c, fmt.Sprintf("article:%s:like", articleID), Userid).Result()
	if err != nil {
		utils.ResponseFail(c, "addLike failed")
		return
	}
	if isMember {
		utils.ResponseFail(c, "已经点赞过了")
		return
	}
	global.Rdb.SAdd(c, fmt.Sprintf("article:%s:like", articleID), Userid)
	utils.ResponseSuccess(c, "点赞成功")
}

func DeleteLike(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	c.Param("topicID")
	articleID := c.Param("articleID")
	isMember, err := global.Rdb.SIsMember(c, fmt.Sprintf("article:%s:like", articleID), Userid).Result()
	if err != nil {
		utils.ResponseFail(c, "deleteLike failed")
		return
	}
	if !isMember {
		utils.ResponseFail(c, "没有点赞过")
		return
	}
	global.Rdb.SRem(c, fmt.Sprintf("article:%s:like", articleID), Userid)
	utils.ResponseSuccess(c, "取消点赞成功")
}
