package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/chenqinghe/redis-desktop/i18n"
	"github.com/chenqinghe/walk"
	"github.com/chenqinghe/walk/declarative"
	"github.com/lxn/win"
	"github.com/sirupsen/logrus"
)

type TextEditEx struct {
	*walk.TextEdit

	root   *MainWindowEX
	parent *TabPageEx

	offset int // input rune count, NOT byte

	historyCmds []string
	cmdCursor   int
}

func NewTextEdit(root *MainWindowEX, p *TabPageEx) (*TextEditEx, error) {
	textEditEx := &TextEditEx{
		root:     root,
		parent:   p,
		TextEdit: new(walk.TextEdit),
	}

	builder := declarative.NewBuilder(p)
	if err := (declarative.TextEdit{
		AssignTo: &textEditEx.TextEdit,
		ContextMenuItems: []declarative.MenuItem{
			declarative.Action{
				Text:        i18n.Tr("widget.textedit.menu.execselected"),
				OnTriggered: textEditEx.RunSelectCmd,
			},
			declarative.Action{
				Text: i18n.Tr("widget.textedit.menu.copy"),
				OnTriggered: func() {

				},
			},
			declarative.Separator{},
			declarative.Action{
				Text:        i18n.Tr("widget.textedit.menu.clear"),
				OnTriggered: textEditEx.ClearScreen,
			},
		},
		OnKeyPress: textEditEx.OnKeyPress,
		OnKeyUp:    textEditEx.OnKeyUp,
		ReadOnly:   true,
		TextColor:  walk.RGB(255, 255, 255),
		VScroll:    true,
	}).Create(builder); err != nil {
		logrus.Errorln("create textedit error:", err)
		return nil, err
	}

	textEditEx.FocusedChanged().Attach(func() {
		time.AfterFunc(time.Millisecond*5, func() {
			textEditEx.AsWindowBase().Synchronize(func() {
				if textEditEx.Visible() {
					textEditEx.SetTextSelection(textEditEx.TextLength(), textEditEx.TextLength())
				}
			})
		})
	})

	if err := walk.InitWrapperWindow(textEditEx); err != nil {
		return nil, err
	}

	return textEditEx, nil
}

func (te *TextEditEx) OnKeyPress(key walk.Key) {
	te.SetTextSelection(te.TextLength(), te.TextLength())
	if key == walk.KeyReturn {
		content := []rune(te.Text())
		cmd := string(content[len(content)-te.offset:])
		te.runCmd(cmd)
		return
	}
}

func (te *TextEditEx) ClearScreen() {
	te.TextEdit.SetText("> ")
	te.TextEdit.SetTextSelection(2, 2)
	te.offset = 0
}

func (te *TextEditEx) runCmd(cmd string) {
	cmd = strings.TrimSpace(cmd)

	cmds := make([]string, 0, len(te.historyCmds))
	for _, v := range te.historyCmds {
		if v != cmd {
			cmds = append(cmds, v)
		}
	}
	te.historyCmds = cmds
	te.historyCmds = append(te.historyCmds, cmd)

	resp := execCmd(te.parent.conn, cmd)
	te.AppendText("\r\n")
	te.AppendText(resp)
	te.AppendText("\r\n\r\n> ")
	te.SetTextSelection(te.TextLength(), te.TextLength())
	te.offset = 0
}

func (te *TextEditEx) OnKeyUp(key walk.Key) {
	if key == walk.KeyReturn {
		l := te.TextLength()
		te.SetTextSelection(l-1, l-1)
		te.ReplaceSelectedText(" ", false)
	}
}

func (te *TextEditEx) RunSelectCmd() {
	start, end := te.TextSelection()
	logrus.Debugln("RunSelectCmd: start:", start, "end:", end)
	if end-start <= 0 {
		walk.MsgBox(nil, "INFO", i18n.Tr("alert.selectedcmdempty"), walk.MsgBoxIconInformation)
		fmt.Println("nothing to run")
		return
	}
	cmd := string([]rune(te.Text())[start:end])
	logrus.Debugln("selected cmd:", cmd)

	te.clearCmdBuffer()
	te.AppendText(cmd)
	te.moveCursorToEnd()
	te.offset = end - start
	te.runCmd(cmd)
}

func (te *TextEditEx) moveCursorToEnd() {
	length := te.TextLength()
	te.SetTextSelection(length, length)
}

func (te *TextEditEx) clearCmdBuffer() {
	text := te.Text()
	start, end := len([]rune(text))-te.offset, len([]rune(text))
	logrus.Debugln("clearCmdBuffer: start:", start, "end:", end, "to be cleared:", text[start:end])
	te.SetTextSelection(start, end)
	te.ReplaceSelectedText("", false)

	te.offset = 0
}

func (te *TextEditEx) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == win.WM_CHAR {
		logrus.Debugln("WndProc: WM_CHAR:", wParam, lParam)

		switch key := walk.Key(wParam); key {
		case walk.KeyBack:
			if te.offset <= 0 {
				return 0
			}
			te.offset--
		case walk.KeyUp, walk.KeyDown:
			if key == walk.KeyUp {
				te.cmdCursor++
			} else if key == walk.KeyDown {
				if te.cmdCursor > 0 {
					te.cmdCursor--
				}
			}
			cmd := te.historyCmds[len(te.historyCmds)-te.cmdCursor]
			te.SetTextSelection(te.TextLength()-te.offset, te.TextLength())
			te.ReplaceSelectedText(cmd, false)
			return 0
		default:
			te.cmdCursor = 0
			te.offset++
		}
	} else if msg == win.WM_KEYDOWN {
		switch key := walk.Key(wParam); key {
		case walk.KeyUp, walk.KeyDown:
			if key == walk.KeyUp {
				if te.cmdCursor < len(te.historyCmds) {
					te.cmdCursor++
				}
			} else if key == walk.KeyDown {
				if te.cmdCursor > 1 {
					te.cmdCursor--
				}
			}
			cmd := te.historyCmds[len(te.historyCmds)-te.cmdCursor]
			te.SetTextSelection(te.TextLength()-te.offset, te.TextLength())
			te.ReplaceSelectedText(cmd, false)
			te.offset = len([]rune(cmd))
			return 0
		}
	}
	return te.TextEdit.WndProc(hwnd, msg, wParam, lParam)
}
