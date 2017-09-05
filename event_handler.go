package deepstreamio

type EventHandler struct {
    client          *Client
    clientConfig    *ClientConfig
    connection      *connection
}

func newEventHandler(client *Client, clientConfig *ClientConfig) *EventHandler {
    var e = &EventHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
