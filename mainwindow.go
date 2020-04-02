package main

import (
	"github.com/lxn/walk"
)

type MainWindowEX struct {
	*walk.MainWindow

	logFile string

	LE_host     *walk.LineEdit
	LE_port     *walk.LineEdit
	LE_password *walk.LineEdit

	PB_connect *PushButtonEx

	sessionFile string
	LB_sessions *ListBoxEX

	TW_screenGroup *TabWidgetEx
}
