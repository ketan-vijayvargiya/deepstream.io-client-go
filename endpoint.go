package deepstreamio

import (
    "github.com/gorilla/websocket"
    "github.com/Sirupsen/logrus"
)

type endpoint struct {
    url                 string
    connection          *connection

    isWebsocketClosed   chan bool
    websocketConn       *websocket.Conn
}

func newEndpoint(url string, connection *connection) *endpoint {
    return &endpoint{url: url, connection: connection, isWebsocketClosed: make(chan bool, 1)}
}

func (e *endpoint) send(msg string) {
    logrus.WithField("msg", msg).Debug("Sent message")
    go func() {
        var err = e.websocketConn.WriteMessage(websocket.TextMessage, []byte(msg))
        if err != nil {
            e.connection.onError(err.Error())
        }
    }()
}

func (e *endpoint) close(forceClose bool) {
    e.isWebsocketClosed <- true

    go func() {
        var err error
        if forceClose {
            err = e.websocketConn.Close()
        } else {
            err = e.websocketConn.WriteMessage(
                websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
        }

        if err != nil {
            e.connection.onError(err.Error())
        } else {
            e.connection.onClose()
        }
    }()
}

func (e *endpoint) open() {
    go func() {
        var conn, _, err = websocket.DefaultDialer.Dial(e.url, nil)

        if err != nil {
            e.connection.onError(err.Error())
        } else {
            logrus.Debug("Opened endpoint")

            conn.SetCloseHandler(e.websocketCloseHandler)
            e.websocketConn = conn

            go e.readMessagesInLoop()
            e.connection.onOpen()
        }
    }()
}

func (e *endpoint) websocketCloseHandler(code int, text string) error {
    e.isWebsocketClosed <- true
    defer e.connection.onClose()

    return nil
}

func (e *endpoint) readMessagesInLoop() {
    for {
        select {
        case <- e.isWebsocketClosed:
            return

        default:
            var _, rawMsg, err = e.websocketConn.ReadMessage()
            if err != nil {
                e.connection.onError(err.Error())
                return
            } else if rawMsgStr := string(rawMsg); len(rawMsgStr) > 0 {
                logrus.WithField("msg", rawMsgStr).Debug("Read message")
                e.connection.onMessage(rawMsgStr)
            }            
        }
    }
}
