package deepstreamio

import (
    "bytes"
    "time"
)

type connection struct {
    url                             string
    clientConfig                    *ClientConfig
    client                          *Client

    connectionState                 ConnectionState
    globalConnectivityState         GlobalConnectivityState

    originalUrl                     string
    authParams                      string
    loginCallback                   func(loginResult *LoginResult)
    deliberateClose                 bool
    redirecting                     bool
    tooManyAuthAttempts             bool
    challengeDenied                 bool
    reconnectTimeout                *time.Timer
    reconnectionAttempt             int
    endpoint                        *endpoint
    messageBuffer                   *bytes.Buffer
    connectionStateListeners        []ConnectionStateListener
}

func newConnection(url string, clientConfig *ClientConfig, client *Client) *connection {
    var conn = &connection{url: url,
        clientConfig: clientConfig,
        client: client,

        connectionState: ConnectionState_Closed,
        globalConnectivityState: GlobalConnectivityState_Disconnected,

        originalUrl: url,
        reconnectTimeout: nil,
        messageBuffer: &bytes.Buffer{},
        connectionStateListeners: []ConnectionStateListener{},
    }
    conn.createEndpoint()
    return conn
}

func (c *connection) authenticate(authParams string, loginCallback func(loginResult *LoginResult)) {
    c.loginCallback = loginCallback
    c.authParams = authParams

    if c.tooManyAuthAttempts || c.challengeDenied {
        c.client.onError(Topic_Error, Event_IsClosed, "This client's connection was closed")
        c.loginCallback(getLoginResultFailure(Event_IsClosed, "This client's connection was closed"))
        return
    }

    if c.connectionState == ConnectionState_AwaitingAuthentication {
        c.sendAuthMessage()
    }
}

func (c *connection) send(message string) {
    if c.connectionState != ConnectionState_Open {
        c.messageBuffer.WriteString(message)
    } else {
        c.endpoint.send(message)
    }
}

func (c *connection) sendMsg(topic Topic, action Action, data []string) {
    c.send(getMsg(topic, action, data))
}

func (c *connection) sendAuthMessage()  {
    c.setState(ConnectionState_Authenticating)
    c.endpoint.send(getMsg(Topic_Auth, Action_Request, []string{c.authParams}))
}

func (c *connection) close(forceClose bool) {
    c.deliberateClose = true

    if forceClose && c.endpoint != nil {
        c.endpoint.close(true)
    } else if c.endpoint != nil {
        c.endpoint.close(false)
        c.endpoint = nil
    }

    if c.reconnectTimeout != nil {
        c.reconnectTimeout.Stop()
        c.reconnectTimeout = nil
    }
}

func (c *connection) onOpen() {
    c.setState(ConnectionState_AwaitingConnection)
}

func (c *connection) onError(err string) {
    c.setState(ConnectionState_Error)

    var timer = time.NewTimer(time.Second)
    go func() {
        <- timer.C
        c.client.onError(Topic(""), Event_ConnectionError, err)
    }()
}

func (c *connection) onMessage(rawMessage string) {
    var parsedMessages = parse(rawMessage, c.client)

    for _, message := range parsedMessages {
        switch message.Topic {
        case Topic_Connection:
            c.handleConnectionResponse(message)

        case Topic_Auth:
            c.handleAuthResponse(message)

        default:
            c.client.onError(Topic_Error, Event_UnsolicitedMessage, string(message.Action))
        }
    }
}

func (c *connection) onClose() {
    if c.redirecting {
        c.redirecting = false
        c.createEndpoint()
    } else if c.deliberateClose {
        c.setState(ConnectionState_Closed)
    } else {
        if c.originalUrl != c.url {
            c.url = c.originalUrl
            c.createEndpoint()
            return
        }
        c.tryReconnect()
    }
}

func (c *connection) handleConnectionResponse(message *Message) {
    switch message.Action {
    case Action_Ping:
        c.endpoint.send(getMsg(Topic_Connection, Action_Pong, nil))

    case Action_Ack:
        c.setState(ConnectionState_AwaitingAuthentication)

    case Action_Challenge:
        c.setState(ConnectionState_Challenging)
        c.endpoint.send(getMsg(Topic_Connection, Action_ChallengeResponse, []string{c.originalUrl}))

    case Action_Rejection:
        c.challengeDenied = true
        c.close(false)

    case Action_Redirect:
        c.url = message.Data[0]
        c.redirecting = true
        c.endpoint.close(false)
        c.endpoint = nil
    }
}

func (c *connection) handleAuthResponse(message *Message) {
    switch message.Action {
    case Action_Error:
        if message.Data[0] == string(Event_TooManyAuthAttempts) {
            c.deliberateClose = true;
            c.tooManyAuthAttempts = true;
        } else {
            c.authParams = "";
            c.setState(ConnectionState_AwaitingAuthentication);
        }

        if c.loginCallback != nil {
            c.loginCallback(getLoginResultFailure(
                Event(message.Data[0]), convertTyped(message.Data[1], c.client)))
        }

    case Action_Ack:
        c.setState(ConnectionState_Open)

        if c.messageBuffer.Len() > 0 {
            c.endpoint.send(c.messageBuffer.String())
            c.messageBuffer = &bytes.Buffer{}
        }

        if c.loginCallback != nil {
            c.loginCallback(getLoginResultSucess(
                convertTyped(message.Data[0], c.client)))
        }
    }
}

func (c *connection) setState(connectionState ConnectionState) {
    c.connectionState = connectionState

    for _, l := range c.connectionStateListeners {
        l.connectionStateChanged(connectionState)
    }

    if c.connectionState == ConnectionState_AwaitingAuthentication && c.authParams != "" {
        c.sendAuthMessage()
    }
}

func (c *connection) setGlobalConnectivityState(globalConnectivityState GlobalConnectivityState) {
    c.globalConnectivityState = globalConnectivityState

    if globalConnectivityState == GlobalConnectivityState_Connected {
        if c.connectionState == ConnectionState_Closed || c.connectionState == ConnectionState_Error {
            c.tryReconnect()
        }
    } else {
        if c.reconnectTimeout != nil {
            c.reconnectTimeout.Stop()
        }

        c.reconnectTimeout = nil
        c.reconnectionAttempt = 0
        c.endpoint.close(true)
        c.setState(ConnectionState_Closed)
    }
}

func (c *connection) createEndpoint() {
    c.endpoint = newEndpoint(c.url, c)
    c.endpoint.open()
}

func (c *connection) tryReconnect() {
    if c.reconnectTimeout != nil {
        return
    }

    if c.reconnectionAttempt < c.clientConfig.MaxReconnectAttempts {
        if c.globalConnectivityState == GlobalConnectivityState_Connected {
            c.setState(ConnectionState_Reconnecting)

            var delayTimeMillis = min(c.clientConfig.ReconnectIntervalIncrement * c.reconnectionAttempt,
                c.clientConfig.MaxReconnectInterval)
            c.reconnectTimeout = time.NewTimer(time.Duration(delayTimeMillis) * time.Millisecond)
            go func() {
                <-c.reconnectTimeout.C
                c.tryOpen()
            }()

            c.reconnectionAttempt++
        }
    } else {
        c.clearReconnect();
        c.close(true);
    }
}

func (c *connection) tryOpen() {
    c.reconnectTimeout.Stop()
    c.reconnectTimeout = nil
    c.endpoint.open()
}

func (c *connection) clearReconnect() {
    c.reconnectTimeout = nil
    c.reconnectionAttempt = 0
}
