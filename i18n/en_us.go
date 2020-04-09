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
			"widiget.button.yes": "YES",
			"widiget.button.no":  "NO",
			"widiget.button.ok":  "OK",
		},
	},
}
