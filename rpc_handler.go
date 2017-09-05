package deepstreamio

type RpcHandler struct {
    client          *Client
    clientConfig    *ClientConfig
    connection      *connection
}

func newRpcHandler(client *Client, clientConfig *ClientConfig) *RpcHandler {
    var e = &RpcHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
