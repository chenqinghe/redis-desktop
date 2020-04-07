package main

import (
	"encoding/json"
	"fmt"
	"github.com/chenqinghe/redis-desktop/i18n"
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

	lang i18n.Lang

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

func (mw *MainWindowEX) SetSessionFile(file string) {
	mw.sessionFile = file
}

func (mw *MainWindowEX) LoadSession() error {
	data, err := ioutil.ReadFile(mw.sessionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	sessions := make([]session, 0)
	if err := json.Unmarshal(data, &sessions); err != nil {
		return err
	}

	mw.LB_sessions.AddSessions(sessions)
	return nil
}

func createMainWindow(lang i18n.Lang) *MainWindowEX {
	mw := &MainWindowEX{
		lang:           lang,
		PB_connect:     new(PushButtonEx),
		LB_sessions:    new(ListBoxEX),
		TW_screenGroup: new(TabWidgetEx),
	}
	mw.PB_connect.root = mw
	mw.LB_sessions.root = mw
	mw.TW_screenGroup.root = mw
	err := MainWindow{
		Title:    mw.lang.Tr("mainwindow.title"),
		MinSize:  Size{600, 400},
		AssignTo: &mw.MainWindow,
		Layout:   VBox{MarginsZero: true},
		MenuItems: []MenuItem{
			Menu{
				Text: mw.lang.Tr("mainwindow.menu.file"),
				Items: []MenuItem{
					Action{
						Text:        mw.lang.Tr("mainwindow.menu.file.export"),
						OnTriggered: nil,
					},
					Action{
						Text:        mw.lang.Tr("mainwindow.menu.file.import"),
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: mw.lang.Tr("mainwindow.menu.edit"),
				Items: []MenuItem{
					Action{
						Text:        mw.lang.Tr("mainwindow.menu.edit.clear"),
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: mw.lang.Tr("mainwindow.menu.setting"),
				Items: []MenuItem{
					Action{
						Text:        mw.lang.Tr("mainwindow.menu.setting.theme"),
						OnTriggered: nil,
					},
					Action{
						Text:        mw.lang.Tr("mainwindow.menu.logpath"),
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: mw.lang.Tr("mainwindow.menu.run"),
				Items: []MenuItem{
					Action{
						Text:        mw.lang.Tr("mainwindow.menu.run.batch"),
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: mw.lang.Tr("mainwindow.menu.help"),
				Items: []MenuItem{
					Action{
						Text: mw.lang.Tr("mainwindow.menu.help.source"),
						OnTriggered: func() {
							startPage("https://github.com/chenqinghe/redis-desktop")
						},
					},
					Action{
						Text:        mw.lang.Tr("mainwindow.menu.help.bug"),
						OnTriggered: startIssuePage,
					},
					Separator{},
					Action{
						Text: mw.lang.Tr("mainwindow.menu.help.donate"),
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
							Label{Text: mw.lang.Tr("mainwindow.labelhost")},
							LineEdit{AssignTo: &mw.LE_host},
							Label{Text: mw.lang.Tr("mainwindow.labelport")},
							LineEdit{AssignTo: &mw.LE_port},
							Label{Text: mw.lang.Tr("mainwindow.labelpassword")},
							LineEdit{AssignTo: &mw.LE_password, PasswordMode: true},
							PushButton{
								Text:      mw.lang.Tr("mainwindow.PBconnect"),
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
										Text:        mw.lang.Tr("mainwindow.LBsessions.menu.deletesession"),
										OnTriggered: mw.LB_sessions.RemoveSelectedSession,
									},
								},
							},
							TabWidget{
								AssignTo: &mw.TW_screenGroup.TabWidget,
								Pages: []TabPage{
									TabPage{
										Title: "tab1",
										Content: ImageView{
											Mode:  ImageViewModeStretch,
											Image: "img/cover.png",
										},
									},
								},
								ContentMarginsZero: true,
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
	if _, err := (Dialog{
		Title: "捐赠",
		Layout: VBox{
			MarginsZero: true,
		},
		Children: []Widget{
			ImageView{
				Image:   "img/cover.jpg",
				Mode:    ImageViewModeStretch,
				MinSize: Size{500, 500},
				MaxSize: Size{500, 500},
			},
		},
	}).Run(p); err != nil {
		logrus.Errorln("showDonate: run dialog error:", err)
	}
}
