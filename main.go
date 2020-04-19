package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chenqinghe/redis-desktop/config"
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
	appPath := filepath.Join(rootPath, "RedisDesktop")

	// redirect stdout & stderr to ensure capture panic info
	//logFile := filepath.Join(appPath, "RedisDesktop.log")
	//if err := os.MkdirAll(filepath.Dir(logFile), os.ModePerm); err != nil {
	//	coredump(err.Error())
	//}
	//f, err := NewFileDescriptor(logFile)
	//if err != nil {
	//	coredump(err.Error())
	//}
	//if err := SetStdHandle(f, f); err != nil {
	//	coredump(err.Error())
	//}
	var err error

	configPath := filepath.Join(appPath, "config.toml")
	_, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := config.Save(configPath); err != nil {
				logrus.Fatalln("save config:", err)
			}
		} else {
			logrus.Fatalln("stat config file:", err)
		}
	} else {
		if err := config.Load(configPath); err != nil {
			logrus.Fatalln("load config:", err)
		}
	}

	lv, err := logrus.ParseLevel(config.Get().LogConfig.Level)
	if err != nil {
		lv = logrus.InfoLevel
	}
	logrus.SetLevel(lv)

	lang, ok := i18n.GetLang(config.Get().Lang)
	if ok {
		i18n.SetDefaultLang(lang)
	}
	mw := createMainWindow()

	mw.SetSessionFile(filepath.Join(appPath, "sessions.dat"))
	if err := mw.LoadSession(); err != nil {
		walk.MsgBox(mw, "ERROR", i18n.Tr("alert.loadsessionfailed", err.Error()), walk.MsgBoxIconError)
		return
	}

	mw.Run()
}
