package main

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/lxn/walk"
)

func coredump(msg string) {
	ioutil.WriteFile("coredump", []byte(msg), os.ModePerm)
	os.Exit(1)
}

func main() {
	rootPath := os.Getenv("APPDATA")

	logrus.SetLevel(logrus.DebugLevel)

	// redirect stdout & stderr to ensure capture panic info
	//mw.logFile = filepath.Join(rootPath, "RedisDesktop", "log")
	//if err := os.MkdirAll(filepath.Dir(mw.logFile), os.ModePerm); err != nil {
	//	coredump(err.Error())
	//}
	//f, err := NewFileDescriptor(mw.logFile)
	//if err != nil {
	//	coredump(err.Error())
	//}
	//if err := SetStdHandle(f, f); err != nil {
	//	coredump(err.Error())
	//}

	mw := createMainWindow()
	mw.sessionFile = filepath.Join(rootPath, "RedisDesktop", "data")
	sessions, err := loadSession(mw.sessionFile)
	if err != nil {
		walk.MsgBox(nil, "ERROR", "加载会话文件失败："+err.Error(), walk.MsgBoxIconError)
		return
	}
	mw.LB_sessions.AddSessions(sessions)

	mw.Run()
}
