package deepstreamio

type RpcHandler struct {
    client          *client
    clientConfig    *ClientConfig
    connection      *connection
}

func newRpcHandler(client *client, clientConfig *ClientConfig) *RpcHandler {
    var e = &RpcHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
