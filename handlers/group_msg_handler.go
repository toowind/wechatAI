package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/gpt"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/pkg/logger"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/service"
	"github.com/eatmoreapple/openwechat"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var _ MessageHandlerInterface = (*GroupMessageHandler)(nil)

type ReqMessage struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	To      string `json:"to"`
	Content string `json:"content"`
}

type ResponseMessage struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Plugin  string     `json:"plugin"`
	Data    []Dataitem `json:"data"`
}
type Dataitem struct {
	Content string `json:"content"`
}

// GroupMessageHandler 群消息处理
type GroupMessageHandler struct {
	// 获取自己
	self *openwechat.Self
	// 群
	group *openwechat.Group
	// 接收到消息
	msg *openwechat.Message
	// 发送的用户
	sender *openwechat.User
	// 实现的用户业务
	service service.UserServiceInterface
}

func GroupMessageContextHandler() func(ctx *openwechat.MessageContext) {
	return func(ctx *openwechat.MessageContext) {
		msg := ctx.Message
		// 获取用户消息处理器
		handler, err := NewGroupMessageHandler(msg)
		if err != nil {
			logger.Warning(fmt.Sprintf("init group message handler error: %s", err))
			//if "can not found sender from system message" == err.Error() {
			//	msg.ReplyText("大家好，我是AI机器人，大家有什么问题都可以问我。另外，如果需要我作画，请@我，然后说 请画xxx")
			//}
			return
		}

		// 处理用户消息
		err = handler.handle()
		if err != nil {
			logger.Warning(fmt.Sprintf("handle group message error: %s", err))
		}
	}
}

// NewGroupMessageHandler 创建群消息处理器
func NewGroupMessageHandler(msg *openwechat.Message) (MessageHandlerInterface, error) {
	sender, err := msg.Sender()
	if err != nil {
		return nil, err
	}
	group := &openwechat.Group{User: sender}
	groupSender, err := msg.SenderInGroup()
	if err != nil {
		return nil, err
	}

	userService := service.NewUserService(c, groupSender)
	handler := &GroupMessageHandler{
		self:    sender.Self,
		msg:     msg,
		group:   group,
		sender:  groupSender,
		service: userService,
	}
	return handler, nil

}

// handle 处理消息
func (g *GroupMessageHandler) handle() error {
	if g.msg.IsText() {
		return g.ReplyText()
	}
	return nil
}

