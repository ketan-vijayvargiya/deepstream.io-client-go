package deepstreamio

type PresenceHandler struct {
    client          *Client
    clientConfig    *ClientConfig
    connection      *connection
}

func newPresenceHandler(client *Client, clientConfig *ClientConfig) *PresenceHandler {
    var e = &PresenceHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
