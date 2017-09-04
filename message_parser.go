package deepstreamio

import (
    "strings"
    "strconv"
    "encoding/json"
)

func parse(message string, client *client) []*message {
    var messages = []*message{}

    for _, rawMessage := range strings.Split(message, messageRecordSeparator) {
        if parsedMessage := parseMessage(rawMessage, client); parsedMessage != nil {
            messages = append(messages, parsedMessage)
        }
    }

    return messages
}

func parseMessage(message string, client *client) *message {
    var parts = strings.Split(message, messageUnitSeparator)

    if len(parts) < 2 {
        client.onError(Topic(""), Event_MessageParseError, "Insufficient message parts")
        return nil
    }
    if !isTopicValid(parts[0]) {
        client.onError("", Event_MessageParseError, "Received message for unknown topic " + parts[0])
        return nil
    }
    if !isActionValid(parts[1]) {
        client.onError("", Event_MessageParseError, "Unknown action " + parts[1])
        return nil
    }
    return &message{topic: Topic(parts[0]), action: Action(parts[1]), data: parts[2:]}
}

func convertTyped(value string, client *client) interface{} {
    switch value[0] {
    case Type_String:
        return value[1:]

    case Type_Null:
        return nil

    case Type_Number:
        var f, _ = strconv.ParseFloat(value[1:], 64)
        return f

    case Type_True:
        return true

    case Type_False:
        return false

    case Type_Object:
        return parseObject(value[1:])

    case Type_Undefined:
        return Type_Undefined
    }

    client.onError( Topic_Error, Event_MessageParseError, "UNKNOWN_TYPE (" + value + ")")
    return nil
}

func parseObject(value string) interface{} {
    var j interface{}
    json.Unmarshal([]byte(value), j)

    return j
}
