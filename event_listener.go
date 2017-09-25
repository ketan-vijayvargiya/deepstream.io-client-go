package deepstreamio

type EventListener interface {
    OnEvent(eventName string, args interface{})
}
