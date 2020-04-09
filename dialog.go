package main

import (
	"github.com/chenqinghe/redis-desktop/i18n"
	"github.com/chenqinghe/walk"
	. "github.com/chenqinghe/walk/declarative"
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
						Text: i18n.Tr("widiget.button.ok"),
						OnClicked: func() {
							content = input.Text()
							dlg.Close(0)
						},
					},
					PushButton{
						Text: i18n.Tr("widiget.button.cancel"),
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

// Custom show a custom dialog, only the confirm button was pushed, accepted become true, otherwise
// no matter the dialog was closed or cancel button was pushed, accepted returns false.
func (sd *SimpleDialog) Custom(owner walk.Form, widget Widget) (accepted bool, err error) {
	var (
		dlg *walk.Dialog
	)

	if _, err := (Dialog{
		AssignTo: &dlg,
		Layout:   VBox{Margins: Margins{}},
		Children: []Widget{
			widget,
			Composite{
				Layout: HBox{Margins: Margins{}},
				Children: []Widget{
					PushButton{
						Text: i18n.Tr("widiget.button.ok"),
						OnClicked: func() {
							// some stuff here...
							dlg.Close(0)
						},
					},
					PushButton{
						Text: i18n.Tr("widiget.button.cancel"),
						OnClicked: func() {
							dlg.Close(0)
						},
					},
				},
			},
		},
		Title:     sd.Title,
		Size:      sd.Size,
		FixedSize: sd.FixedSize,
	}).Run(owner); err != nil {
		return false, err
	}

	return
}
