package deepstreamio

type utilResubscribeNotifier struct {
    client          *Client
    resubscribe     func()
    isReconnecting  bool
}

func newUtilResubscribeNotifier(client *Client, resubscribe func()) *utilResubscribeNotifier {
    u := &utilResubscribeNotifier{client: client, resubscribe: resubscribe}
    client.AddConnectionChangeListener(u)

    return u
}

func (u *utilResubscribeNotifier) destroy() {
    u.client.RemoveConnectionChangeListener(u)
    u.client = nil
    u.resubscribe = nil
}

func (u *utilResubscribeNotifier) ConnectionStateChanged(state ConnectionState) {
    if state == ConnectionState_Reconnecting && !u.isReconnecting {
        u.isReconnecting = true
    }
    if state == ConnectionState_Open && u.isReconnecting {
        u.isReconnecting = false
        u.resubscribe()
    }
}
