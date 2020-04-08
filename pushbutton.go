package main

import (
	"fmt"
	"strconv"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

type PushButtonEx struct {
	*walk.PushButton

	root *MainWindowEX
}

func (pb *PushButtonEx) OnClick() {
	h := pb.root.LE_host.Text()
	if len(h) == 0 {
		walk.MsgBox(nil, "ERROR", "host不能为空", walk.MsgBoxIconError)
		return
	}
	p, err := strconv.Atoi(pb.root.LE_port.Text())
	if err != nil {
		walk.MsgBox(nil, "ERROR", "端口格式错误:"+err.Error(), walk.MsgBoxIconError)
		return
	}
	s := session{Password: pb.root.LE_password.Text(), Host: pb.root.LE_host.Text(), Port: p}

	var exist bool
	for _, v := range pb.root.LB_sessions.GetSessions() {
		if fmt.Sprintf("%s:%d", v.Host, v.Port) == fmt.Sprintf("%s:%d", s.Host, s.Port) {
			exist = true
			break
		}
	}
	if !exist {
		ret := walk.MsgBox(pb.root, "INFO", "是否保存当前会话？", walk.MsgBoxIconQuestion|walk.MsgBoxYesNo)
		if ret == win.IDYES { // save session
			pb.root.LB_sessions.AddSession(s)
		}
	}

	pb.root.TW_screenGroup.startNewSession(s)
}
