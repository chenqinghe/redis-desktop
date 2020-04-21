package i18n

func init() {
	AddLang(en_us)
}

var en_us = Lang{
	name: "en_us",
	Translation: &Translation{
		sections: make([]string, 0),
		words: map[string]string{
			// for main window
			"mainwindow.title":                         "Redis Cli Desktop",
			"mainwindow.menu.file":                     "File",
			"mainwindow.menu.file.export":              "export",
			"mainwindow.menu.file.import":              "import",
			"mainwindow.menu.edit":                     "Edit",
			"mainwindow.menu.edit.clear":               "clear",
			"mainwindow.menu.setting":                  "Setting",
			"mainwindow.menu.setting.theme":            "theme",
			"mainwindow.menu.logpath":                  "Log Path",
			"mainwindow.menu.run":                      "Run",
			"mainwindow.menu.run.batch":                "run batch command",
			"mainwindow.menu.help":                     "Help",
			"mainwindow.menu.help.source":              "view source code",
			"mainwindow.menu.help.bug":                 "new issue",
			"mainwindow.menu.help.donate":              "donate",
			"mainwindow.LBsessions.menu.deletesession": "delete session",
			"mainwindow.labelhost":                     "Host",
			"mainwindow.labelport":                     "Port",
			"mainwindow.labelpassword":                 "Password",
			"mainwindow.PBconnect":                     "Connect",

			// for widiget
			"widget.button.yes":    "YES",
			"widget.button.no":     "NO",
			"widget.button.ok":     "OK",
			"widget.button.cancel": "Cancel",

			// textedit
			"widget.textedit.menu.execselected": "Run Selected Command",
			"widget.textedit.menu.copy":         "Copy",
			"widget.textedit.menu.clear":        "Clear",

			// treeview
			"widget.treeview.menu.opensession":     "Open",
			"widget.treeview.menu.addsession":      "New Session",
			"widget.treeview.menu.adddirectory":    "New Directory",
			"widget.treeview.menu.editsession":     "Edit Session",
			"widget.treeview.menu.editdirectory":   "Edit Directory",
			"widget.treeview.menu.deletesession":   "Remove Session",
			"widget.treeview.menu.deletedirectory": "Remove Directory",

			// alert msg
			"alert.loadsessionfailed": "Load Session Failed: %s",
			"alert.selectedcmdempty":  "No Selected Command",
			"alert.noopenedsession":   "No Opened Session",
			"alert.hostcannotempty":"Host cannot be empty",
			"alert.invalidport":"Invalid Port: %s",
		},
	},
}
