package gpt

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/config"
	"gitee.com/lmuiotctf/chatGpt_wechat/tree/master/pkg/logger"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// const BASEURL = "https://api.openai.com/v1/"
const BASEURL = "https://api.openai.com/v1/"

/*
{
"discord":{
"id":"1111158971491954689",
"channel_id":"1110596223653654528",
"guild_id":"1110591114882322512",
"content":"**请画一幅流星雨？ --v 5.1 raw --s 250** - \u003c@1097461907797057676\u003e (fast)",
"timestamp":"2023-05-25T05:09:03.139Z",
"edited_timestamp":null,
"mention_roles":[
],
"tts":false,
"mention_everyone":false,
"author":{
"id":"936929561302675456",
"email":"",
"username":"Midjourney Bot",
"avatar":"4a79ea7cd151474ff9f6e08339d69380",
"locale":"",
"discriminator":"9282",
"token":"",
"verified":false,
"mfa_enabled":false,
"banner":"",
"accent_color":0,
"bot":true,
"public_flags":589824,
"premium_type":0,
"system":false,
"flags":0
},
"attachments":[
{
"id":"1111158970833457322",
"url":"https://cdn.discordapp.com/attachments/1110596223653654528/1111158970833457322/HelenAnderson__bf52ee6c-6c9a-42c6-af3b-f77e67ae4fa6.png",
"proxy_url":"https://media.discordapp.net/attachments/1110596223653654528/1111158970833457322/HelenAnderson__bf52ee6c-6c9a-42c6-af3b-f77e67ae4fa6.png",
"filename":"HelenAnderson__bf52ee6c-6c9a-42c6-af3b-f77e67ae4fa6.png",
"content_type":"image/png",
"width":2048,
"height":2048,
"size":7250543,
"ephemeral":false
}
],
"embeds":[
],
"mentions":[
{
"id":"1097461907797057676",
"email":"",
"username":"HelenAnderson",
"avatar":"5091e12197877c3ab03667bbcaa26e2f",
"locale":"",
"discriminator":"1388",
"token":"",
"verified":false,
"mfa_enabled":false,
"banner":"",
"accent_color":0,
"bot":false,
"public_flags":0,
"premium_type":0,
"system":false,
"flags":0
}
],
"reactions":null,
"pinned":false,
"type":0,
"webhook_id":"",
"member":{
"guild_id":"",
"joined_at":"2023-05-23T15:38:59.07952Z",
"nick":"",
"deaf":false,
"mute":false,
"avatar":"",
"user":null,
"roles":[
"1110592723217559666"
],
"premium_since":null,
"pending":false,
"permissions":"0",
"communication_disabled_until":null
},
"mention_channels":null,
"activity":null,
"application":null,
"message_reference":null,
"referenced_message":null,
"interaction":null,
"flags":0,
"sticker_items":null
},
"type":"GenerateEnd"
}
{
    "discord": {
        "id": "1111158971491954689",
        "channel_id": "1110596223653654528",
        "guild_id": "1110591114882322512",
        "content": "**请画一幅流星雨？ --v 5.1 raw --s 250** - \u003c@1097461907797057676\u003e (fast)",
        "timestamp": "2023-05-25T05:09:03.139Z",
        "edited_timestamp": null,
        "mention_roles": [

        ],
        "tts": false,
        "mention_everyone": false,
        "author": {
            "id": "936929561302675456",
            "email": "",
            "username": "Midjourney Bot",
            "avatar": "4a79ea7cd151474ff9f6e08339d69380",
            "locale": "",
            "discriminator": "9282",
            "token": "",
            "verified": false,
            "mfa_enabled": false,
            "banner": "",
            "accent_color": 0,
            "bot": true,
            "public_flags": 589824,
            "premium_type": 0,
            "system": false,
            "flags": 0
        },
        "attachments": [
            {
                "id": "1111158970833457322",
                "url": "https://cdn.discordapp.com/attachments/1110596223653654528/1111158970833457322/HelenAnderson__bf52ee6c-6c9a-42c6-af3b-f77e67ae4fa6.png",
                "proxy_url": "https://media.discordapp.net/attachments/1110596223653654528/1111158970833457322/HelenAnderson__bf52ee6c-6c9a-42c6-af3b-f77e67ae4fa6.png",
                "filename": "HelenAnderson__bf52ee6c-6c9a-42c6-af3b-f77e67ae4fa6.png",
                "content_type": "image/png",
                "width": 2048,
                "height": 2048,
                "size": 7250543,
                "ephemeral": false
            }
        ],
        "embeds": [

        ],
        "mentions": [
            {
                "id": "1097461907797057676",
                "email": "",
                "username": "HelenAnderson",
                "avatar": "5091e12197877c3ab03667bbcaa26e2f",
                "locale": "",
                "discriminator": "1388",
                "token": "",
                "verified": false,
                "mfa_enabled": false,
                "banner": "",
                "accent_color": 0,
                "bot": false,
                "public_flags": 0,
                "premium_type": 0,
                "system": false,
                "flags": 0
            }
        ],
        "reactions": null,
        "pinned": false,
        "type": 0,
        "webhook_id": "",
        "member": {
            "guild_id": "",
            "joined_at": "2023-05-23T15:38:59.07952Z",
            "nick": "",
            "deaf": false,
            "mute": false,
            "avatar": "",
            "user": null,
            "roles": [
                "1110592723217559666"
            ],
            "premium_since": null,
            "pending": false,
            "permissions": "0",
            "communication_disabled_until": null
        },
        "mention_channels": null,
        "activity": null,
        "application": null,
        "message_reference": null,
        "referenced_message": null,
        "interaction": null,
        "flags": 0,
        "sticker_items": null
    },
    "type": "GenerateEnd"
}
*/

