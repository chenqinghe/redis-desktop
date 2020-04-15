package main

import (
	"github.com/chenqinghe/walk"
	. "github.com/chenqinghe/walk/declarative"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type TabPageEx struct {
	*walk.TabPage

	content *TextEditEx

	conn redis.Conn
}

func (tw *TabWidgetEx) NewTabPageEx() (*TabPageEx, error) {
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
	}).Create(NewBuilder(tw.root)); err != nil {
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
