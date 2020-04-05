package main

import (
	"encoding/json"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

type MainWindowEX struct {
	*walk.MainWindow

	logFile string

	LE_host     *walk.LineEdit
	LE_port     *walk.LineEdit
	LE_password *walk.LineEdit

	LE_command *walk.LineEdit

	PB_connect *PushButtonEx

	sessionFile string
	LB_sessions *ListBoxEX

	TW_screenGroup *TabWidgetEx
}

func (mw *MainWindowEX) saveSessions(sessions []session) error {
	data, err := json.Marshal(sessions)
	if err != nil {
		return err
	}
RETRY:
	if err := ioutil.WriteFile(mw.sessionFile, data, os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(mw.sessionFile), os.ModePerm); err != nil {
				return err
			}
			goto RETRY
		}
		return err
	}
	return nil
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

func createMainWindow() *MainWindowEX {
	mw := &MainWindowEX{
		PB_connect:     new(PushButtonEx),
		LB_sessions:    new(ListBoxEX),
		TW_screenGroup: new(TabWidgetEx),
	}
	mw.PB_connect.root = mw
	mw.LB_sessions.root = mw
	mw.TW_screenGroup.root = mw
	err := MainWindow{
		Title:    "redis命令行工具",
		MinSize:  Size{600, 400},
		AssignTo: &mw.MainWindow,
		Layout:   VBox{MarginsZero: true},
		MenuItems: []MenuItem{
			Menu{
				Text: "文件",
				Items: []MenuItem{
					Action{
						Text:        "导出会话...",
						OnTriggered: nil,
					},
					Action{
						Text:        "导入会话...",
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: "编辑",
				Items: []MenuItem{
					Action{
						Text:        "清屏",
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: "设置",
				Items: []MenuItem{
					Action{
						Text:        "主题",
						OnTriggered: nil,
					},
					Action{
						Text:        "日志路径",
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: "运行",
				Items: []MenuItem{
					Action{
						Text:        "批量运行命令",
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: "帮助",
				Items: []MenuItem{
					Action{
						Text: "查看源码",
						OnTriggered: func() {
							startPage("https://github.com/chenqinghe/redis-desktop")
						},
					},
					Action{
						Text:        "报bug",
						OnTriggered: startIssuePage,
					},
					Separator{},
					Action{
						Text: "捐赠",
						OnTriggered: func() {
							showDonate(mw)
						},
					},
				},
			},
		},
		Children: []Widget{
			LineEdit{
				AssignTo: &mw.LE_command,
				Visible:  false,
			},
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
							LineEdit{AssignTo: &mw.LE_password, PasswordMode: true},
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
								OnItemActivated: func() {
									if mw.LB_sessions.CurrentIndex() >= 0 {
										mw.TW_screenGroup.startNewSession(mw.LB_sessions.CurrentSession())
									}
								},
								OnSelectedIndexesChanged: func() { mw.LB_sessions.EnsureItemVisible(0) },
								OnCurrentIndexChanged:    func() {},
								MultiSelection:           false,
								ContextMenuItems: []MenuItem{
									Action{
										Text:        "删除会话",
										OnTriggered: mw.LB_sessions.RemoveSelectedSession,
									},
								},
							},
							TabWidget{
								AssignTo:           &mw.TW_screenGroup.TabWidget,
								ContentMarginsZero: true,
								Pages: []TabPage{
									TabPage{
										Title: "tab1",
										Content: ImageView{
											Mode:  ImageViewModeStretch,
											Image: "img/cover.png",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}.Create()
	if err != nil {
		log.Fatalln(err)
	}

	icon, _ := walk.NewIconFromFile("img/redis.ico")
	mw.SetIcon(icon)

	return mw
}

func startIssuePage() {
	body := url.QueryEscape(fmt.Sprintf(issueTemplate, VERSION))
	uri := fmt.Sprintf("https://github.com/chenqinghe/redis-desktop/issues/new?body=%s", body)
	startPage(uri)
}

func startPage(uri string) {
	cmd := exec.Command("cmd", "/C", "start", uri)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorln("exec cmd error:", err)
	}
}

func showDonate(p walk.Form) {
	dialog, err := walk.NewDialog(p)
	if err != nil {
		panic(err)
	}
	layout:= walk.NewVBoxLayout()
	layout.SetMargins(walk.Margins{})
	dialog.SetLayout(layout)
	dialog.SetTitle("捐赠")

	iv, err := walk.NewImageView(dialog)
	if err != nil {
		panic(err)
	}

	f, err := walk.NewImageFromFile("img/cover.jpg")
	if err != nil {
		panic(err)
	}
	iv.SetImage(f)
	iv.SetMode(walk.ImageViewModeStretch)
	if err := iv.SetSize(walk.Size{500, 500}); err != nil {
		logrus.Errorln("set size error:", err)
	}
	if err := dialog.SetSize(walk.Size{500, 500}); err != nil {
		logrus.Errorln("set size error:", err)
	}
	if err:=dialog.SetMinMaxSize(walk.Size{500,500},walk.Size{500,500});err!=nil {
		logrus.Errorln("set size error:",err)
	}

	dialog.Show()

}

var menuItems = []MenuItem{
	Menu{
		Text: "文件",
		Items: []MenuItem{
			Action{
				Text:        "导出会话...",
				OnTriggered: nil,
			},
			Action{
				Text:        "导入会话...",
				OnTriggered: nil,
			},
		},
	},
	Menu{
		Text: "编辑",
		Items: []MenuItem{
			Action{
				Text:        "清屏",
				OnTriggered: nil,
			},
		},
	},
	Menu{
		Text: "设置",
		Items: []MenuItem{
			Action{
				Text:        "主题",
				OnTriggered: nil,
			},
			Action{
				Text:        "日志路径",
				OnTriggered: nil,
			},
		},
	},
	Menu{
		Text: "运行",
		Items: []MenuItem{
			Action{
				Text:        "批量运行命令",
				OnTriggered: nil,
			},
		},
	},
	Menu{
		Text: "帮助",
		Items: []MenuItem{
			Action{
				Text: "查看源码",
				OnTriggered: func() {
					startPage("https://github.com/chenqinghe/redis-desktop")
				},
			},
			Action{
				Text:        "报bug",
				OnTriggered: startIssuePage,
			},
			Separator{},
			Action{
				Text: "捐赠",
				OnTriggered: func() {
				},
			},
		},
	},
}
