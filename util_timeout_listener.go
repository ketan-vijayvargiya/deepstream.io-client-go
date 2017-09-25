package deepstreamio

type UtilTimeoutListener interface {
    OnTimeout(topic Topic, action Action, event Event, name string)
}
