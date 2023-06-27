package handlers

import (
	"fmt"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/config"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/pkg/logger"
	"github.com/eatmoreapple/openwechat"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var c = cache.New(config.LoadConfig().SessionTimeout, time.Minute*5)

// MessageHandlerInterface 消息处理接口
type MessageHandlerInterface interface {
	handle() error
	ReplyText() error
}

// 声明全局channel
var Ch chan string

// 声明处理回调的http服务器
func listen() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("被调用了")
		//读取回调数据
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		if !strings.Contains(string(data), "FirstTrigger") {
			fmt.Println("AI生成图片成功")
			Ch <- string(data)
			w.Write([]byte("finish"))
		}

	})
	fmt.Println("Server is listening on port 8888")
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

}

// QrCodeCallBack 登录扫码回调，
func QrCodeCallBack(uuid string) {
	if runtime.GOOS == "windows" {
		// 运行在Windows系统上
		openwechat.PrintlnQrcodeUrl(uuid)
	} else {
		openwechat.PrintlnQrcodeUrl(uuid)

		//log.Println("login in linux")
		//url := "https://login.weixin.qq.com/l/" + uuid
		//log.Printf(url)
		//log.Printf("如果二维码无法扫描，请缩小控制台尺寸，或更换命令行工具，缩小二维码像素")
		//q, _ := qrcode.New(url, qrcode.High)
		//fmt.Println(q.ToSmallString(true))
	}
}

func NewHandler() (msgFunc func(msg *openwechat.Message), err error) {
	Ch = make(chan string)
	go listen()
	dispatcher := openwechat.NewMessageMatchDispatcher()

	// 清空会话
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return strings.Contains(message.Content, config.LoadConfig().SessionClearToken)
	}, TokenMessageContextHandler())

	// 处理群消息
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsSendByGroup()
	}, GroupMessageContextHandler())

	// 好友申请
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return message.IsFriendAdd()
	}, func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		if config.LoadConfig().AutoPass {
			_, err := msg.Agree("")
			if err != nil {
				logger.Warning(fmt.Sprintf("add friend agree error : %v", err))
				return
			}
		}
	})

	// 私聊
	// 获取用户消息处理器
	dispatcher.RegisterHandler(func(message *openwechat.Message) bool {
		return !(strings.Contains(message.Content, config.LoadConfig().SessionClearToken) || message.IsSendByGroup() || message.IsFriendAdd())
	}, UserMessageContextHandler())
	return openwechat.DispatchMessage(dispatcher), nil
}
