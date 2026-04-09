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
	tr := i18n.New()
	term := terminal.New()
	ui := menu.New(term, tr)

	cfgSvc := config.NewJSONService("config.json")
	listSvc := storage.NewJSONStore("lists.json")
	instSvc := installer.NewService()

	a := app.New(cfgSvc, listSvc, instSvc, ui, term, tr)
	a.Run()
}
