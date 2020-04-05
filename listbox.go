package main

import (
	"fmt"
	"github.com/lxn/walk"
	"github.com/sirupsen/logrus"
)

type session struct {
	Host     string
	Port     int
	Password string
}

type ListBoxEX struct {
	*walk.ListBox

	root *MainWindowEX

	sessions []session
	Model    []string
}

func (lb *ListBoxEX) AddSession(sess session) {
	lb.addSession(sess)
	lb.ReloadModel()
	if err:=lb.root.saveSessions(lb.sessions);err!=nil {
		logrus.Errorln("save sessions error:",err)
	}
}
func (lb *ListBoxEX) addSession(sess session) {
	lb.Model = append(lb.Model, fmt.Sprintf("%s:%d", sess.Host, sess.Port))
	lb.sessions = append(lb.sessions, sess)
}

func (lb *ListBoxEX) AddSessions(sesses []session) {
	for _, sess := range sesses {
		lb.addSession(sess)
	}
	lb.ReloadModel()
	lb.root.saveSessions(lb.sessions)
}

func (lb *ListBoxEX) GetSessions() []session {
	return lb.sessions
}

func (lb *ListBoxEX) CurrentSession() session {
	return lb.sessions[lb.CurrentIndex()]
}

func (lb *ListBoxEX) SessionCount() int {
	return len(lb.sessions)
}

func (lb *ListBoxEX) RemoveSession(sess session) {
	var index int = -1
	for k, v := range lb.sessions {
		if v.Port == sess.Port && v.Host == sess.Host {
			index = k
		}
	}
	if index > 0 {
		lb.sessions[index] = lb.sessions[len(lb.sessions)-1]
		lb.sessions = lb.sessions[:len(lb.sessions)-1]
		lb.Model[index] = lb.Model[len(lb.Model)-1]
		lb.Model = lb.Model[:len(lb.Model)-1]
	}
	lb.ReloadModel()
	lb.root.saveSessions(lb.sessions)
}

func (lb *ListBoxEX) RemoveSelectedSession() {
	index := lb.CurrentIndex()
	if index > 0 {
		lb.sessions[index] = lb.sessions[len(lb.sessions)-1]
		lb.sessions = lb.sessions[:len(lb.sessions)-1]
		lb.Model[index] = lb.Model[len(lb.Model)-1]
		lb.Model = lb.Model[:len(lb.Model)-1]
		lb.ReloadModel()
		lb.root.saveSessions(lb.sessions)
	}
}

func (lb *ListBoxEX) ReloadModel() {
	lb.SetModel(lb.Model)
}
