package deepstreamio

type utilListener struct {
    topic               Topic
    client              *Client
    clientConfig        *ClientConfig
    connection          *connection
    resubscribeNotifier *utilResubscribeNotifier
    ackTimoutRegistry   *utilAckTimeoutRegistry
    pattern             string
    listenerCallback    ListenListener
}

func newUtilListener(topic Topic, pattern string, listenerCallback ListenListener, client *Client,
    clientConfig *ClientConfig) *utilListener {

    u := &utilListener{topic: topic, client: client, clientConfig: clientConfig, connection: client.connection,
        ackTimoutRegistry: client.utilAckTimeoutRegistry, pattern: pattern, listenerCallback: listenerCallback}
    u.resubscribeNotifier = newUtilResubscribeNotifier(client, u.resubscribe)

    return u
}

func (u *utilListener) start() {
    u.scheduleAckTimeout()
    u.sendListen()
}

func (u *utilListener) destroy() {
    u.connection.sendMsg(u.topic, Action_Unlisten, []string{u.pattern})
    u.resubscribeNotifier.destroy()
    u.listenerCallback = nil
    u.pattern = ""
    u.client = nil
    u.connection = nil
    u.ackTimoutRegistry = nil
}

func (u *utilListener) onMessage(message *Message) {
    if message.Action == Action_Ack {
        u.ackTimoutRegistry.clearMsg(message)
    } else {
        if message.Action == Action_SubscriptionForPatternFound {
            accepted := u.listenerCallback.OnSubscriptionForPatternAdded(message.Data[1])
            if (accepted) {
                u.sendAccept(message.Data[1]);
            } else {
                u.sendReject(message.Data[1]);
            }
        } else if message.Action == Action_SubscriptionForPatternRemoved {
            u.listenerCallback.OnSubscriptionForPatternRemoved(message.Data[1])
        }
    }
}

func (u *utilListener) sendListen() {
    u.connection.sendMsg(u.topic, Action_Listen, []string{u.pattern})
}

func (u *utilListener) sendAccept(subscription string) {
    u.connection.sendMsg(u.topic, Action_ListenAccept, []string{u.pattern, subscription})
}

func (u *utilListener) sendReject(subscription string) {
    u.connection.sendMsg(u.topic, Action_ListenReject, []string{u.pattern, subscription})
}

func (u *utilListener) scheduleAckTimeout() {
    u.ackTimoutRegistry.add(u.topic, Action_Listen, u.pattern, "", nil, u.clientConfig.SubscriptionTimeout)
}

func (u *utilListener) resubscribe() {
    u.scheduleAckTimeout()
    u.sendListen()
}
