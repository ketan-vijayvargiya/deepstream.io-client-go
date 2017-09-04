package deepstreamio

type ClientConfig struct {
    Path                                string
    ReconnectIntervalIncrement          int
    MaxReconnectInterval                int
    MaxReconnectAttempts                int
    RpcAckTimeout                       int
    RpcResponseTimeout                  int
    SubscriptionTimeout                 int
    MaxMessagesPerPacket                int
    TimeBetweenSendingQueuedPackages    int
    RecordReadAckTimeout                int
    RecordReadTimeout                   int
    RecordDeleteTimeout                 int
    RecordMergeStrategy                 RecordMergeStrategy
}

func (orig *ClientConfig) cloneWithDefaults() *ClientConfig {
    var cloned = &ClientConfig{}

    cloned.Path                             = getStringOrDefault(orig.Path, "/deepstream")
    cloned.ReconnectIntervalIncrement       = getIntOrDefault(orig.ReconnectIntervalIncrement, 4000)
    cloned.MaxReconnectInterval             = getIntOrDefault(orig.MaxReconnectInterval, 1500)
    cloned.MaxReconnectAttempts             = getIntOrDefault(orig.MaxReconnectAttempts, 5)
    cloned.RpcAckTimeout                    = getIntOrDefault(orig.RpcAckTimeout, 6000)
    cloned.RpcResponseTimeout               = getIntOrDefault(orig.RpcResponseTimeout, 10000)
    cloned.SubscriptionTimeout              = getIntOrDefault(orig.SubscriptionTimeout, 2000)
    cloned.MaxMessagesPerPacket             = getIntOrDefault(orig.MaxMessagesPerPacket, 100)
    cloned.TimeBetweenSendingQueuedPackages = getIntOrDefault(orig.TimeBetweenSendingQueuedPackages, 16)
    cloned.RecordReadAckTimeout             = getIntOrDefault(orig.RecordReadAckTimeout, 1000)
    cloned.RecordReadTimeout                = getIntOrDefault(orig.RecordReadTimeout, 3000)
    cloned.RecordDeleteTimeout              = getIntOrDefault(orig.RecordDeleteTimeout, 3000)

    if isRecordMergeStrategyValid(string(orig.RecordMergeStrategy)) {
        cloned.RecordMergeStrategy = orig.RecordMergeStrategy
    } else {
        cloned.RecordMergeStrategy = RecordMergeStrategy_RemoteWins
    }

    return cloned
}
