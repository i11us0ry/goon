package main

import (
	"fmt"
	"goon3/check"
	"goon3/public"
	"time"
)

func main() {
	sTime := time.Now().Unix()
	// 检查配置文件
	check.ConfigCheck()
	public.Init(check.ConfigRead(&public.Conf{}))
	// 初始化输入指令和打印banner
	public.Flag()

	//处理输入指令
	check.CheckInput()
	// 根据mode下发任务
	check.CheckMode()
	fmt.Println()
	eTime := time.Now().Unix()
	public.Info.Println(fmt.Sprintf("running time：%ds", eTime-sTime))
}
