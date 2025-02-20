package world

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/deminzhang/qimen-go/gui"
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
func UIChatLogLn(msg string, a ...any) {
	if uiChat != nil {
		uiChat.ChatLogLn(msg, a...)
	}
}
func UIChatLog(msg string, a ...any) {
	if uiChat != nil {
		uiChat.ChatLog(msg, a...)
	} else {
		print(msg)
	}
}

func NewUIChat() *UIChat {
	p := &UIChat{
		BaseUI: gui.BaseUI{Visible: true,
			X: 0, Y: ScreenHeight - 250, W: 550, H: 250,
		},
		showChatUI: false,
	}
	textBoxLog := gui.NewTextBox(16, 0, 470, 180)
	inputBoxChat := gui.NewInputBox(16, 190, 470, 24)

	checkBoxGM := gui.NewCheckBox(490, 200, "令")
	btnChatSend := gui.NewButton(490, 230, 32, 16, "发")

	btnChatSwitch := gui.NewButton(0, 230, 32, 16, "隐")
	btnClear := gui.NewButton(32, 230, 32, 16, "清")
	btnCopy := gui.NewButton(64, 230, 32, 16, "复")

	p.AddChildren(btnChatSwitch, btnClear, btnCopy,
		textBoxLog,
		inputBoxChat,
		checkBoxGM,
		btnChatSend)

	inputBoxChat.DefaultText = "输入内容发送.."

	inputBoxChat.SetOnPressEnter(func(i *gui.InputBox) {
		if !p.showChatUI {
			return
		}
		if i.Focused() {
			btnChatSend.Click()
		} else {
			i.SetFocused(true)
		}
		if i.TextField.IsFocused() {
			i.TextField.Blur()
		} else {
			i.TextField.Focus()
		}
	})
	checkBoxGM.SetOnCheckChanged(func(c *gui.CheckBox) {
		var msg string
		if c.Checked() {
			msg = " 命令模式(help查看命令)"
		} else {
			msg = " 聊天模式(AI)"
		}
		textBoxLog.AppendTextLn(msg)
	})
	checkBoxGM.SetChecked(true)
	btnChatSend.SetOnClick(func() {
		i := inputBoxChat
		if i.Text() != "" {
			i.SetFocused(false)
			if checkBoxGM.Checked() { //调试命令
				textBoxLog.AppendTextLn("cmd: " + i.Text())
				parseCmd(i.Text())
			} else { //聊天
				p.ChatLogLn("me: %s\n", i.Text())
				go SendChat(i.Text())
			}
			i.AppendTextHistory(i.Text())
			i.SetText("")
		} else {
			i.SetFocused(false)
			textBoxLog.AppendTextLn("no input msg")
		}
	})
	resetView := func() {
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
	}
	btnChatSwitch.SetOnClick(func() {
		p.showChatUI = !p.showChatUI
		resetView()
	})
	btnClear.SetOnClick(func() {
		textBoxLog.Text = ""
	})
	btnCopy.SetOnClick(func() {
		clipboard.WriteAll(textBoxLog.Text)
	})

	p.textBoxLog = textBoxLog
	uiChat = p
	resetView()
	return p
}

func (p *UIChat) Update() {
	p.BaseUI.Update()
}

func (p *UIChat) OnClose() {
	uiChat = nil
}

func (p *UIChat) ChatLogLn(msg string, a ...any) {
	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a...)
	}
	p.textBoxLog.AppendTextLn(msg)
}
func (p *UIChat) ChatLog(msg string, a ...any) {
	if len(a) > 0 {
		msg = fmt.Sprintf(msg, a...)
	}
	p.textBoxLog.AppendText(msg)
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
		UIChatLogLn("help: show help")
		UIChatLogLn("showBattle: 显战斗")
		UIChatLogLn("hideBattle: 隐战斗")
		UIChatLogLn("showWave: 显示引力波")
		UIChatLogLn("hideWave: 隐引力波")
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
