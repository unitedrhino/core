package main

import (
	"fmt"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("请输入用户id和需要生成的密码")
		return
	}
	fmt.Printf("需要生成的用户id:%v  密码:%v\n", args[1], args[2])
	pwd := utils.MakePwd(args[2], cast.ToInt64(args[1]), false)
	fmt.Printf("生成的密码为:%v\n", pwd)
}
