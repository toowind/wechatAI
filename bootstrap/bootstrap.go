package bootstrap

import (
	"fmt"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/handlers"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/pkg/logger"
	"github.com/eatmoreapple/openwechat"
	"os"
)

func sendMessageToSomeOne(f *openwechat.Friend) {
	//for {
	//	now := time.Now()
	//	t := time.Date(now.Year(), now.Month(), now.Day(), 6, 20, 0, 0, now.Location())
	//	timer := time.NewTimer(now.Sub(t))
	//	<-timer.C
	//	f.SendText("现在是机器人自动给你发的问候~晚安~")
	//	break
	//}
	f.SendText("现在是机器人自动给你发的问候~晚安~")
}
func sendImageToSomeOne(f *openwechat.Friend, path string) {
	//for {
	//	now := time.Now()
	//	t := time.Date(now.Year(), now.Month(), now.Day(), 6, 20, 0, 0, now.Location())
	//	timer := time.NewTimer(now.Sub(t))
	//	<-timer.C
	//	f.SendText("现在是机器人自动给你发的问候~晚安~")
	//	break
	//}
	img, _ := os.Open(path)
	defer img.Close()
	f.SendImage(img)
}
func Run() {
	//bot := openwechat.DefaultBot()
	bot := openwechat.DefaultBot(openwechat.Desktop) // 桌面模式，上面登录不上的可以尝试切换这种模式

	// 注册消息处理函数
	handler, err := handlers.NewHandler()
	if err != nil {
		logger.Danger("register error: %v", err)
		return
	}
	bot.MessageHandler = handler

	// 注册登陆二维码回调
	bot.UUIDCallback = handlers.QrCodeCallBack

	// 创建热存储容器对象
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")

	// 执行热登录
	err = bot.HotLogin(reloadStorage, true)
	if err != nil {
		logger.Warning(fmt.Sprintf("login error: %v ", err))
		return
	}
	self, err := bot.GetCurrentUser()
	if err != nil {
		logger.Warning(fmt.Sprintf("login error: %v ", err))
		return
	}
	fmt.Println(self.ID())
	fmt.Println(self.NickName)
	fmt.Println(self.RemarkName)
	fmt.Println(self.UserName)
	fmt.Println(self.HeadImgUrl)
	//friends, err := self.Friends()
	//if err != nil {
	//	logger.Warning(fmt.Sprintf("login error: %v ", err))
	//	return
	//}
	//
	//someone := friends.SearchByNickName(2, "陈小瑞")
	//logger.Warning(fmt.Sprintf("查找到的好友数: %d ", someone.Count()))
	//if someone.Count() > 0 {
	//	//go sendMessageToSomeOne(someone.First())
	//	path, _ := gpt.Text2img("流星雨")
	//	go sendImageToSomeOne(someone.First(), path)
	//}
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}
