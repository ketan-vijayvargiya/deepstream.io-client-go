package deepstreamio

import (
    "strings"
    "strconv"
    "encoding/json"
    "github.com/Sirupsen/logrus"
)

func getMsg(topic Topic, action Action, data []string) string {
    c := []string{string(topic), string(action)}
    for _, d := range data {
        c = append(c, d)
    }

    return strings.Join(c, messageUnitSeparator) + messageRecordSeparator
}

func typed(vi interface{}) string {
    if vi == nil {
        return Type_Null
    }

    switch v := vi.(type) {
    case string:
        return string(Type_String) + v

    case bool:
        if v {
            return string(Type_True)
        } else {
            return string(Type_False)
        }

    case int:   // Need to cases for other numbers, such as uint, float etc.
        return string(Type_Number) + strconv.Itoa(v)

    default:
        b, err := json.Marshal(v)
        if err != nil {
            logrus.WithField("input", vi).Error("Error in marshalling object to JSON")
            return string(Type_Undefined)
        }
        return string(Type_Object) + string(b)
    }
}
