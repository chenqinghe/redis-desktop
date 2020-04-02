package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func coredump(msg string) {
	ioutil.WriteFile("coredump", []byte(msg), os.ModePerm)
	os.Exit(1)
}

func loadSession(file string) ([]session, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	sessions := make([]session, 0)
	if err := json.Unmarshal(data, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func saveSession(sessions []session, file string) error {
	data, err := json.Marshal(sessions)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, os.ModePerm)
}

func main() {
	mw := &MainWindowEX{
		PB_connect:     new(PushButtonEx),
		LB_sessions:    new(ListBoxEX),
		TW_screenGroup: new(TabWidgetEx),
	}
	mw.PB_connect.root = mw
	mw.LB_sessions.root = mw
	mw.TW_screenGroup.root = mw

	rootPath := os.Getenv("APPDATA")
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

	mw.sessionFile = filepath.Join(rootPath, "RedisDesktop", "data")
	sessions, err := loadSession(mw.sessionFile)
	if err != nil {
		walk.MsgBox(nil, "ERROR", "加载会话文件失败："+err.Error(), walk.MsgBoxIconError)
		return
	}
	mw.LB_sessions.AddSessions(sessions)

	n, err := MainWindow{
		Title:    "redis命令行工具",
		MinSize:  Size{600, 400},
		AssignTo: &mw.MainWindow,
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			VSplitter{
				Children: []Widget{
					Composite{
						MaxSize: Size{0, 50},
						Layout:  HBox{},
						Children: []Widget{
							Label{Text: "host"},
							LineEdit{AssignTo: &mw.LE_host},
							Label{Text: "port"},
							LineEdit{AssignTo: &mw.LE_port},
							Label{Text: "password"},
							LineEdit{AssignTo: &mw.LE_password},
							PushButton{
								Text:      "连接",
								AssignTo:  &mw.PB_connect.PushButton,
								OnClicked: mw.PB_connect.OnClick,
							},
						},
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							ListBox{
								MaxSize:  Size{200, 0},
								AssignTo: &mw.LB_sessions.ListBox,
								Model:    mw.LB_sessions.Model,
								OnCurrentIndexChanged: func() {
								},
								OnItemActivated: func() {
									curSession := mw.LB_sessions.CurrentSession()
									mw.LE_host.SetText(curSession.Host)
									mw.LE_port.SetText(strconv.Itoa(curSession.Port))
									mw.LE_password.SetText(curSession.Password)
									mw.PB_connect.OnClick()
								},
								OnSelectedIndexesChanged: func() {
								},
								MultiSelection: false,
							},
							TabWidget{
								AssignTo:           &mw.TW_screenGroup.TabWidget,
								ContentMarginsZero: true,
								Pages: []TabPage{
									TabPage{
										Title: "tab1",
										Content: TextEdit{
											VScroll:  true,
											ReadOnly: true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}.Run()

	fmt.Println(n, err)
}
