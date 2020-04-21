package main

import (
	"github.com/chenqinghe/redis-desktop/i18n"
	"github.com/chenqinghe/walk"
	"github.com/lxn/win"
	"strconv"
)

type PushButtonEx struct {
	*walk.PushButton

	root *MainWindowEX
}

func (pb *PushButtonEx) OnClick() {
	h := pb.root.LE_host.Text()
	if len(h) == 0 {
		walk.MsgBox(nil, "ERROR", i18n.Tr("alert.hostcannotempty"), walk.MsgBoxIconError)
		return
	}
	p, err := strconv.Atoi(pb.root.LE_port.Text())
	if err != nil {
		walk.MsgBox(nil, "ERROR", i18n.Tr("alert.invalidport", err.Error()), walk.MsgBoxIconError)
		return
	}
	s := Session{Password: pb.root.LE_password.Text(), Host: pb.root.LE_host.Text(), Port: p}

	if !pb.root.TV_sessions.SessionExist(s) {
		ret := walk.MsgBox(pb.root, "INFO", "是否保存当前会话？", walk.MsgBoxIconQuestion|walk.MsgBoxYesNo)
		if ret == win.IDYES { // save session
			pb.root.TV_sessions.AddSession(s)
		}
	}

	pb.root.TW_pages.startNewSession(s)
}
