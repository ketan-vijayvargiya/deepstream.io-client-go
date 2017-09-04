package deepstreamio

import "testing"

var testTopic = Topic_Connection
var testAction = Action_Error
var testMessagePrefix = string(testTopic) + messageUnitSeparator + string(testAction)

func TestMessageToStringNil(t *testing.T) {
    var expected = testMessagePrefix + messageRecordSeparator

    ValidateStringsEqual(t, getMsg(testTopic, testAction, nil), expected)
    ValidateStringsEqual(t, getMsg(testTopic, testAction, []string{}), expected)
}

func TestMessageToStringNonNil(t *testing.T) {
    ValidateStringsEqual(t,
        getMsg(testTopic, testAction, []string{"1", "2"}),
        testMessagePrefix + messageUnitSeparator+ "1" + messageUnitSeparator+ "2" + messageRecordSeparator)
}
