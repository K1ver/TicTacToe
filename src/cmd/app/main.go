package main

import (
	"app/internal/di"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(di.Module)
	app.Run()
}
