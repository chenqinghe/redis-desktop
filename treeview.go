package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lxn/walk"
	"github.com/sirupsen/logrus"
)

type Directory struct {
	parent   *Directory
	Name     string
	Children []walk.TreeItem
}

func NewDirectory(name string, parent *Directory) *Directory {
	return &Directory{Name: name, parent: parent}
}

func (d *Directory) Text() string {
	return d.Name
}

func (d *Directory) Parent() walk.TreeItem {
	if d.parent == nil {
		return nil
	}
	return d.parent
}

func (d *Directory) ChildCount() int {
	return len(d.Children)
}

func (d *Directory) ChildAt(i int) walk.TreeItem {
	return d.Children[i]
}

func (d *Directory) Image() interface{} {
	return "img/dir.ico"
}

type Session struct {
	Host     string
	Port     int
	Password string
	parent   *Directory
}

func (s *Session) Image() interface{} {
	return "img/redis.ico"
}

func (s *Session) Text() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *Session) Parent() walk.TreeItem {
	if s.parent == nil {
		return nil
	}
	return s.parent
}

func (s *Session) ChildCount() int {
	return 0
}

func (s *Session) ChildAt(i int) walk.TreeItem {
	return nil
}

var _ walk.TreeItem = new(Directory)
var _ walk.TreeItem = new(Session)

type SessionTreeModel struct {
	walk.TreeModelBase
	roots []walk.TreeItem
}

func (m *SessionTreeModel) RootCount() int {
	return len(m.roots)
}

func (m *SessionTreeModel) RootAt(index int) walk.TreeItem {
	return m.roots[index]
}

func NewSessionTreeModel() *SessionTreeModel {
	m := &SessionTreeModel{
		TreeModelBase: walk.TreeModelBase{},
		roots:         []walk.TreeItem{},
	}

	return m
}

type TreeViewEx struct {
	*walk.TreeView

	root  *MainWindowEX
	model *SessionTreeModel
}

type Facade struct {
	// Directory
	Name     string
	Children []*Facade

	// Session
	Host     string
	Port     int
	Password string
}

func (f Facade) IsDirectory() bool {
	return f.Name != ""
}

func (tv *TreeViewEx) LoadSession(data []byte) error {
	facades := make([]Facade, 0)
	if err := json.Unmarshal(data, &facades); err != nil {
		return err
	}

	for _, v := range facades {
		tv.model.roots = append(tv.model.roots, buildModel(nil, v))
	}

	return tv.SetModel(tv.model)
}

func (tv *TreeViewEx) SaveSession(file string) error {
	data, err := json.Marshal(tv.model.roots)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, os.ModePerm)
}

func buildModel(parent *Directory, facade Facade) walk.TreeItem {
	if !facade.IsDirectory() {
		return &Session{
			Host:     facade.Host,
			Port:     facade.Port,
			Password: facade.Password,
			parent:   parent,
		}
	}
	dir := &Directory{
		parent: parent,
		Name:   facade.Name,
	}
	for _, v := range facade.Children {
		dir.Children = append(dir.Children, buildModel(dir, *v))
	}

	return dir
}

func (tv *TreeViewEx) AddSession(s *Session) {
	tv.addSession(s)
	tv.ReloadModel()
	if err := tv.root.saveSessions(nil); err != nil {
		logrus.Errorln("save sessions error:", err)
	}
}
func (tv *TreeViewEx) addSession(s *Session) {
	//item := tv.CurrentItem()
	//var dir *Directory
	//if d, ok := item.(*Directory); ok {
	//	dir = d
	//} else {
	//	dir = d.parent
	//}
	//
	//s.parent = dir
	//dir.children = append(dir.children, s)
}

func (tv *TreeViewEx) AddSessions(sesses []Session) {
	//for _, sess := range sesses {
	//	tv.addSession(sess)
	//}
	//tv.ReloadModel()
	//tv.root.saveSessions(tv.sessions)
}

func (tv *TreeViewEx) GetSessions() []Session {
	return nil
	//return tv.sessions
}

func (tv *TreeViewEx) RemoveSelectedSession() {
	//index := tv.CurrentIndex()
	//if index > 0 {
	//	tv.sessions[index] = tv.sessions[len(tv.sessions)-1]
	//	tv.sessions = tv.sessions[:len(tv.sessions)-1]
	//	tv.Model[index] = tv.Model[len(tv.Model)-1]
	//	tv.Model = tv.Model[:len(tv.Model)-1]
	//	tv.ReloadModel()
	//	tv.root.saveSessions(tv.sessions)
	//}
}

func (tv *TreeViewEx) ReloadModel() {
	tv.SetModel(tv.model)
}

func (tv *TreeViewEx) AddDirectory() {
	var parent *Directory
	var item walk.TreeItem

	curItem := tv.CurrentItem()
	switch t := curItem.(type) {
	case *Directory:
		parent = t
	case *Session:
		walk.MsgBox(tv.root, "ERROR", "不能在会话中创建目录", walk.MsgBoxIconError)
		return
	default: // nil
		name := (&SimpleDialog{}).Prompt(tv.root, "请输入目录名称")
		if name == "" {
			return
		}
		item = &Directory{
			parent:   nil,
			Name:     name,
			Children: nil,
		}
		tv.model.roots = append(tv.model.roots, item)
		tv.ReloadModel()
		tv.EnsureVisible(item)
		return
	}

	name := (&SimpleDialog{}).Prompt(tv.root, "请输入目录名称")
	if name == "" {
		return
	}
	item = &Directory{
		parent:   parent,
		Name:     name,
		Children: nil,
	}
	parent.Children = append(parent.Children, item)

	tv.ReloadModel()
	tv.EnsureVisible(item)
}
