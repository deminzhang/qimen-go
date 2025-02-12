package world

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/deminzhang/qimen-go/gui"
	"io"
	"net/http"
	"strings"
	"sync"
)

type UIChat struct {
	gui.BaseUI
	textBoxLog *gui.TextBox
	showChatUI bool
}

var uiChat *UIChat

func UIShowChat() {
	if uiChat == nil {
		uiChat = NewUIChat()
		gui.ActiveUI(uiChat)
	}
}
func UIHideChat() {
	if uiChat != nil {
		gui.CloseUI(uiChat)
		uiChat = nil
	}
}
func UIChatLog(msg string, a ...any) {
	if uiChat != nil {
		uiChat.ChatLog(msg, a...)
	}
}

func NewUIChat() *UIChat {
	p := &UIChat{
		BaseUI: gui.BaseUI{Visible: true,
			X: 0, Y: ScreenHeight - 250, W: 350, H: 250,
		},
		showChatUI: true,
	}
	textBoxLog := gui.NewTextBox(16, 0, 270, 180)
	inputBoxChat := gui.NewInputBox(16, 190, 270, 32)

	checkBoxGM := gui.NewCheckBox(290, 200, "令")
	btnChatSend := gui.NewButton(290, 230, 32, 16, "发")

	btnChatSwitch := gui.NewButton(0, 230, 32, 16, "隐")
	btnClear := gui.NewButton(32, 230, 32, 16, "清")
	btnCopy := gui.NewButton(64, 230, 32, 16, "复")

	p.AddChildren(btnChatSwitch, btnClear, btnCopy,
		textBoxLog,
		inputBoxChat,
		checkBoxGM,
		btnChatSend)

	inputBoxChat.DefaultText = "输入信息指令..[help查看命令]"

	inputBoxChat.SetOnPressEnter(func(i *gui.InputBox) {
		if !p.showChatUI {
			return
		}
		if i.Focused() {
			btnChatSend.Click()
		} else {
			i.SetFocused(true)
		}
	})
	checkBoxGM.SetOnCheckChanged(func(c *gui.CheckBox) {
		msg := "debug command"
		if c.Checked() {
			msg += " (On)"
		} else {
			msg += " (Off)"
		}
		textBoxLog.AppendLine(msg)
	})
	checkBoxGM.SetChecked(true)
	btnChatSend.SetOnClick(func() {
		i := inputBoxChat
		if i.Text() != "" {
			i.SetFocused(false)
			if checkBoxGM.Checked() { //调试命令
				textBoxLog.AppendLine("cmd: " + i.Text())
				parseCmd(i.Text())
			} else { //聊天
				textBoxLog.AppendLine("me:" + i.Text())
				go sendChat(i.Text())
			}
			i.AppendTextHistory(i.Text())
			i.SetText("")
		} else {
			i.SetFocused(false)
			textBoxLog.AppendLine("no input msg")
		}
	})
	btnChatSwitch.SetOnClick(func() {
		p.showChatUI = !p.showChatUI
		if p.showChatUI {
			btnChatSwitch.Text = "隐"
		} else {
			btnChatSwitch.Text = "令"
		}
		textBoxLog.Visible = p.showChatUI
		inputBoxChat.Visible = p.showChatUI
		btnChatSend.Visible = p.showChatUI
		checkBoxGM.Visible = p.showChatUI
		btnClear.Visible = p.showChatUI
		btnCopy.Visible = p.showChatUI
	})
	btnClear.SetOnClick(func() {
		textBoxLog.Text = ""
	})
	btnCopy.SetOnClick(func() {
		clipboard.WriteAll(textBoxLog.Text)
	})

	p.textBoxLog = textBoxLog
	uiChat = p
	return p
}

func (p *UIChat) Update() {
	p.BaseUI.Update()
}

func (p *UIChat) OnClose() {
	uiChat = nil
}

func (p *UIChat) ChatLog(msg string, a ...any) {
	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a...)
	}
	p.textBoxLog.AppendLine(msg)
}

func parseCmd(str string) {
	arr := strings.Split(str, " ")
	var args []string
	for _, ss := range arr {
		if len(strings.TrimSpace(ss)) > 0 {
			args = append(args, strings.TrimSpace(ss))
		}
	}
	switch strings.ToLower(args[0]) {
	case strings.ToLower("help"):
		UIChatLog("help: show help")
		UIChatLog("showBattle: 显战斗")
		UIChatLog("hideBattle: 隐战斗")
		UIChatLog("showWave: 显示引力波")
		UIChatLog("hideWave: 隐引力波")
	case strings.ToLower("showBattle"):
		ThisGame.showBattle = true
	case strings.ToLower("hideBattle"):
		ThisGame.showBattle = false
	case strings.ToLower("showWave"):
		ThisGame.showWave = true
	case strings.ToLower("hideWave"):
		ThisGame.showWave = false
	}
}

