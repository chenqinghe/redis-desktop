package main

import (
	"fmt"
	"github.com/lxn/walk"
)

type TabWidgetEx struct {
	*walk.TabWidget

	root *MainWindowEX

	pages []*TabPageEx
}

func (tw *TabWidgetEx) startNewSession(sess session) {
	tabPage, err := tw.NewTabPageEx()
	if err != nil {
		walk.MsgBox(nil, "ERROR", "新建标签页失败："+err.Error(), walk.MsgBoxIconError)
		return
	}
	tw.SetCurrentIndex(tw.Pages().Len() - 1)
	tabPage.SetTitle(fmt.Sprintf("%s:%d", sess.Host, sess.Port))

	tabPage.content.SetText("")
	tabPage.content.AppendText(fmt.Sprintf("connecting to %s:%d ......\r\n", sess.Host, sess.Port))
	conn, err := DialRedis(sess.Host, sess.Port, sess.Password)
	if err != nil {
		tabPage.content.AppendText(err.Error())
		return
	}
	tabPage.conn = conn

	if r := execCmd(conn, "PING"); r != "PONG" {
		tabPage.content.AppendText(r + "\r\n")
		return
	}
	tabPage.content.AppendText("连接成功！\r\n\r\n")
	tabPage.content.AppendText("> ")
	tabPage.content.SetReadOnly(false)
}

func (tw *TabWidgetEx) AddPage(page *TabPageEx) {
	tw.Pages().Add(page.TabPage)
	tw.pages = append(tw.pages, page)
}

func (tw *TabWidgetEx) CurrentPage() *TabPageEx {
	idx := tw.CurrentIndex()
	if idx > 0 {
		return tw.pages[idx-1] // TODO: tw.Pages() contains home page, if home page could be remove, idx not -1
	}
	return nil
}
