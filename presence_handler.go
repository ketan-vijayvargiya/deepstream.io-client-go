package deepstreamio

type PresenceHandler struct {
    client          *client
    clientConfig    *ClientConfig
    connection      *connection
}

func newPresenceHandler(client *client, clientConfig *ClientConfig) *PresenceHandler {
    var e = &PresenceHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
