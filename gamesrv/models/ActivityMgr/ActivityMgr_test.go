package ActivityMgr

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego"
)

func Test_GetActivity(t *testing.T) {
	fmt.Print("---------------------开始测试-------------------------\n")

	beego.Debug(Instance().ActivityInfoArray)

	fmt.Print("----------------------测试结束--------------------------\n")
}
