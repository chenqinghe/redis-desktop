package main

import (
	"fmt"
	"strings"

	"github.com/chenqinghe/walk"
	"github.com/chenqinghe/walk/declarative"
	"github.com/lxn/win"
	"github.com/sirupsen/logrus"
)

type TextEditEx struct {
	root   *MainWindowEX
	parent *TabPageEx

	*walk.TextEdit

	offset int
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
				Text:        "执行选中命令",
				OnTriggered: textEditEx.RunSelectCmd,
			},
			declarative.Action{
				Text: "复制",
				OnTriggered: func() {

				},
			},
			declarative.Separator{},
			declarative.Action{
				Text:        "清屏",
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

	// TODO: when textedit visible change to true, move cursor to the end
	//textEditEx.TextEdit.VisibleChanged().Attach(func() {
	//	fmt.Println("textedit visible changed:", textEditEx.TextEdit.Visible())
	//	if textEditEx.TextEdit.Visible() {
	//		textEditEx.TextEdit.SetTextSelection(0, 1)
	//	}
	//})

	walk.InitWrapperWindow(textEditEx)

	return textEditEx, nil
}

func (te *TextEditEx) OnKeyPress(key walk.Key) {
	te.SetTextSelection(te.TextLength(), te.TextLength())
	if key == walk.KeyReturn {
		content := te.Text()
		cmd := content[len(content)-te.offset:]
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
		walk.MsgBox(nil, "INFO", "选中的命令为空！", walk.MsgBoxIconInformation)
		fmt.Println("nothing to run")
		return
	}
	cmd := string([]rune(te.Text())[start:end])
	logrus.Debugln("selected cmd:", cmd)

	te.clearCmdBuffer()
	te.AppendText(cmd)
	te.moveCursorToEnd()
	te.offset = len(cmd)
	te.runCmd(cmd)
}

func (te *TextEditEx) moveCursorToEnd() {
	te.SetTextSelection(te.TextLength(), te.TextLength())
}

func (te *TextEditEx) clearCmdBuffer() {
	text := te.Text()
	start, end := len(text)-te.offset, len(text)
	bufferedCmd := text[start:end]
	runeLen := len([]rune(bufferedCmd))
	start, end = len([]rune(text))-runeLen, len([]rune(text))
	logrus.Debugln("clearCmdBuffer: start:", start, "end:", end, "to be cleared:", text[start:end])
	te.SetTextSelection(start, end)
	te.ReplaceSelectedText("", false)

	te.offset = 0
}

func (te *TextEditEx) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == win.WM_CHAR {
		logrus.Debugln("WndProc: WM_CHAR:", wParam, lParam)
		if walk.Key(wParam) == walk.KeyBack {
			fmt.Println("backspace pressed!")
			if te.offset <= 0 {
				return 0
			}
			te.offset--
			fmt.Println("offset:", te.offset)
		} else {
			te.offset++
		}
	}
	return te.TextEdit.WndProc(hwnd, msg, wParam, lParam)
}
