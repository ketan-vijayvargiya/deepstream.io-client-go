package deepstreamio

type Client interface {
    Login(authParams string) *LoginResult
    Close()
    AddConnectionChangeListener(listener ConnectionStateListener)
    RemoveConnectionChangeListener(listener ConnectionStateListener)
    GetConnectionState() ConnectionState
    SetGlobalConnectivityState(state GlobalConnectivityState)
}

func NewClient(url string, clientConfig *ClientConfig,
    runtimeErrorHandler func (topic Topic, event Event, errorMessage string)) Client {

    var clonedClientConfig  = clientConfig.cloneWithDefaults()

    var client = &client{
        url                 : url,
        clientConfig        : clonedClientConfig,
        runtimeErrorHandler : runtimeErrorHandler,
    }
    client.connection       = newConnection(url, clonedClientConfig, client)

    client.EventHandler     = newEventHandler(client, clonedClientConfig)
    client.RpcHandler       = newRpcHandler(client, clonedClientConfig)
    client.RecordHandler    = newRecordHandler(client, clonedClientConfig)
    client.PresenceHandler  = newPresenceHandler(client, clonedClientConfig)

    return client
}

type client struct {
    url                 string
    clientConfig        *ClientConfig
    connection          *connection
    runtimeErrorHandler func (topic Topic, event Event, errorMessage string)

    EventHandler        *EventHandler
    RpcHandler          *RpcHandler
    RecordHandler       *RecordHandler
    PresenceHandler     *PresenceHandler
}

func (c *client) Login(authParams string) *LoginResult {
    loginResultChan := make(chan *LoginResult)
    c.connection.authenticate(authParams, func(loginResult *LoginResult) {
        loginResultChan <- loginResult
    })

    return <- loginResultChan
}

func (c *client) Close() {
    c.connection.close(false)
}

func (c *client) AddConnectionChangeListener(listener ConnectionStateListener) {
    c.connection.connectionStateListeners = append(c.connection.connectionStateListeners, listener)
}

func (c *client) RemoveConnectionChangeListener(listener ConnectionStateListener) {
    var arr = make([]ConnectionStateListener, len(c.connection.connectionStateListeners) - 1)
    for _, l := range c.connection.connectionStateListeners {
        if l != listener {
            arr = append(arr, l)
        }
    }
    c.connection.connectionStateListeners = arr
}

func (c *client) GetConnectionState() ConnectionState {
    return c.connection.connectionState
}

func (c *client) SetGlobalConnectivityState(state GlobalConnectivityState) {
    c.connection.setGlobalConnectivityState(state)
}

func (c *client) onError(topic Topic, event Event, msg string) {
    if event == Event_AckTimeout || event == Event_ResponseTimeout {
        if c.connection.connectionState == ConnectionState_AwaitingAuthentication {
            c.onError(Topic_Error, Event_NotAuthenticated,
                "Your message timed out because you're not authenticated. Have you called login()?")
            return
        }
    }

    if c.runtimeErrorHandler != nil {
        c.runtimeErrorHandler(topic, event, msg)
    } else {
        panic("Unhandled exception for Topic:" + string(topic) + ", Event:" + string(event) + ", msg:" + msg)
    }
}
