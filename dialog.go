package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type SimpleDialog struct {
	Title string

	Size      Size
	FixedSize bool
}

func (sd *SimpleDialog) Alert(p walk.Form, msg string) {

}

func (sd *SimpleDialog) Confirm(p walk.Form, msg string) bool {

	return true
}

func (sd *SimpleDialog) Prompt(p walk.Form, msg string) string {
	var dlg *walk.Dialog
	var input *walk.LineEdit
	var content string

	if _, err := (Dialog{
		AssignTo: &dlg,
		Layout:   VBox{Margins: Margins{}},
		Children: []Widget{
			LineEdit{AssignTo: &input},
			Composite{
				Layout: HBox{Margins: Margins{}},
				Children: []Widget{
					PushButton{
						Text: "确认",
						OnClicked: func() {
							content = input.Text()
							dlg.Close(0)
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
		Title:     msg,
		Size:      sd.Size,
		FixedSize: sd.FixedSize,
	}).Run(p); err != nil {
		return ""
	}

	return content
}
