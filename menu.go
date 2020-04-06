package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"reflect"
)

func buildMenu(declMenu declarative.Menu) (*walk.Menu, error) {
	menu, err := walk.NewMenu()
	if err != nil {
		return nil, err
	}

	for _, v := range declMenu.Items {
		switch t := v.(type) {
		case declarative.Menu:
			m, err := buildMenu(t)
			if err != nil {
				return nil, err
			}
			menu.Actions().AddMenu(m)
		case *declarative.Menu:
			m, err := buildMenu(*t)
			if err != nil {
				return nil, err
			}
			menu.Actions().AddMenu(m)
		case declarative.Action:
			action := walk.NewAction()
			action.SetText(t.Text)
			action.Triggered().Attach(t.OnTriggered)
			menu.Actions().Add(action)
		case *declarative.Action:
			action := walk.NewAction()
			action.SetText(t.Text)
			action.Triggered().Attach(t.OnTriggered)
			menu.Actions().Add(action)
		case declarative.Separator, *declarative.Separator:
			menu.Actions().Add(walk.NewSeparatorAction())
		default:
			panic("unsupported MenuItem:" + reflect.TypeOf(v).String())
		}
	}

	return menu, nil
}
