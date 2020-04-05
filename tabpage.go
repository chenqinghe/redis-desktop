package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/lxn/walk"
	"github.com/sirupsen/logrus"
)

type TabPageEx struct {
	*walk.TabPage

	content *TextEditEx

	conn redis.Conn
}

func (tw *TabWidgetEx) NewTabPageEx() (*TabPageEx, error) {
	// FIXME: 移除默认tab，下面的方式可以移除，但是会造成程序无响应
	//if tw.Pages().At(0).Title() == "tab1" {
	//	tw.Pages().RemoveAt(0)
	//}
	tabpage, err := walk.NewTabPage()
	if err != nil {
		return nil, err
	}
	layout := walk.NewHBoxLayout()
	layout.SetMargins(walk.Margins{})
	layout.SetSpacing(0)
	tabpage.SetLayout(layout)

	tabpageEx := &TabPageEx{
		TabPage: tabpage,
	}

	textedit, err := NewTextEdit(tw.root, tabpageEx)
	if err != nil {
		return nil, err
	}
	tabpageEx.content = textedit

	tabpageEx.content.VisibleChanged().Attach(func() {
		// FIXME: 将光标移到最后而不是选中全部
		//logrus.Debugln("tabpage visible changed!", tabpageEx.TabPage.Visible())
		if tabpageEx.TabPage.Visible() {
			te := tabpageEx.content.TextEdit
			start, end := te.TextLength(), te.TextLength()
			logrus.Debugln("start:", start, "end:", end)

			te.SetFocus()
			te.SetTextSelection(start-1, end)
		}
	})

	tw.Pages().Add(tabpage)

	tw.pages = append(tw.pages, tabpageEx)
	return tabpageEx, nil
}