const (
	//官方 DeepSeek
	// https://api-docs.deepseek.com/zh-cn/
	DeepSeekChatURL = "https://api.deepseek.com/chat/completions"
	DeepSeekAPIKey  = "sk-11ec325f6c164193a0bc3a98b8c56fe3" // Replace with your API key

	//本地 Ollama
	//https://github.com/datawhalechina/handy-ollama/blob/main/docs/C4/1.%20Ollama%20API%20%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97.md
	ChatURL = "http://10.100.136.238:11434/v1/chat/completions"
	APIKey  = "" // Replace with your API key
)

func sendChat(str string) {
	UIChatLog("chat: " + str)
	response, err := sendDeepSeekRequest(ChatURL, APIKey, str)
	if err != nil {
		UIChatLog("error: " + err.Error())
		println("error: " + err.Error())
		return
	}
	res := string(response)
	UIChatLog("response: " + res)
	println("response: " + res)
	//{"error":{"message":"Authentication Fails (no such user)","type":"authentication_error","param":null,"code":"invalid_request_error"}} //认证失败
	//{"error":{"message":"Insufficient Balance","type":"unknown_error","param":null,"code":"invalid_request_error"}} //余额不足

	//{"error":{"message":"model \"deepseek-chat\" not found, try pulling it first","type":"api_error","param":null,"code":null}}
	// {"id":"chatcmpl-386","object":"chat.completion","created":1739351889,"model":"deepseek-r1:7b","system_fingerprint":"fp_ollama","choices":[{"index":0,"message":{"role":"assistant","content":"\u003cthink\u003e\nOkay, the user says \"Hello!\" and that's it. I should respond politely.\n\nI'll greet them back in a friendly manner.\n\nMaybe ask how they can help today to keep the conversation going.\n\nKeep it simple and open-ended so they feel comfortable sharing more.\n\u003c/think\u003e\n\nHello! How can I assist you today?"},"finish_reason":"stop"}],"usage":{"prompt_tokens":11,"completion_tokens":68,"total_tokens":79}}

	//TODO: 解析返回结果
	var data map[string]interface{}
	err = json.Unmarshal(response, &data)
	if err != nil {
		UIChatLog("error: " + err.Error())
		println("error: " + err.Error())
		return
	}
	if data["error"] != nil {
		UIChatLog("error: " + data["error"].(map[string]interface{})["message"].(string))
		println("error: " + data["error"].(map[string]interface{})["message"].(string))
		return
	}
	choices := data["choices"].([]interface{})
	for _, choice := range choices {
		msg := choice.(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		UIChatLog("assistant: " + msg)
		payload["messages"] = append(payload["messages"].([]map[string]string), map[string]string{"role": "system", "content": msg})
	}
}

var payload = map[string]interface{}{
	//"model": "deepseek-chat",
	"model": "deepseek-r1:7b",
	"messages": []map[string]string{
		{"role": "system", "content": "You are a helpful assistant."},
		//{"role": "user", "content": "Hello!"},
	},
	//"stream": false, //非流式 一次性返回
	"stream": true, //流式
}

func sendDeepSeekRequest(url, apikey, str string) ([]byte, error) {
	his := payload["messages"]
	his = append(his.([]map[string]string), map[string]string{"role": "user", "content": str})

	// Marshal the payload into JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	if payload["stream"].(bool) { //TODO
		var wg sync.WaitGroup
		done := make(chan struct{})
		// 处理流式响应
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer resp.Body.Close()
			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println("Error reading response:", err)
					return
				}
				fmt.Print(line)
				// data: {"id":"chatcmpl-225","object":"chat.completion.chunk","created":1739356082,"model":"deepseek-r1:7b","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":"?"},"finish_reason":null}]}
				//	data: {"id":"chatcmpl-225","object":"chat.completion.chunk","created":1739356082,"model":"deepseek-r1:7b","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":"stop"}]}
				//	data: [DONE]
			}
		}()
		// 等待流结束
		select {
		case <-done:
			fmt.Println("Stream ended or closed")
		}
	}

	defer resp.Body.Close()
	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return responseBody, nil
}
