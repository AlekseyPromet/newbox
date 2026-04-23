package main

import (
	"netbox_go/internal/app"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		app.ModuleInfra,
		app.ModuleRepository,
		app.ModuleAPI,
		app.ModuleServer,
	).Run()
}
