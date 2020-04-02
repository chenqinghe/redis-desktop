package main

import (
	"fmt"
	"github.com/lxn/win"
	"strings"

	"github.com/lxn/walk"
)

type TextEditEx struct {
	root *MainWindowEX

	parent *TabPageEx

	*walk.TextEdit
}

func NewTextEdit(root *MainWindowEX, p *TabPageEx) (*TextEditEx, error) {
	textEditEx := &TextEditEx{
		root:     root,
		parent:   p,
		TextEdit: nil,
	}

	var style uint32
	style |= win.WS_VSCROLL
	textedit, err := walk.NewTextEditWithStyle(p.TabPage, style)
	if err != nil {
		return nil, err
	}
	textEditEx.TextEdit = textedit
	textedit.SetVisible(true)
	textedit.SetReadOnly(true)
	textedit.KeyPress().Attach(textEditEx.OnKeyPress)
	textedit.KeyUp().Attach(textEditEx.OnKeyUp)
	return textEditEx, nil
}

func (te *TextEditEx) OnKeyPress(key walk.Key) {
	te.SetTextSelection(te.TextLength(), te.TextLength())
	if key == walk.KeyReturn {
		content := te.Text()
		content = strings.TrimSpace(content)
		for i := len(content) - 1; i > 0; i-- {
			if content[i] == '>' {
				fmt.Println("command:", content[i+1:])
				resp := execCmd(te.parent.conn, content[i+1:])
				te.AppendText("\r\n")
				te.AppendText(resp)
				te.AppendText("\r\n\r\n> ")
				te.SetTextSelection(te.TextLength(), te.TextLength())
				break
			}
		}
		return
	}
}

func (te *TextEditEx) OnKeyUp(key walk.Key) {
	if key == walk.KeyReturn {
		l := te.TextLength()
		te.SetTextSelection(l-1, l-1)
		te.ReplaceSelectedText(" ", false)
	}
}
