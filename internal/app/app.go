package app

import (
	"github.com/stalwartgiraffe/cmr/internal/elog"
	"github.com/stalwartgiraffe/cmr/internal/otel"
)

type App struct {
	otel.Otel

	Shutdowns

	elog.ElapsedLogger
}

type AppErr struct {
	App *App
	Err error
}

func NewApp() AppErr {
	return AppErr{
		App: &App{
			Shutdowns: newShutdowns(),

			ElapsedLogger: *elog.New(),
		},
	}
}
