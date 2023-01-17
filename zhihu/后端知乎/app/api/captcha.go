package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"main/app/global"
	"main/dao/redis"
	"main/model"
	"main/utils"
	"net/http"
	"strconv"
	"time"
)

// CreateCode @Result id 验证码id
// @Result bse64s 图片base64编码
// @Result err 错误
func CreateCode(c *gin.Context) (string, string, error) {
	var driver base64Captcha.Driver

	driver = model.StringConfig()

	if driver == nil {
		panic("生成验证码的类型没有配置，请在yaml文件中配置完再次重试启动项目")
	}
	// 创建验证码并传入创建的类型的配置，以及存储的对象
	c1 := base64Captcha.NewCaptcha(driver, model.Result)
	id, b64s, err := c1.Generate()
	if err != nil {
		utils.ResponseFail(c, "create captcha failed")
		return "", "", err
	}
	return id, b64s, nil
}
func CreateCaptcha(c *gin.Context) {
	uid, b64s, err := CreateCode(c)
	if err != nil {
		utils.ResponseFail(c, "create captcha failed")
		return
	}
	key := 1
	for {
		if global.Rdb.Get(c, fmt.Sprintf("key:%s", strconv.Itoa(key))).Val() == "" {
			break
		} else {
			key++
		}
	}
	err = redis.Set(c, fmt.Sprintf("key:%s", strconv.Itoa(key)), uid, 2*time.Minute)
	if err != nil {
		utils.ResponseFail(c, "write password into redis failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"b64s":   b64s,
		"key":    key,
	})
}

func CheckCaptcha(c *gin.Context) {
	var captcha model.Captcha
	if err := c.ShouldBind(&captcha); err != nil {
		fmt.Println(err)
		utils.ResponseFail(c, "verification failed")
		return
	}
	key := captcha.Key
	UID := captcha.UID
	if UID == redis.GetCodeAnswer(global.Rdb.Get(c, fmt.Sprintf("key:%s", strconv.Itoa(key))).Val()) {
		utils.ResponseSuccess(c, "right captcha")
	} else {
		utils.ResponseFail(c, "wrong captcha")
	}

}
