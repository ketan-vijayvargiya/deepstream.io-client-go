package deepstreamio

import "github.com/chuckpreslar/emission"

type EventHandler struct {
    client                  *Client
    clientConfig            *ClientConfig
    connection              *connection

    subscriptionTimeout     int
    emitter                 *emission.Emitter
    listeners               map[string]*utilListener
    subscriptions           map[string]bool     // keeping map, instead of array, because removal is easier in a map.
    ackTimeoutRegistry      *utilAckTimeoutRegistry
}

func newEventHandler(client *Client, clientConfig *ClientConfig) *EventHandler {
    e := &EventHandler{client: client,
        clientConfig: clientConfig,
        connection: client.connection,
        subscriptionTimeout: clientConfig.SubscriptionTimeout,
        emitter: emission.NewEmitter(),
        listeners: make(map[string]*utilListener),
        subscriptions: make(map[string]bool, 0),
        ackTimeoutRegistry: client.utilAckTimeoutRegistry,
    }

    newUtilResubscribeNotifier(client, func() {
        for eventName, _ := range e.subscriptions {
            e.connection.sendMsg(Topic_Event, Action_Subscribe, []string{eventName})
        }
    })
    return e
}

func (e *EventHandler) Subscribe(eventName string, eventListener *EventListener) {
    if e.emitter.GetListenerCount(eventName) > 0 {
        e.subscriptions[eventName] = true
        e.ackTimeoutRegistry.add(Topic_Event, Action_Subscribe, eventName, "", nil, e.subscriptionTimeout)
        e.connection.send(getMsg(Topic_Event, Action_Subscribe, []string{eventName}))
    }
    e.emitter.On(eventName, eventListener)
}

func (e *EventHandler) Unsubscribe(eventName string, eventListener *EventListener) {
    delete(e.subscriptions, eventName)
    e.emitter.Off(eventName, eventListener)

    if e.emitter.GetListenerCount(eventName) > 0 {
        e.ackTimeoutRegistry.add(Topic_Event, Action_Unsubscribe, eventName, "", nil, e.subscriptionTimeout)
        e.connection.send(getMsg(Topic_Event, Action_Unsubscribe, []string{eventName}))
    }
}

func (e *EventHandler) Emit(eventName string, data interface{}) {
    dataArr := []string{eventName}
    if data != nil {
        dataArr = append(dataArr, typed(data))
    }

    e.connection.send(getMsg(Topic_Event, Action_Event, dataArr))
    e.broadcastEvent(eventName, data)
}

func (e *EventHandler) Listen(pattern string, listenListener ListenListener) {
    if e.listeners[pattern] != nil {
        e.client.onError(Topic_Event, Event_ListenerExists, pattern)
    } else {
        eventListener := newUtilListener(Topic_Event, pattern, listenListener, e.client, e.clientConfig)
        e.listeners[pattern] = eventListener
        eventListener.start()
    }
}

func (e *EventHandler) Unlisten(pattern string) {
    listener := e.listeners[pattern]
    if listener != nil {
        e.ackTimeoutRegistry.add(Topic_Event, Action_Unlisten, pattern, "", nil, e.subscriptionTimeout)
        listener.destroy()
        delete(e.listeners, pattern)
    } else {
        e.client.onError(Topic_Event, Event_NotListening, pattern)
    }
}

func (e *EventHandler) handle(message *Message) {
    var eventName string

    if message.Action == Action_Ack {
        eventName = message.Data[1]
    } else {
        eventName = message.Data[0]
    }

    if message.Action == Action_Event {
        if len(message.Data) == 2 {
            e.broadcastEvent(eventName, convertTyped(message.Data[1], e.client))
        } else {
            e.broadcastEvent(eventName, nil)
        }
    } else if e.listeners[eventName] != nil {
        e.listeners[eventName].onMessage(message)
    } else if message.Action == Action_Ack {
        e.ackTimeoutRegistry.clearMsg(message)
    } else if message.Action == Action_Error {
        e.client.onError(Topic_Event, Event(message.Data[0]), message.Data[1])
    } else {
        e.client.onError(Topic_Event, Event_UnsolicitedMessage, eventName)
    }
}

func (e *EventHandler) broadcastEvent(eventName string, args interface{})()  {
    e.emitter.Emit(eventName, args)
}
