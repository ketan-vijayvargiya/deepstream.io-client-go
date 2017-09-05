package deepstreamio

import (
    "strings"
    "strconv"
    "encoding/json"
)

func parse(message string, client *Client) []*Message {
    var messages = make([]*Message, 0)

    for _, rawMessage := range strings.Split(message, messageRecordSeparator) {
        if parsedMessage := parseMessage(rawMessage, client); parsedMessage != nil {
            messages = append(messages, parsedMessage)
        }
    }

    return messages
}

func parseMessage(message string, client *Client) *Message {
    if len(message) == 0 {
        return nil
    }

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
    return &Message{Topic: Topic(parts[0]), Action: Action(parts[1]), Data: parts[2:]}
}

func convertTyped(value string, client *Client) interface{} {
    switch Type(value[0]) {
    case Type_String:
        return value[1:]

    case Type_Null:
        return nil

    case Type_Number:
        var f, _ = strconv.ParseFloat(string(value[1:]), 64)
        return f

    case Type_True:
        return true

    case Type_False:
        return false

    case Type_Object:
        return parseObject(string(value[1:]))

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
