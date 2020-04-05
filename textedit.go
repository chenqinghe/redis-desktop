package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"strings"
)

type TextEditEx struct {
	root *MainWindowEX

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

	bg, err := walk.NewSolidColorBrush(walk.RGB(0, 0, 0))
	if err != nil {
		return nil, err
	}
	textedit.SetBackground(bg)
	//font, err := walk.NewFont("微软雅黑", 12, 0)
	//if err != nil {
	//	return nil, err
	//}
	//textedit.SetFont(font)
	textedit.SetTextColor(walk.RGB(255,255,255))

	walk.InitWrapperWindow(textEditEx)

	return textEditEx, nil
}

func (te *TextEditEx) OnKeyPress(key walk.Key) {
	te.SetTextSelection(te.TextLength(), te.TextLength())
	if key == walk.KeyReturn {
		content := te.Text()
		content = strings.TrimSpace(content)
		for i := len(content) - 1; i > 0; i-- {
			if content[i] == '>' {
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

func (te *TextEditEx) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	//if msg == win.WM_CHAR {
	//	//r := te.TextEdit.WndProc(hwnd, msg, wParam, lParam)
	//	fmt.Println("WM_CHAR:", wParam, lParam)
	//	return 1
	//}
	//if msg == win.WM_KEYDOWN {
	//	logrus.Debugln("TextEditEx WndProc:", hwnd, msg, wParam, lParam)
	//	_, file, line, ok := runtime.Caller(1)
	//	fmt.Printf("%s:%d  %t\n", file, line, ok)
	//	//ne:= -1
	//	r := []uintptr{0, 1}[rand.Intn(1)]
	//	fmt.Println("ret:", r)
	//	return r
	//}
	//if msg == win.WM_KEYUP {
	//	return 0
	//}
	ret := te.TextEdit.WndProc(hwnd, msg, wParam, lParam)
	return ret
}
