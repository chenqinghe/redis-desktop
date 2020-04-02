package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/lxn/walk"
)

type TabWidgetEx struct {
	*walk.TabWidget

	root *MainWindowEX

	pages []*TabPageEx
}

func (tw *TabWidgetEx) NewTabPageEx() (*TabPageEx, error) {
	if tw.Pages().At(0).Title() == "tab1" {
		tw.Pages().RemoveAt(0)
	}
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

	tabpageEx.TextEditEx = textedit

	tw.Pages().Add(tabpage)

	tw.pages = append(tw.pages, tabpageEx)
	return tabpageEx, nil
}

type TabPageEx struct {
	*walk.TabPage

	*TextEditEx

	conn redis.Conn
}
