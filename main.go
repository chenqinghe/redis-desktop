package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chenqinghe/redis-desktop/i18n"
	"github.com/chenqinghe/walk"
	"github.com/sirupsen/logrus"
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

	// TODO: how to config languages?
	lang, ok := i18n.GetLang("zh_cn")
	if ok {
		i18n.SetDefaultLang(lang)
	}
	mw := createMainWindow()

	mw.SetSessionFile(filepath.Join(rootPath, "RedisDesktop", "sessions.dat"))
	if err := mw.LoadSession(); err != nil {
		walk.MsgBox(mw, "ERROR", "加载会话文件失败："+err.Error(), walk.MsgBoxIconError)
		return
	}

	mw.Run()
}
