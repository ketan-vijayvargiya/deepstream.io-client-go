package deepstreamio

import "strings"

func getMsg(topic Topic, action Action, data []string) string {
    c := []string{string(topic), string(action)}
    for _, d := range data {
        c = append(c, d)
    }

    return strings.Join(c, messageUnitSeparator) + messageRecordSeparator
}
