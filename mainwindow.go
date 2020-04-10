package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/chenqinghe/redis-desktop/i18n"
	"github.com/chenqinghe/walk"
	. "github.com/chenqinghe/walk/declarative"
	"github.com/sirupsen/logrus"
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
	TV_sessions *TreeViewEx
	//LB_sessions *ListBoxEX

	TW_pages *TabWidgetEx
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

	return mw.TV_sessions.LoadSession(data)
}

func (mw *MainWindowEX) importSession(file string) error {
	return mw.TV_sessions.ImportSessions(file)
}

func createMainWindow() *MainWindowEX {
	mw := &MainWindowEX{
		PB_connect:  new(PushButtonEx),
		TV_sessions: &TreeViewEx{model: NewSessionTreeModel()},
		TW_pages:    new(TabWidgetEx),
	}
	mw.PB_connect.root = mw
	mw.TV_sessions.root = mw
	mw.TW_pages.root = mw
	err := MainWindow{
		Title:    i18n.Tr("mainwindow.title"),
		MinSize:  Size{600, 400},
		AssignTo: &mw.MainWindow,
		Layout:   VBox{MarginsZero: true},
		//Background: SolidColorBrush{Color: walk.RGB(132, 34, 234)},
		MenuItems: []MenuItem{
			Menu{
				Text: i18n.Tr("mainwindow.menu.file"),
				Items: []MenuItem{
					Action{
						Text: i18n.Tr("mainwindow.menu.file.import"),
						OnTriggered: func() {
							dlg := &walk.FileDialog{
								Title: "choose a file", // string
							}
							accepted, err := dlg.ShowOpen(mw)
							if err != nil {
								walk.MsgBox(mw, "ERROR", "Open FileDialog:"+err.Error(), walk.MsgBoxIconError)
								return
							}
							if accepted {
								if err := mw.importSession(dlg.FilePath); err != nil {
									walk.MsgBox(mw, "ERROR", "Import Session:"+err.Error(), walk.MsgBoxIconError)
									return
								}
							}
						},
					},
					Action{
						Text: i18n.Tr("mainwindow.menu.file.export"),
						OnTriggered: func() {
							dlg := &walk.FileDialog{
								Title: "save to file",
							}
							accepted, err := dlg.ShowSave(mw)
							if err != nil {
								walk.MsgBox(mw, "ERROR", "Open FileDialog:"+err.Error(), walk.MsgBoxIconError)
								return
							}
							if accepted {
								if err := mw.TV_sessions.ExportSessions(dlg.FilePath); err != nil {
									walk.MsgBox(mw, "ERROR", "Export Session Error:"+err.Error(), walk.MsgBoxIconError)
									return
								}
							}
						},
					},
				},
			},
			Menu{
				Text: i18n.Tr("mainwindow.menu.edit"),
				Items: []MenuItem{
					Action{
						Text: i18n.Tr("mainwindow.menu.edit.clear"),
						OnTriggered: func() {
							mw.TW_pages.CurrentPage().content.ClearScreen()
						},
					},
				},
			},
			Menu{
				Text: i18n.Tr("mainwindow.menu.setting"),
				Items: []MenuItem{
					Action{
						Text:        i18n.Tr("mainwindow.menu.setting.theme"),
						OnTriggered: nil,
					},
					Action{
						Text:        i18n.Tr("mainwindow.menu.logpath"),
						OnTriggered: nil,
					},
				},
			},
			Menu{
				Text: i18n.Tr("mainwindow.menu.run"),
				Items: []MenuItem{
					Action{
						Text: i18n.Tr("mainwindow.menu.run.batch"),
						OnTriggered: func() {
							curTabpage := mw.TW_pages.CurrentPage()
							if curTabpage == nil {
								walk.MsgBox(mw, "INFO", "当前没有打开的会话", walk.MsgBoxIconInformation)
								return
							}
							batchRun(mw)
						},
					},
				},
			},
			Menu{
				Text: i18n.Tr("mainwindow.menu.help"),
				Items: []MenuItem{
					Action{
						Text: i18n.Tr("mainwindow.menu.help.source"),
						OnTriggered: func() {
							startPage("https://github.com/chenqinghe/redis-desktop")
						},
					},
					Action{
						Text:        i18n.Tr("mainwindow.menu.help.bug"),
						OnTriggered: startIssuePage,
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
							Label{Text: i18n.Tr("mainwindow.labelhost")},
							LineEdit{AssignTo: &mw.LE_host},
							Label{Text: i18n.Tr("mainwindow.labelport")},
							LineEdit{AssignTo: &mw.LE_port},
							Label{Text: i18n.Tr("mainwindow.labelpassword")},
							LineEdit{AssignTo: &mw.LE_password, PasswordMode: true},
							PushButton{
								Text:      i18n.Tr("mainwindow.PBconnect"),
								AssignTo:  &mw.PB_connect.PushButton,
								OnClicked: mw.PB_connect.OnClick,
							},
						},
					},
					Composite{
						Layout: HBox{MarginsZero: true},
						Children: []Widget{
							TreeView{
								AssignTo:   &mw.TV_sessions.TreeView,
								MaxSize:    Size{200, 0},
								ItemHeight: 20,
								Model:      NewSessionTreeModel(),
								ContextMenuItems: []MenuItem{
									Action{
										Text:        "添加会话",
										OnTriggered: mw.TV_sessions.AddSession,
									},
									Action{
										Text:        "添加目录",
										OnTriggered: mw.TV_sessions.AddDirectory,
									},
									Action{
										Text:        "删除会话",
										OnTriggered: mw.TV_sessions.RemoveSelectedSession,
									},
									Action{
										Text:        "删除目录",
										OnTriggered: mw.TV_sessions.RemoveSelectedDirectory,
									},
								},
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									switch button {
									case walk.LeftButton:
										item := mw.TV_sessions.ItemAt(x, y)
										mw.TV_sessions.SetCurrentItem(item)
									case walk.RightButton:
										actionList := mw.TV_sessions.ContextMenu().Actions()
										var showedMenu = make([]int, actionList.Len())
										switch item := mw.TV_sessions.CurrentItem(); item.(type) {
										case *Directory:
											showedMenu = []int{1, 1, 0, 1}
										case *Session:
											showedMenu = []int{0, 0, 1, 0}
										default: // nil
											showedMenu = []int{1, 1, 0, 0}
										}
										for i := 0; i < actionList.Len(); i++ {
											if showedMenu[i] == 1 {
												actionList.At(i).SetVisible(true)
											} else {
												actionList.At(i).SetVisible(false)
											}
										}
									}
								},
								OnItemActivated: func() {
									item := mw.TV_sessions.CurrentItem()
									switch t := item.(type) {
									case *Session:
										mw.TW_pages.startNewSession(*t)
									}
								},
							},
							TabWidget{
								AssignTo: &mw.TW_pages.TabWidget,
								Pages: []TabPage{
									TabPage{
										Title: "home",
										Image: "img/home.ico",
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

func batchRun(p *MainWindowEX) {
	var dlg *walk.Dialog
	var cmdContent *walk.TextEdit

	if _, err := (Dialog{
		Title:    "批量运行命令",
		AssignTo: &dlg,
		MinSize:  Size{500, 500},
		Layout: VBox{Margins: Margins{
			Left:   10, //int
			Top:    10, //int
			Right:  10, //int
			Bottom: 10, //int
		}},
		Children: []Widget{
			Label{Text: "请在下面输入要执行的命令，每行一条..."},
			TextEdit{
				AssignTo: &cmdContent,
			},
			Composite{
				Layout: HBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						Text: "确定",
						OnClicked: func() {
							content := cmdContent.Text()
							dlg.Close(0)
							cmds := strings.Split(content, "\r\n")
							curTabpage := p.TW_pages.CurrentPage()
							if curTabpage == nil {
								walk.MsgBox(p, "INFO", "当前没有打开的会话", walk.MsgBoxIconInformation)
								return
							}
							for _, v := range cmds {
								v = strings.TrimSpace(v)
								if len(v) > 0 {
									curTabpage.content.AppendText(v)
									curTabpage.content.runCmd(v)
								}
							}
						},
					},
					PushButton{
						Text: "取消",
						OnClicked: func() {
							dlg.Close(0)
						},
					},
				},
			},
		},
	}).Run(p); err != nil {
		logrus.Errorln("show batch run dialog error:", err)
	}
}
