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

func PostQuestion(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)

	var topic model.Topic
	if err := c.ShouldBind(&topic); err != nil {
		fmt.Println(err)
		utils.ResponseFail(c, "verification failed")
		return
	}
	title := topic.Title
	context := topic.Context
	//tags := topic.Tags
	if utf8.RuneCountInString(title) > 15 {
		utils.ResponseFail(c, "标题太长")
		return
	}
	if utf8.RuneCountInString(context) > 1000 {
		utils.ResponseFail(c, "描述太长")
		return
	}
	//tagsArr := strings.Split(tags, ",")
	flag, _ := redis.GetPassword(c, fmt.Sprintf("topic:%s", title))
	if flag != "" {
		utils.ResponseFail(c, "topic already exists")
		return
	}

	err := redis.Set(c, fmt.Sprintf("topic:%d:context", mysql.FindTopicID()), context, 0)
	_ = redis.Set(c, fmt.Sprintf("topic:%d:userid", mysql.FindTopicID()), Userid, 0)
	_ = redis.Set(c, fmt.Sprintf("topic:%d:title", mysql.FindTopicID()), title, 0)
	_ = redis.Set(c, fmt.Sprintf("topic:%v:time", mysql.FindTopicID()), time.Now().Format("2006-01-02 15:04:05"), 0)
	if err != nil {
		utils.ResponseFail(c, "write into redis failed")
		return
	}
	//_ = redis.ListSet(c, fmt.Sprintf("topic:%s:tags", title), tagsArr)
	flag1, msg := mysql.AddTopic(UserID, title, context) //写入数据库
	if flag1 {
		utils.ResponseSuccess(c, "post topic success")
	} else {
		utils.ResponseFail(c, fmt.Sprintf("post topic failed,%s", msg))
	}
}

func GetQuestion(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)

	//var ArrTitle []string
	//var ArrName []string
	//var ArrTime []string

	for i := 1; global.Rdb.Get(c, fmt.Sprintf("topic:%d:context", i)).Val() != ""; i++ {
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"title":  global.Rdb.Get(c, fmt.Sprintf("topic:%d:title", i)).Val(),
			"name":   username,
			"time":   global.Rdb.Get(c, fmt.Sprintf("topic:%d:time", i)).Val(),
		})
		//ArrTitle[i-1] = global.Rdb.Get(c, fmt.Sprintf("topic:%d:title", i)).Val()
		//ArrName[i-1] = fmt.Sprintf("%s", username)
		//ArrTime[i-1] = global.Rdb.Get(c, fmt.Sprintf("topic:%d:time", i)).Val()
	}
	//c.JSON(http.StatusOK, gin.H{
	//	"status": 200,
	//	"title":  ArrTitle,
	//	"name":   ArrName,
	//	"time":   ArrTime,
	//})
}
func AddCollection(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	topicID := c.Param("topicID")
	global.Rdb.LPush(c, fmt.Sprintf("%s:collection", username), global.Rdb.Get(c, fmt.Sprintf("topic:%s:context", topicID)).Val())
	utils.ResponseSuccess(c, "")
}
func DeleteCollection(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	topicID := c.Param("topicID")
	err := global.Rdb.LRem(c, fmt.Sprintf("%s:collection", username), 1, global.Rdb.Get(c, fmt.Sprintf("topic:%s:context", topicID)).Val())
	if err != nil {
		utils.ResponseFail(c, "")
	}
	utils.ResponseSuccess(c, "")
}
func Search(c *gin.Context) {
	username, _ := c.Get("username")
	Userid := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Val()
	//UserID, _ := global.Rdb.Get(c, fmt.Sprintf("%s:%s", username, "id")).Int64()
	c.Param(Userid)
	var topic model.Topic
	if err := c.ShouldBind(&topic); err != nil {
		fmt.Println(err)
		utils.ResponseFail(c, "verification failed")
		return
	}
	title := topic.Title

	keysMatch, _, _ := global.Rdb.Scan(c, 0, "*"+title+"*", 20).Result()
	c.JSON(http.StatusOK, gin.H{
		"result": keysMatch,
	})

}

