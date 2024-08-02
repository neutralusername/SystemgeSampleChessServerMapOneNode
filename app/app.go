package app

import (
	"sync"

	"github.com/neutralusername/Systemge/Node"
)

type App struct {
	games map[string]*ChessGame
	mutex sync.Mutex
}

func New() *App {
	return &App{
		games: make(map[string]*ChessGame),
	}
}

func (app *App) GetCommandHandlers() map[string]Node.CommandHandler {
	return map[string]Node.CommandHandler{}
}