// ReplyText 发送文本消息到群
func (g *GroupMessageHandler) ReplyText() error {
	if "778131193" != g.group.ID() && "778132299" != g.group.ID() {
		logger.Info(fmt.Sprintf("Received Group %v Text Msg : %v", g.group.NickName, g.msg.Content))
	}
	// 1.不是@的不处理
	if !g.msg.IsAt() {
		return nil
	}

	// 2.获取请求的文本，如果为空字符串不处理
	requestText := g.getRequestText()
	if requestText == "" {
		logger.Info("user message is null")
		return nil
	}
	logger.Info(fmt.Sprintf("requestText: %v", requestText))
	match := strings.Contains(requestText, "请画")
	wmatch := strings.Contains(requestText, "天气")
	jokematch := strings.Contains(requestText, "笑话")

	if match == true {
		g.service.ClearUserSessionContext()
		g.msg.ReplyText(g.buildReplyText("图片正在生成中，请耐心等待......"))
		// 3.请求GPT获取回复
		reply, err := gpt.Text2img(requestText, Ch)
		if err != nil {
			// 2.1 将GPT请求失败信息输出给用户，省得整天来问又不知道日志在哪里。
			errMsg := fmt.Sprintf("gpt request error: %v", err)
			_, err = g.msg.ReplyText(errMsg)
			//_, err = g.msg.ReplyText("网络问题，请稍后重试")
			if err != nil {
				return errors.New(fmt.Sprintf("response group error: %v ", err))
			}
			return err
		}
		// 4.设置上下文，并响应信息给用户
		img, _ := os.Open(reply)
		defer img.Close()
		_, err = g.msg.ReplyImage(img)
		if err != nil {
			return errors.New(fmt.Sprintf("response user error: %v ", err))
		}

		// 5.返回错误信息
		return err
	} else {
		if wmatch == true || jokematch == true {
			apiURL := "https://api.mlyai.com/reply"
			apiKey := "hwdzm95u6acos8nw"
			apiSecret := "untdwlsh"

			// 构造 HTTP 请求参数
			//reqBody := []byte(`{
			//"type": "2",
			//"from": "2",
			//"to": "2",
			//"content":"讲个笑话吧"
			//}`)
			reqBody := ReqMessage{
				Type:    "2",
				From:    "2",
				To:      "2",
				Content: requestText,
			}
			requestData, _ := json.Marshal(reqBody)
			req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestData))
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Api-Key", apiKey)
			req.Header.Set("Api-Secret", apiSecret)

			// 发送 HTTP 请求并获取响应
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}
			defer resp.Body.Close()
			// 读取响应数据
			respBodyData, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}
			responseMessage := &ResponseMessage{}
			err = json.Unmarshal(respBodyData, responseMessage)
			if err != nil {
				fmt.Println("Error: ", err)
				return err
			}
			var rep string
			if len(responseMessage.Data) > 0 {
				rep = responseMessage.Data[0].Content
			}
			// 输出响应状态码和响应数据
			fmt.Println("Status code:", resp.StatusCode)
			fmt.Println("Response body:", string(respBodyData))

			_, err = g.msg.ReplyText(g.buildReplyText(rep))
			if err != nil {
				return errors.New(fmt.Sprintf("response user error: %v ", err))
			}

			// 5.返回错误信息
			return err

		} else {
			// 3.请求GPT获取回复
			reply, err := gpt.Completions(requestText)
			if err != nil {
				// 2.1 将GPT请求失败信息输出给用户，省得整天来问又不知道日志在哪里。
				errMsg := fmt.Sprintf("gpt request error: %v", err)
				//_, err = g.msg.ReplyText("网络问题，请稍后重试")
				_, err = g.msg.ReplyText(errMsg)
				if err != nil {
					return errors.New(fmt.Sprintf("response group error: %v ", err))
				}
				return err
			}

			// 4.设置上下文，并响应信息给用户
			g.service.SetUserSessionContext(requestText, reply)
			_, err = g.msg.ReplyText(g.buildReplyText(reply))
			if err != nil {
				return errors.New(fmt.Sprintf("response user error: %v ", err))
			}

			// 5.返回错误信息
			return err
		}
	}
	return nil
}

// getRequestText 获取请求接口的文本，要做一些清洗
func (g *GroupMessageHandler) getRequestText() string {
	// 1.去除空格以及换行
	requestText := strings.TrimSpace(g.msg.Content)
	requestText = strings.Trim(g.msg.Content, "\n")

	// 2.替换掉当前用户名称
	replaceText := "@" + g.self.NickName
	requestText = strings.TrimSpace(strings.ReplaceAll(g.msg.Content, replaceText, ""))
	if requestText == "" {
		return ""
	}

	// 3.获取上下文，拼接在一起，如果字符长度超出4000，截取为4000。（GPT按字符长度算），达芬奇3最大为4068，也许后续为了适应要动态进行判断。
	sessionText := g.service.GetUserSessionContext()
	if sessionText != "" {
		requestText = sessionText + "\n" + requestText
	}
	if len(requestText) >= 4000 {
		requestText = requestText[:4000]
	}

	// 4.检查用户发送文本是否包含结束标点符号
	punctuation := ",.;!?，。！？、…"
	runeRequestText := []rune(requestText)
	lastChar := string(runeRequestText[len(runeRequestText)-1:])
	if strings.Index(punctuation, lastChar) < 0 {
		requestText = requestText + "。" // 判断最后字符是否加了标点，没有的话加上句号，避免openai自动补齐引起混乱。
	}

	// 5.返回请求文本
	return requestText
}

// buildReply 构建回复文本
func (g *GroupMessageHandler) buildReplyText(reply string) string {
	// 1.获取@我的用户
	atText := "@" + g.sender.NickName
	textSplit := strings.Split(reply, "\n\n")
	if len(textSplit) > 1 {
		trimText := textSplit[0]
		reply = strings.Trim(reply, trimText)
	}
	reply = strings.TrimSpace(reply)
	if reply == "" {
		return atText + " 请求得不到任何有意义的回复，请具体提出问题。"
	}

	// 2.拼接回复,@我的用户，问题，回复
	//replaceText := "@" + g.self.NickName
	//question := strings.TrimSpace(strings.ReplaceAll(g.msg.Content, replaceText, ""))
	//reply = atText + "\n" + question + "\n" + reply
	reply = atText + "\n" + reply

	reply = strings.Trim(reply, "\n")

	// 3.返回回复的内容
	return reply
}
