package world

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	//DeepSeekChatURL https://api-docs.deepseek.com/zh-cn/
	DeepSeekChatURL = "https://api.deepseek.com/chat/completions" //官方 DeepSeek
	DeepSeekAPIKey  = ""                                          // Replace with your API key

	//ChatURL
	//https://github.com/datawhalechina/handy-ollama/blob/main/docs/C4/1.%20Ollama%20API%20%E4%BD%BF%E7%94%A8%E6%8C%87%E5%8D%97.md
	ChatURL = "http://10.100.136.238:11434/v1/chat/completions" //本地 Ollama
	APIKey  = ""                                                // Replace with your API key
)

type Chat struct {
	Payload map[string]interface{}
	outFunc func(fmt string, a ...any)
}

func NewChat(model string, outFunc func(fmt string, a ...any)) *Chat {
	return &Chat{
		Payload: map[string]interface{}{
			"model":    model,
			"messages": []map[string]string{},
			//"stream": false, //非流式 一次性返回
			"stream": true, //流式
		},
		outFunc: outFunc,
	}
}

func (c *Chat) SetModel(str string) {
	c.Payload["model"] = str
}

func (c *Chat) ApplyChat(role, content string) {
	c.Payload["messages"] = append(c.Payload["messages"].([]map[string]string), map[string]string{"role": role, "content": content})
}

func (c *Chat) SendChat(role, content string) {
	outFunc := c.outFunc
	c.ApplyChat(role, content)
	err := c.sendAIRequest(ChatURL, APIKey)
	if err != nil {
		outFunc("error: %s\n", err.Error())
		println("error: " + err.Error())
	}
}

func (c *Chat) sendAIRequest(url, apikey string) error {
	payload := c.Payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if payload["stream"].(bool) {
		return c.handleStreamResponse(resp)
	} else {
		return c.handleNonStreamResponse(resp)
	}
}

func (c *Chat) handleStreamResponse(resp *http.Response) error {
	payload := c.Payload
	outFunc := c.outFunc
	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
			err := resp.Body.Close()
			if err != nil {
				fmt.Println("Error closing response body:", err)
				return
			}
		}()
		var role string
		var contentX string
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
			if line == "\n" {
				continue
			}
			//example
			//  data: {"id":"chatcmpl-225","object":"chat.completion.chunk","created":1739356082,"model":"deepseek-r1:7b","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":"?"},"finish_reason":null}]}
			//	data: {"id":"chatcmpl-225","object":"chat.completion.chunk","created":1739356082,"model":"deepseek-r1:7b","system_fingerprint":"fp_ollama","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":"stop"}]}
			//	data: [DONE]
			print(line)
			if strings.HasPrefix(line, "data: ") { //data 没引号 解析
				line = line[6:]
			}
			if line == "[DONE]\n" {
				break
			}
			var data map[string]interface{}
			err = json.Unmarshal([]byte(line), &data)
			if err != nil {
				outFunc("error: %s\n", err.Error())
				println("error json.Unmarshal: " + err.Error())
				return
			}
			if data["error"] != nil {
				errs := data["error"].(map[string]interface{})["message"].(string)
				outFunc("error: %s\n", errs)
				println("error: " + errs)
				return
			}
			choices := data["choices"].([]interface{})
			var thinking int64
			for _, choice := range choices {
				delta := choice.(map[string]interface{})["delta"].(map[string]interface{})
				msg := delta["content"].(string)
				if role == "" {
					role = delta["role"].(string)
					outFunc(role + ": ")
				}
				switch msg {
				case "<think>":
					thinking = time.Now().Unix()
					msg = "思考中..."
				case "</think>":
					thinking = 0
					during := time.Now().Unix() - thinking
					switch {
					case during < 60:
						msg = fmt.Sprintf("思考结束: (用时 %d 秒)", during)
					case during < 3600:
						msg = fmt.Sprintf("思考结束: (用时 %d 分 %d 秒)", during/60, during%60)
					default:
						msg = fmt.Sprintf("思考结束: (用时 %d 小时 %d 分 %d 秒)", during/3600, during%3600/60, during%60)
					}
				}
				outFunc(msg)
				if thinking > 0 {
					continue
				}
				contentX += msg
			}
			//"finish_reason":null}]}
			//"finish_reason":"stop"}]}
			//finish_reason := data["finish_reason"].(string)
		}
		payload["messages"] = append(payload["messages"].([]map[string]string), map[string]string{"role": role, "content": contentX})
	}()
	select {
	case <-done:
		print("Stream ended or closed")
	}
	return nil
}

func (c *Chat) handleNonStreamResponse(resp *http.Response) error {
	payload := c.Payload
	outFunc := c.outFunc
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	res := string(response)
	print(res)
	// example response:
	//{"error":{"message":"Authentication Fails (no such user)","type":"authentication_error","param":null,"code":"invalid_request_error"}} //认证失败
	//{"error":{"message":"Insufficient Balance","type":"unknown_error","param":null,"code":"invalid_request_error"}} //余额不足
	//{"error":{"message":"model \"deepseek-chat\" not found, try pulling it first","type":"api_error","param":null,"code":null}}
	// {"id":"chatcmpl-386","object":"chat.completion","created":1739351889,"model":"deepseek-r1:7b","system_fingerprint":"fp_ollama",
	//"choices":[{"index":0,"message":{"role":"assistant","content":"\u003cthink\u003e\nOkay, the user says \"Hello!\" and that's it. I should respond politely.\n\nI'll greet them back in a friendly manner.\n\nMaybe ask how they can help today to keep the conversation going.\n\nKeep it simple and open-ended so they feel comfortable sharing more.\n\u003c/think\u003e\n\nHello! How can I assist you today?"},"finish_reason":"stop"}],"usage":{"prompt_tokens":11,"completion_tokens":68,"total_tokens":79}}

	var data map[string]interface{}
	err = json.Unmarshal(response, &data)
	if err != nil {
		outFunc("error: %s\n", err.Error())
		println("error: " + err.Error())
		return err
	}
	if data["error"] != nil {
		errMap := data["error"].(map[string]interface{})
		errMsg := errMap["message"].(string)
		//errTyp := errMap["type"].(string)
		errCode := errMap["code"].(string)
		switch errCode {
		case "invalid_request_error":
		}
		outFunc("error: %s\n", errMsg)
		println("error: " + errMsg)
		return nil
	}
	choices := data["choices"].([]interface{})
	for _, choice := range choices {
		message := choice.(map[string]interface{})["message"].(map[string]interface{})
		role := message["role"].(string)
		content := message["content"].(string)
		outFunc("%s: %s\n", role, content)
		payload["messages"] = append(payload["messages"].([]map[string]string), map[string]string{"role": role, "content": content})
	}
	return nil
}

var chat *Chat

func SendChat(str string) {
	if chat == nil {
		chat = NewChat("deepseek-chat", UIChatLog) //deepseek官方
		chat.SetModel("deepseek-r1:7b")            //Ollama 本地
		//chat.ApplyChat("system", "You are a helpful assistant.")
		chat.ApplyChat("system", "假如你是一个玄学大师,精通八字命理,奇门遁甲,大六壬,梅花易数,星盘解读.")
	}
	chat.ApplyChat("system", fmt.Sprintf("当前时间是: %s", time.Now().Format(time.DateTime)))
	chat.SendChat("user", str)
}
