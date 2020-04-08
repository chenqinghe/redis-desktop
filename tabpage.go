package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/sirupsen/logrus"
)

type TabPageEx struct {
	*walk.TabPage

	content *TextEditEx

	conn redis.Conn
}

func (tw *TabWidgetEx) NewTabPageEx() (*TabPageEx, error) {
	// FIXME: 移除默认tab，下面的方式可以移除，但是会造成程序无响应
	//if tw.Pages().At(0).Title() == "home" {
	//	tw.Pages().RemoveAt(0)
	//}

	tabpageEx := &TabPageEx{
		content: nil,
		conn:    nil,
	}

	if err := (TabPage{
		AssignTo: &tabpageEx.TabPage,
		Image:    "img/redis.ico",
		Layout: HBox{
			MarginsZero: true,
			SpacingZero: true,
		},
		ContextMenuItems: []MenuItem{
			Action{
				Text: "关闭会话",
				OnTriggered: func() {
					// TODO: remove the page
				},
			},
		},
	}).Create(NewBuilder(nil)); err != nil {
		logrus.Errorln("create tabpage error:", err)
		return nil, err
	}

	textedit, err := NewTextEdit(tw.root, tabpageEx)
	if err != nil {
		return nil, err
	}
	tabpageEx.SetContent(textedit)

	tw.AddPage(tabpageEx)
	return tabpageEx, nil
}

func (p *TabPageEx) SetContent(textedit *TextEditEx) {
	p.content = textedit
}
