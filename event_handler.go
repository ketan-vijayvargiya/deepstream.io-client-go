package deepstreamio

type EventHandler struct {
    client          *client
    clientConfig    *ClientConfig
    connection      *connection
}

func newEventHandler(client *client, clientConfig *ClientConfig) *EventHandler {
    var e = &EventHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
