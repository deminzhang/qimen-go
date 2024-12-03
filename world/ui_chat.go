package world

import (
	"fmt"
	"github.com/deminzhang/qimen-go/gui"
)

type UIChat struct {
	gui.BaseUI
	btnChatSwitch *gui.Button
	btnClear      *gui.Button
	textBoxLog    *gui.TextBox
	inputBoxChat  *gui.InputBox
	checkBoxGM    *gui.CheckBox
	btnChatSend   *gui.Button
	showChatUI    bool
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
	p.textBoxLog = gui.NewTextBox(16, 0, 270, 180)
	p.inputBoxChat = gui.NewInputBox(16, 190, 270, 32)
	p.checkBoxGM = gui.NewCheckBox(290, 200, "令")
	p.btnChatSend = gui.NewButton(290, 230, 32, 16, "发")

	p.btnChatSwitch = gui.NewButton(0, 230, 32, 16, "隐")
	p.btnClear = gui.NewButton(32, 230, 32, 16, "清")

	p.AddChild(p.btnChatSwitch)
	p.AddChild(p.btnClear)
	p.AddChild(p.textBoxLog)
	p.AddChild(p.inputBoxChat)
	p.AddChild(p.checkBoxGM)
	p.AddChild(p.btnChatSend)

	p.inputBoxChat.DefaultText = "输入信息指令.."

	p.inputBoxChat.SetOnPressEnter(func(i *gui.InputBox) {
		if !p.showChatUI {
			return
		}
		if i.Focused() {
			p.btnChatSend.Click()
		} else {
			i.SetFocused(true)
		}
	})
	p.checkBoxGM.SetOnCheckChanged(func(c *gui.CheckBox) {
		msg := "debug command"
		if c.Checked() {
			msg += " (On)"
		} else {
			msg += " (Off)"
		}
		p.textBoxLog.AppendLine(msg)
	})
	p.checkBoxGM.SetChecked(true)
	p.btnChatSend.SetOnClick(func(b *gui.Button) {
		i := p.inputBoxChat
		if i.Text() != "" {
			i.SetFocused(false)
			if p.checkBoxGM.Checked() { //调试命令
				p.textBoxLog.AppendLine("gm " + i.Text())
				fmt.Println("GM:" + i.Text())
				//sendGMCmd(i.Text())
			} else { //聊天
				p.textBoxLog.AppendLine("me:" + i.Text())
				//sendChat(World.self, i.Text)
			}
			i.AppendTextHistory(i.Text())
			i.SetText("")
		} else {
			i.SetFocused(false)
			p.textBoxLog.AppendLine("no input msg")
		}
	})
	p.btnChatSwitch.SetOnClick(func(b *gui.Button) {
		p.showChatUI = !p.showChatUI
		if p.showChatUI {
			b.Text = "隐"
		} else {
			b.Text = "聊"
		}
	})
	p.btnClear.SetOnClick(func(b *gui.Button) {
		p.textBoxLog.Text = ""
	})
	uiChat = p
	return p
}

func (p *UIChat) Update() {
	p.BaseUI.Update()
	p.textBoxLog.Visible = p.showChatUI
	p.inputBoxChat.Visible = p.showChatUI
	p.btnChatSend.Visible = p.showChatUI
	p.checkBoxGM.Visible = p.showChatUI
	p.btnClear.Visible = p.showChatUI
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
