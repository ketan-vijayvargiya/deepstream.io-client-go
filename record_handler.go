package deepstreamio

type RecordHandler struct {
    client          *client
    clientConfig    *ClientConfig
    connection      *connection
}

func newRecordHandler(client *client, clientConfig *ClientConfig) *RecordHandler {
    var e = &RecordHandler{client: client, clientConfig: clientConfig, connection: client.connection}
    return e
}
