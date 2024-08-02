package app

import (
	"SystemgeSampleChessServer/topics"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		topics.STARTGAME: func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			whiteId := websocketClient.GetId()
			blackId := message.GetPayload()
			if !node.WebsocketClientExists(blackId) {
				return Error.New("Opponent does not exist", nil)
			}
			app.mutex.Lock()
			defer app.mutex.Unlock()
			if app.games[whiteId] != nil || app.games[blackId] != nil {
				return Error.New("Already in a game", nil)
			}
			err := node.AddToWebsocketGroup(whiteId+"-"+blackId, whiteId, blackId)
			if err != nil {
				return Error.New("Error adding to group", err)
			}
			game := newChessGame(whiteId, blackId)
			app.games[whiteId] = game
			app.games[blackId] = game
			node.WebsocketGroupcast(whiteId+"-"+blackId, Message.NewAsync(topics.STARTGAME, game.marshalBoard()))
			return nil
		},
		topics.ENDGAME: func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			game := app.games[websocketClient.GetId()]
			if game == nil {
				return Error.New("Game does not exist", nil)
			}
			delete(app.games, game.whiteId)
			delete(app.games, game.blackId)
			node.WebsocketGroupcast(game.whiteId+"-"+game.blackId, Message.NewAsync(topics.ENDGAME, ""))
			node.RemoveFromWebsocketGroup(game.whiteId+"-"+game.blackId, game.whiteId, game.blackId)
			return nil
		},
		topics.MOVE: func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			move, err := UnmarshalMove(message.GetPayload())
			if err != nil {
				return Error.New("Error unmarshalling move", err)
			}
			move.PlayerId = websocketClient.GetId()
			app.mutex.Lock()
			defer app.mutex.Unlock()
			game := app.games[move.PlayerId]
			if game == nil {
				return Error.New("Game does not exist", nil)
			}
			move, err = game.handleMoveRequest(move)
			if err != nil {
				return err
			}
			node.WebsocketGroupcast(game.whiteId+"-"+game.blackId, Message.NewAsync(topics.MOVE, Helpers.JsonMarshal(move)))
			return nil
		},
	}
}

func (game *ChessGame) handleMoveRequest(move *Move) (*Move, error) {
	if game.isWhiteTurn() && move.PlayerId != game.whiteId {
		return nil, Error.New("Not your turn", nil)
	}
	if !game.isWhiteTurn() && move.PlayerId != game.blackId {
		return nil, Error.New("Not your turn", nil)
	}
	chessMove, err := game.move(move)
	if err != nil {
		return nil, Error.New("Invalid move", err)
	}
	return chessMove, nil
}

func (app *App) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := websocketClient.Send(Message.NewAsync("connected", websocketClient.GetId()).Serialize())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Error sending connected message", err).Error())
		}
		websocketClient.Disconnect()
	}
}

func (app *App) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	game := app.games[websocketClient.GetId()]
	if game != nil {
		delete(app.games, game.whiteId)
		delete(app.games, game.blackId)

		node.WebsocketGroupcast(game.whiteId+"-"+game.blackId, Message.NewAsync(topics.ENDGAME, ""))
		node.RemoveFromWebsocketGroup(game.whiteId+"-"+game.blackId, game.whiteId, game.blackId)
	}
}
