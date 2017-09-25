package deepstreamio

type ListenListener interface {
    OnSubscriptionForPatternAdded(subscription string) bool
    OnSubscriptionForPatternRemoved(subscription string)
}
