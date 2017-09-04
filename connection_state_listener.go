package deepstreamio

type ConnectionStateListener interface {
    connectionStateChanged(connectionState ConnectionState)
}
