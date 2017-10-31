package deepstreamio

import "testing"

var testTopic = Topic_Connection
var testAction = Action_Error
var testMessagePrefix = string(testTopic) + messageUnitSeparator + string(testAction)

type testStruct struct {
    A   string
    B   string
    c   string
}

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

func TestTyped(t *testing.T)  {
    ValidateStringsEqual(t, typed("foo"), "Sfoo")
    ValidateStringsEqual(t, typed(5), "N5")
    ValidateStringsEqual(t, typed(true), "T")
    ValidateStringsEqual(t, typed(false), "F")
    ValidateStringsEqual(t, typed(nil), "L")

    v := &testStruct{A: "a_val", B: "b_val", c: "c_val"}
    // 'c' isn't included, because it isn't visible.
    ValidateStringsEqual(t, typed(v), "O{\"A\":\"a_val\",\"B\":\"b_val\"}")
}
