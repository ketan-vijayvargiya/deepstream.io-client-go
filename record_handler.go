package deepstreamio

type RecordHandler struct {
    client          *Client
    clientConfig    *ClientConfig
    connection      *connection
}

func newRecordHandler(client *Client, clientConfig *ClientConfig) *RecordHandler {
    var e = &RecordHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
