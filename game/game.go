package game

import (
	"encoding/json"
	"github.com/dkpeakbil/taskserver/domain"
	"github.com/dkpeakbil/taskserver/usecase"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Game struct {
	addr  string
	ucase usecase.UseCase
	conn  map[*websocket.Conn]*GameConn
}

type GameConn struct {
	conn          *websocket.Conn
	authenticated bool
}

type GameCmd struct {
	Command string `json:"cmd"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:   2048,
	WriteBufferSize:  2048,
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
	EnableCompression: false,
}

func NewGame(addr string, ucase usecase.UseCase) (*Game, error) {
	return &Game{
		addr:  addr,
		ucase: ucase,
		conn:  make(map[*websocket.Conn]*GameConn),
	}, nil
}

func (g *Game) Run() error {
	http.HandleFunc("/", g.handleWs)
	return http.ListenAndServe(g.addr, nil)
}

func (g *Game) handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			logrus.Debugf("err closing ws conn %s", err)
		}

		delete(g.conn, conn)
	}()

	_, ok := g.conn[conn]
	if !ok {
		g.conn[conn] = &GameConn{
			conn:          conn,
			authenticated: false,
		}
	}

	g.listen(conn)
}

func (g *Game) listen(conn *websocket.Conn) {
	client, ok := g.conn[conn]
	if !ok {
		return
	}

	for {
		command := &GameCmd{}
		_, msg, err := conn.ReadMessage()
		if err != nil || err == io.EOF {
			logrus.Debugf("error reading %s", err)
			conn.WriteMessage(websocket.CloseMessage, []byte("bye"))
			break
		}

		if err := json.Unmarshal(msg, &command); err != nil {
			logrus.Debugf("invalid command")
			continue
		}

		switch command.Command {
		case "auth":
			if client.authenticated {
				continue
			}

			var authRequest *domain.AuthGameRequest
			if err := json.Unmarshal(msg, &authRequest); err != nil {
				logrus.Debugf("error auth request: %s", err)
				continue
			}

			res := g.ucase.AuthGame(authRequest)
			if res.Status {
				client.authenticated = true
			}

			r, _ := json.Marshal(res)
			err := client.conn.WriteMessage(websocket.TextMessage, r)
			if err != nil {
				logrus.Debugf("error writing client %s", err)
			}

		case "getusers":
			if !client.authenticated {
				continue
			}

			var getUsersRequest *domain.GetUsersRequest
			if err := json.Unmarshal(msg, &getUsersRequest); err != nil {
				logrus.Debugf("error get users request: %s", err)
				continue
			}

			res := g.ucase.GetUsers(getUsersRequest)
			r, _ := json.Marshal(res)
			err := client.conn.WriteMessage(websocket.TextMessage, r)
			if err != nil {
				logrus.Debugf("error writing client %s", err)
			}

		default:
			logrus.Debugf("invalid command")
		}
	}
}
