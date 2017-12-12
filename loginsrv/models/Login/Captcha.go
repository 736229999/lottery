package Login

import (
	"io"

	"sync"

	"github.com/astaxie/beego"
	"github.com/lifei6671/gocaptcha"
)

var once sync.Once

func CreateCaptcha(w io.Writer, imageFormat int) string {
	//读取字体库
	once.Do(func() {
		err := gocaptcha.ReadFonts("conf/fonts", ".ttf")
		if err != nil {
			beego.Emergency(err)
		}
	})

	//生成一张验证码图片
	captchaImage, err := gocaptcha.NewCaptchaImage(Dx, Dy, gocaptcha.RandLightColor())
	if err != nil {
		return ""
	}

	//生成噪点，大概消耗3ms时间
	//captchaImage.DrawNoise(gocaptcha.CaptchaComplexHigh)

	//生成文字噪点，大概会消耗4ms时间
	//captchaImage.DrawTextNoise(gocaptcha.CaptchaComplexHigh)

	font := gocaptcha.RandText(4)
	//生成验证码（为了提高效率目前就只生成验证码了
	err_0 := captchaImage.DrawText(font)
	if err_0 != nil {
		return ""
	}

	//随机画3条线；
	//captchaImage.Drawline(3)

	//画边框；
	//captchaImage.DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A))

	//画白色空线条；
	//captchaImage.DrawHollowLine()

	captchaImage.SaveImage(w, gocaptcha.ImageFormatJpeg)

	return font
}

//删除验证码存储记录
func Destory(mp string) {
	delete(VerifyInfo, mp)
}