// ChatGPTResponseBody 响应
type ChatGPTResponseBody struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int                    `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChoiceItem           `json:"choices"`
	Usage   map[string]interface{} `json:"usage"`
}

type ChoiceItem struct {
	Message      message `json:"message"`
	Index        int     `json:"index"`
	Logprobs     int     `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}
type ChatGPTImgResponseBody struct {
	Created int                    `json:"created"`
	Data    []DataItem             `json:"data"`
	Usage   map[string]interface{} `json:"usage"`
}
type MjbotResponseBody struct {
	Message string `json:"message"`
}
type MjcallbackResponseBody struct {
	Discord struct {
		ID              string        `json:"id"`
		ChannelID       string        `json:"channel_id"`
		GuildID         string        `json:"guild_id"`
		Content         string        `json:"content"`
		Timestamp       time.Time     `json:"timestamp"`
		EditedTimestamp interface{}   `json:"edited_timestamp"`
		MentionRoles    []interface{} `json:"mention_roles"`
		Tts             bool          `json:"tts"`
		MentionEveryone bool          `json:"mention_everyone"`
		Author          struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			Username      string `json:"username"`
			Avatar        string `json:"avatar"`
			Locale        string `json:"locale"`
			Discriminator string `json:"discriminator"`
			Token         string `json:"token"`
			Verified      bool   `json:"verified"`
			MfaEnabled    bool   `json:"mfa_enabled"`
			Banner        string `json:"banner"`
			AccentColor   int    `json:"accent_color"`
			Bot           bool   `json:"bot"`
			PublicFlags   int    `json:"public_flags"`
			PremiumType   int    `json:"premium_type"`
			System        bool   `json:"system"`
			Flags         int    `json:"flags"`
		} `json:"author"`
		Attachments []struct {
			Id          string `json:"id"`
			Url         string `json:"url"`
			ProxyURL    string `json:"proxy_url"`
			Filename    string `json:"filename"`
			ContentType string `json:"content_type"`
			Width       int    `json:"width"`
			Height      int    `json:"height"`
			Size        int    `json:"size"`
			Ephemeral   bool   `json:"ephemeral"`
		} `json:"attachments"`
		Embeds   []interface{} `json:"embeds"`
		Mentions []struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			Username      string `json:"username"`
			Avatar        string `json:"avatar"`
			Locale        string `json:"locale"`
			Discriminator string `json:"discriminator"`
			Token         string `json:"token"`
			Verified      bool   `json:"verified"`
			MfaEnabled    bool   `json:"mfa_enabled"`
			Banner        string `json:"banner"`
			AccentColor   int    `json:"accent_color"`
			Bot           bool   `json:"bot"`
			PublicFlags   int    `json:"public_flags"`
			PremiumType   int    `json:"premium_type"`
			System        bool   `json:"system"`
			Flags         int    `json:"flags"`
		} `json:"mentions"`
		Reactions interface{} `json:"reactions"`
		Pinned    bool        `json:"pinned"`
		Type      int         `json:"type"`
		WebhookID string      `json:"webhook_id"`
		Member    struct {
			GuildID                    string      `json:"guild_id"`
			JoinedAt                   time.Time   `json:"joined_at"`
			Nick                       string      `json:"nick"`
			Deaf                       bool        `json:"deaf"`
			Mute                       bool        `json:"mute"`
			Avatar                     string      `json:"avatar"`
			User                       interface{} `json:"user"`
			Roles                      []string    `json:"roles"`
			PremiumSince               interface{} `json:"premium_since"`
			Pending                    bool        `json:"pending"`
			Permissions                string      `json:"permissions"`
			CommunicationDisabledUntil interface{} `json:"communication_disabled_until"`
		} `json:"member"`
		MentionChannels   interface{} `json:"mention_channels"`
		Activity          interface{} `json:"activity"`
		Application       interface{} `json:"application"`
		MessageReference  interface{} `json:"message_reference"`
		ReferencedMessage interface{} `json:"referenced_message"`
		Interaction       interface{} `json:"interaction"`
		Flags             int         `json:"flags"`
		StickerItems      interface{} `json:"sticker_items"`
	} `json:"discord"`
	Type string `json:"type"`
}
type DataItem struct {
	Url string `json:"url"`
}
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTRequestBody 请求体
type ChatGPTRequestBody struct {
	Model            string     `json:"model"`
	Messages         []Messages `json:"messages"`
	MaxTokens        uint       `json:"max_tokens"`
	Temperature      float64    `json:"temperature"`
	TopP             int        `json:"top_p"`
	FrequencyPenalty int        `json:"frequency_penalty"`
	PresencePenalty  int        `json:"presence_penalty"`
}
type RequestBodyImg struct {
	Size   string `json:"size"`
	Prompt string `json:"prompt"`
}
type RequestBodyImgMJ struct {
	Type         string `json:"type"`
	Prompt       string `json:"prompt"`
	DiscordMsgId string `json:"discordMsgId"`
	MsgHash      string `json:"msgHash"`
	Index        int64  `json:"index"`
}
type Messages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Completions gtp文本模型回复
// curl https://api.openai.com/v1/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer your chatGPT key"
// -d '{"model": "text-davinci-003", "prompt": "give me good song", "temperature": 0, "max_tokens": 7}'
func Completions(msg string) (string, error) {
	cfg := config.LoadConfig()

	requestBody := ChatGPTRequestBody{
		Model: cfg.Model,

		MaxTokens: cfg.MaxTokens,
		Messages:  []Messages{{"user", msg}},

		Temperature:      cfg.Temperature,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gpt json string : %v", string(requestData)))
	req, err := http.NewRequest("POST", BASEURL+"chat/completions", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 300 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return "", errors.New(fmt.Sprintf("请求GTP出错了，gpt api status code not equals 200,code is %d ,details:  %v ", response.StatusCode, string(body)))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("response gpt json string : %v", string(body)))

	gptResponseBody := &ChatGPTResponseBody{}
	log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}

	var reply string
	if len(gptResponseBody.Choices) > 0 {
		reply = gptResponseBody.Choices[0].Message.Content
	}
	logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
	return reply, nil
}
func Text2imgGpt(msg string) (string, error) {
	size := "512x512"
	requestBody := RequestBodyImg{
		Size:   size,
		Prompt: msg,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gpt json string : %v", string(requestData)))
	req, err := http.NewRequest("POST", BASEURL+"images/generations", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}

	apiKey := config.LoadConfig().ApiKey
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{Timeout: 300 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return "", errors.New(fmt.Sprintf("请求GTP出错了，gpt api status code not equals 200,code is %d ,details:  %v ", response.StatusCode, string(body)))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	//logger.Info(fmt.Sprintf("response gpt json string : %v", string(body)))
	//fmt.Println("Response body:", string(body))
	//
	gptResponseBody := &ChatGPTImgResponseBody{}
	//log.Println(string(body))
	err = json.Unmarshal(body, gptResponseBody)
	if err != nil {
		return "", err
	}
	//
	var reply string
	if len(gptResponseBody.Data) > 0 {
		reply = gptResponseBody.Data[0].Url
	}
	logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
	// Get the data
	resp, err := http.Get(reply)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//currentPath := GetCurrentDirectory()
	// 创建一个文件用于保存
	tmpData := []byte(gptResponseBody.Data[0].Url)
	sum := md5.Sum(tmpData)
	name := hex.EncodeToString(sum[:])
	out, err := os.Create("/Users/wlb/" + name + ".jpg")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	return "/Users/wlb/" + name + ".jpg", nil
}
func Text2img(msg string, ch chan string) (string, error) {
	requestBody := RequestBodyImgMJ{
		Type:   "generate",
		Prompt: msg,
	}
	requestData, err := json.Marshal(requestBody)

	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("request gpt json string : %v", string(requestData)))
	req, err := http.NewRequest("POST", "http://127.0.0.1:16007/v1/trigger/midjourney-bot", bytes.NewBuffer(requestData))
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 300 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		body, _ := ioutil.ReadAll(response.Body)
		return "", errors.New(fmt.Sprintf("请求出错了，mj api status code not equals 200,code is %d ,details:  %v ", response.StatusCode, string(body)))
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	logger.Info(fmt.Sprintf("response gpt json string : %v", string(body)))
	fmt.Println("Response body:", string(body))

	mjResponseBody := &MjbotResponseBody{}
	//log.Println(string(body))
	err = json.Unmarshal(body, mjResponseBody)
	if err != nil {
		return "", err
	}
	//
	if mjResponseBody.Message == "success" {
		//reply = gptResponseBody.Data[0].Url
		select {
		case res := <-ch:
			fmt.Println("async callback is done:", res)
			mjcallbackResponseBody := &MjcallbackResponseBody{}
			err = json.Unmarshal([]byte(res), mjcallbackResponseBody)
			if err != nil {
				return "", nil
			}

			resp, err := http.Get(mjcallbackResponseBody.Discord.Attachments[0].Url)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			// 创建一个文件用于保存
			tmpData := []byte(mjcallbackResponseBody.Discord.Attachments[0].Url)
			sum := md5.Sum(tmpData)
			name := hex.EncodeToString(sum[:])
			out, err := os.Create("/Users/wlb/" + name + ".jpg")
			if err != nil {
				panic(err)
			}
			defer out.Close()

			// 然后将响应流和文件流对接起来
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				panic(err)
			}
			return "/Users/wlb/" + name + ".jpg", nil
		case <-time.After(600 * time.Second):
			fmt.Println("超时")
			return "", nil

		}
	}
	return "", nil
	//logger.Info(fmt.Sprintf("gpt response text: %s ", reply))
	//resp, err := http.Get(reply)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//// 创建一个文件用于保存
	//tmpData := []byte(gptResponseBody.Data[0].Url)
	//sum := md5.Sum(tmpData)
	//name := hex.EncodeToString(sum[:])
	//out, err := os.Create("/Users/wlb/" + name + ".jpg")
	//if err != nil {
	//	panic(err)
	//}
	//defer out.Close()
	//
	//// 然后将响应流和文件流对接起来
	//_, err = io.Copy(out, resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//return "/Users/wlb/" + name + ".jpg", nil
}

//func GetCurrentDirectory() string {
//	//返回绝对路径 filepath.Dir(os.Args[0])去除最后一个元素的路径
//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	//将\替换成/
//	return strings.Replace(dir, "\\", "/", -1)
//}
