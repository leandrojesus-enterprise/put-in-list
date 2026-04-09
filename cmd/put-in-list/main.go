package main

import (
	"github.com/leandrojesus-enterprise/put-in-list/internal/app"
	"github.com/leandrojesus-enterprise/put-in-list/internal/config"
	"github.com/leandrojesus-enterprise/put-in-list/internal/i18n"
	"github.com/leandrojesus-enterprise/put-in-list/internal/installer"
	"github.com/leandrojesus-enterprise/put-in-list/internal/menu"
	"github.com/leandrojesus-enterprise/put-in-list/internal/storage"
	"github.com/leandrojesus-enterprise/put-in-list/internal/terminal"
)

func main() {
	var tr *i18n.Translator = i18n.New()
	var term *terminal.Terminal = terminal.New()
	var ui *menu.Menu = menu.New(term, tr)

	var cfgSvc *config.JSONService = config.NewJSONService("config.json")
	var listSvc *storage.JSONStore = storage.NewJSONStore("lists.json")
	var instSvc *installer.Service = installer.NewService()

	var a *app.App = app.New(cfgSvc, listSvc, instSvc, ui, term, tr)
	a.Run()
}
