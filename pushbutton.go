package main

import (
	"fmt"
	"github.com/lxn/walk"
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
		walk.MsgBox(nil, "ERROR", "host不能为空", walk.MsgBoxIconError)
		return
	}
	p, err := strconv.Atoi(pb.root.LE_port.Text())
	if err != nil {
		walk.MsgBox(nil, "ERROR", "端口格式错误:"+err.Error(), walk.MsgBoxIconError)
		return
	}
	var exist bool
	for _, v := range pb.root.LB_sessions.sessions {
		if fmt.Sprintf("%s:%d", v.Host, v.Port) == fmt.Sprintf("%s:%d", h, p) {
			exist = true
			break
		}
	}
	if !exist {
		ret := walk.MsgBox(nil, "INFO", "是否保存当前会话？", walk.MsgBoxIconQuestion|walk.MsgBoxYesNo)
		if ret == win.IDYES { // save session
			s := session{Password: pb.root.LE_password.Text(), Host: pb.root.LE_host.Text(), Port: p}
			pb.root.LB_sessions.AddSession(s)
			//sessions = append(sessions, s)
			if err := saveSession(pb.root.LB_sessions.GetSessions(), pb.root.sessionFile); err != nil {
				walk.MsgBox(nil, "ERROR", "保存session失败："+err.Error(), walk.MsgBoxIconError)
			}
			pb.root.LB_sessions.SetModel(pb.root.LB_sessions.Model)
		}
	}

	tabPage, err := pb.root.TW_screenGroup.NewTabPageEx()
	if err != nil {
		walk.MsgBox(nil, "ERROR", "新建标签页失败："+err.Error(), walk.MsgBoxIconError)
		return
	}
	pb.root.TW_screenGroup.SetCurrentIndex(pb.root.TW_screenGroup.Pages().Len() - 1)
	tabPage.SetTitle(fmt.Sprintf("%s:%d", h, p))

	tabPage.TextEditEx.SetText("")
	tabPage.TextEditEx.AppendText(fmt.Sprintf("connecting to %s:%d ......\r\n", h, p))
	conn, err := connectToRedis(h, p, pb.root.LE_password.Text())
	if err != nil {
		tabPage.TextEditEx.AppendText(err.Error())
		return
	}
	tabPage.conn = conn

	if r := execCmd(conn, "PING"); r != "PONG" {
		tabPage.TextEditEx.AppendText(r + "\r\n")
		return
	} else {
		fmt.Println("ping response:", r)
	}
	tabPage.TextEditEx.AppendText("连接成功！\r\n\r\n")
	tabPage.TextEditEx.AppendText("> ")
	tabPage.TextEditEx.SetReadOnly(false)
}
