package main

import (
	"fmt"
	"github.com/lxn/walk"
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
	lb.Model = append(lb.Model, fmt.Sprintf("%s:%d", sess.Host, sess.Port))
	lb.sessions = append(lb.sessions, sess)
}

func (lb *ListBoxEX) AddSessions(sesses []session) {
	for _, sess := range sesses {
		lb.AddSession(sess)
	}
}

func (lb *ListBoxEX) GetSessions() []session {
	return lb.sessions
}


func (lb *ListBoxEX)CurrentSession()session {
	return lb.sessions[lb.CurrentIndex()]
}