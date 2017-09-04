package deepstreamio

const (
    messageUnitSeparator   = "\u001f"
    messageRecordSeparator = "\u001e"
)

type Message struct {
    Topic  Topic
    Action Action
    Data   []string
}
