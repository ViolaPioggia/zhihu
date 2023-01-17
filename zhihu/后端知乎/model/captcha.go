package model

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
	"time"
)

var Result = base64Captcha.NewMemoryStore(20240, 3*time.Minute)

func StringConfig() *base64Captcha.DriverString {
	stringType := &base64Captcha.DriverString{
		Height:          100,
		Width:           50,
		NoiseCount:      0,
		ShowLineOptions: base64Captcha.OptionShowHollowLine | base64Captcha.OptionShowSlimeLine,
		Length:          4,
		Source:          "6", //qwertyuiopasdfghjklzxcvb
		BgColor: &color.RGBA{
			R: 40,
			G: 30,
			B: 89,
			A: 29,
		},
		Fonts: nil,
	}
	return stringType
}

type Captcha struct {
	Base64 string `form:"base64" json:"base64" `
	Key    int    `form:"key" json:"key" `
	UID    string `form:"uid" json:"uid" `
}
