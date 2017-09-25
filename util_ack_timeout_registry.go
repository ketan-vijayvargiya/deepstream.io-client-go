package deepstreamio

import "time"

type utilAckTimeoutRegistry struct {
    client      *Client
    register    map[string]*time.Timer
    ackTimers   []*ackTimeout
    state       ConnectionState
}

func newUtilAckTimeoutRegistry(client *Client) *utilAckTimeoutRegistry {
    atr := &utilAckTimeoutRegistry{
        client: client,
        register: make(map[string]*time.Timer),
        ackTimers: make([]*ackTimeout, 0),
        state: client.GetConnectionState(),
    }
    client.AddConnectionChangeListener(atr)
    return atr
}

func (u *utilAckTimeoutRegistry) clearMsg(message *Message) {
    var action Action
    var name string
    if message.Action == Action_Ack {
        action = Action(message.Data[0])
        name = message.Data[1]
    } else {
        action = message.Action
        name = message.Data[0]
    }

    uniqueName := u.getUniqueName(message.Topic, action, name)
    if ok := u.clear(uniqueName); !ok {
        u.client.onError(message.Topic, Event_UnsolicitedMessage, message.raw)
    }
}

func (u *utilAckTimeoutRegistry) add(topic Topic, action Action, name string,
    event Event, timeoutListener UtilTimeoutListener, timeout int) {

    if timeoutListener == nil {
        timeoutListener = u
    }

    if len(event) == 0 {
        event = Event_AckTimeout
    }

    uniqueName := u.getUniqueName(topic, action, name)
    u.clear(uniqueName)

    ackTimeout := &ackTimeout{topic: topic, action: action, name: name, event: event,
        timeoutListener: timeoutListener, timeout: timeout, client: u.client}

    if u.state == ConnectionState_Open {
        timer := scheduleFuncAfterMillis(ackTimeout.run, timeout)
        u.register[uniqueName] = timer
    } else {
        u.ackTimers = append(u.ackTimers, ackTimeout)
    }
}

func (u *utilAckTimeoutRegistry) ConnectionStateChanged(connectionState ConnectionState) {
    if connectionState == ConnectionState_Open {
        u.scheduleAcks()
    }
    u.state = connectionState
}

func (u *utilAckTimeoutRegistry) clear(uniqueName string) bool {
    timer, ok := u.register[uniqueName]
    if ok {
        timer.Stop()
        delete(u.register, uniqueName)
    }
    return ok
}

func (u *utilAckTimeoutRegistry) OnTimeout(topic Topic, action Action, event Event, name string) {
    uniqueName := u.getUniqueName(topic, action, name)
    delete(u.register, uniqueName)
}

func (u *utilAckTimeoutRegistry) scheduleAcks() {
    for _, a := range u.ackTimers {
        scheduleFuncAfterMillis(a.run, a.timeout)
    }
    u.ackTimers = make([]*ackTimeout, 0)    // Create a new slice so that the existing one becomes eligible to be GCed.
}

func (u *utilAckTimeoutRegistry) getUniqueName(topic Topic, action Action, name string) string {
    return string(topic) + string(action) + name
}

type ackTimeout struct {
    topic           Topic
    action          Action
    name            string
    event           Event
    timeoutListener UtilTimeoutListener
    timeout         int
    client          *Client
}

func (a *ackTimeout) run() {
    a.timeoutListener.OnTimeout(a.topic, a.action, a.event, a.name)

    var msg string
    if a.event == Event_AckTimeout {
        msg = "No ACK message received in time for " + string(a.action) + " " + a.name
    } else {
        msg = "No message received in time for " + string(a.action) + " " + a.name
    }
    a.client.onError(a.topic, a.event, msg)
}
