package redis

import (
	"main/model"
)

// VerifyCaptcha @Pram id 验证码id
// @Pram VerifyValue 用户输入的答案
// @Result true：正确，false：失败
func VerifyCaptcha(id, VerifyValue string) bool {
	return model.Result.Verify(id, VerifyValue, true)

}

// GetCodeAnswer @Pram codeId 验证码id
// @Result 验证码答案
func GetCodeAnswer(codeId string) string {
	return model.Result.Get(codeId, false)
}


