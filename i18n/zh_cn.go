package i18n

func init() {
	AddLang(zh_ch)
}

var zh_ch = Lang{
	name: "zh_cn",
	Translation: &Translation{
		sections: make([]string, 0),
		words: map[string]string{
			"mainwindow.title":                         "redis命令行工具",
			"mainwindow.menu.file":                     "文件",
			"mainwindow.menu.file.export":              "导出会话",
			"mainwindow.menu.file.import":              "导入会话",
			"mainwindow.menu.edit":                     "编辑",
			"mainwindow.menu.edit.clear":               "清屏",
			"mainwindow.menu.setting":                  "设置",
			"mainwindow.menu.setting.theme":            "主题",
			"mainwindow.menu.logpath":                  "日志路径",
			"mainwindow.menu.run":                      "运行",
			"mainwindow.menu.run.batch":                "在当前会话批量运行",
			"mainwindow.menu.help":                     "帮助",
			"mainwindow.menu.help.source":              "查看源码",
			"mainwindow.menu.help.bug":                 "提bug",
			"mainwindow.menu.help.donate":              "捐赠",
			"mainwindow.LBsessions.menu.deletesession": "删除会话",
			"mainwindow.labelhost":                     "主机",
			"mainwindow.labelport":                     "端口",
			"mainwindow.labelpassword":                 "密码",
			"mainwindow.PBconnect":                     "连接",
		},
	},
}
