package deepstreamio

type ConnectionStateListener interface {
    ConnectionStateChanged(connectionState ConnectionState)
}
