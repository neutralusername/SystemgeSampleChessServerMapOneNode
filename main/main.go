package main

import (
	"SystemgeSampleChessServer/app"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Tools"
)

const LOGGER_PATH = "logs.log"

func main() {
	Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	Dashboard.New(&Config.Dashboard{
		NodeConfig: &Config.Node{
			Name:           "dashboard",
			RandomizerSeed: Tools.GetSystemTime(),
		},
		ServerConfig: &Config.TcpServer{
			Port: 8081,
		},
		NodeStatusIntervalMs: 1000,

		NodeSystemgeClientCounterIntervalMs:             1000,
		NodeSystemgeClientRateLimitCounterIntervalMs:    1000,
		NodeSystemgeClientConnectionCounterIntervalMs:   1000,
		NodeSystemgeClientAsyncMessageCounterIntervalMs: 1000,
		NodeSystemgeClientSyncResponseCounterIntervalMs: 1000,
		NodeSystemgeClientSyncRequestCounterIntervalMs:  1000,
		NodeSystemgeClientTopicCounterIntervalMs:        1000,

		NodeSystemgeServerCounterIntervalMs:             1000,
		NodeSystemgeServerRateLimitCounterIntervalMs:    1000,
		NodeSystemgeServerConnectionCounterIntervalMs:   1000,
		NodeSystemgeServerAsyncMessageCounterIntervalMs: 1000,
		NodeSystemgeServerSyncResponseCounterIntervalMs: 1000,
		NodeSystemgeServerSyncRequestCounterIntervalMs:  1000,
		NodeSystemgeServerTopicCounterIntervalMs:        1000,

		NodeWebsocketCounterIntervalMs: 1000,
		HeapUpdateIntervalMs:           1000,
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		GoroutineUpdateIntervalMs:      1000,
		AutoStart:                      true,
		AddDashboardToDashboard:        true,
	},
		Node.New(&Config.NewNode{
			NodeConfig: &Config.Node{
				Name:              "node",
				RandomizerSeed:    Tools.GetSystemTime(),
				InfoLoggerPath:    LOGGER_PATH,
				WarningLoggerPath: LOGGER_PATH,
				ErrorLoggerPath:   LOGGER_PATH,
			},
			WebsocketConfig: &Config.Websocket{
				Pattern: "/ws",
				ServerConfig: &Config.TcpServer{
					Port:      8443,
					Blacklist: []string{},
					Whitelist: []string{},
				},
				HandleClientMessagesSequentially: false,
				ClientWatchdogTimeoutMs:          20000,
				Upgrader: &websocket.Upgrader{
					ReadBufferSize:  1024,
					WriteBufferSize: 1024,
					CheckOrigin: func(r *http.Request) bool {
						return true
					},
				},
			},
			HttpConfig: &Config.HTTP{
				ServerConfig: &Config.TcpServer{
					Port: 8080,
				},
			},
		}, app.New()),
	).StartBlocking()
}
