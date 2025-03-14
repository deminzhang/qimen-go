package world

import (
	"fmt"
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
		showChatUI: false,
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
		//checkBoxGM,
		btnChatSend)

	inputBoxChat.DefaultText = "输入信息指令.."

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
				//sendGMCmd(i.Text())
			} else { //聊天
				textBoxLog.AppendLine("me:" + i.Text())
				//sendChat(World.self, i.Text)
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
