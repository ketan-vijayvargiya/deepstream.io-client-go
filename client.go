package deepstreamio

func NewClient(url string, clientConfig *ClientConfig) *Client {

    var clonedClientConfig  = clientConfig.cloneWithDefaults()

    var client = &Client{
        url                 : url,
        clientConfig        : clonedClientConfig,
    }
    client.connection       = newConnection(url, clonedClientConfig, client)

    client.EventHandler     = newEventHandler(client, clonedClientConfig)
    client.RpcHandler       = newRpcHandler(client, clonedClientConfig)
    client.RecordHandler    = newRecordHandler(client, clonedClientConfig)
    client.PresenceHandler  = newPresenceHandler(client, clonedClientConfig)

    return client
}

type Client struct {
    url                 string
    clientConfig        *ClientConfig
    connection          *connection

    RuntimeErrorHandler func (topic Topic, event Event, errorMessage string)

    EventHandler        *EventHandler
    RpcHandler          *RpcHandler
    RecordHandler       *RecordHandler
    PresenceHandler     *PresenceHandler
}

func (c *Client) Login(authParams string) *LoginResult {
    loginResultChan := make(chan *LoginResult)
    c.connection.authenticate(authParams, func(loginResult *LoginResult) {
        loginResultChan <- loginResult
    })

    return <- loginResultChan
}

func (c *Client) Close() {
    c.connection.close(false)
}

func (c *Client) AddConnectionChangeListener(listener ConnectionStateListener) {
    c.connection.connectionStateListeners = append(c.connection.connectionStateListeners, listener)
}

func (c *Client) RemoveConnectionChangeListener(listener ConnectionStateListener) {
    var arr = make([]ConnectionStateListener, len(c.connection.connectionStateListeners) - 1)
    for _, l := range c.connection.connectionStateListeners {
        if l != listener {
            arr = append(arr, l)
        }
    }
    c.connection.connectionStateListeners = arr
}

func (c *Client) GetConnectionState() ConnectionState {
    return c.connection.connectionState
}

func (c *Client) SetGlobalConnectivityState(state GlobalConnectivityState) {
    c.connection.setGlobalConnectivityState(state)
}

func (c *Client) onError(topic Topic, event Event, msg string) {
    if event == Event_AckTimeout || event == Event_ResponseTimeout {
        if c.connection.connectionState == ConnectionState_AwaitingAuthentication {
            c.onError(Topic_Error, Event_NotAuthenticated,
                "Your message timed out because you're not authenticated. Have you called login()?")
            return
        }
    }

    if c.RuntimeErrorHandler != nil {
        c.RuntimeErrorHandler(topic, event, msg)
    } else {
        panic("Unhandled exception for Topic:" + string(topic) + ", Event:" + string(event) + ", msg:" + msg)
    }
}
