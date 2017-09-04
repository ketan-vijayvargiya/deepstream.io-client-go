package deepstreamio

const (
    messageUnitSeparator   = "\u001f"
    messageRecordSeparator = "\u001e"
)

type message struct {
    topic  Topic
    action Action
    data   []string
}
