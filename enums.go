package deepstreamio

type Action string
const (
    Action_Error Action = "E"
    Action_Ping = "PI"
    Action_Pong = "PO"
    Action_Ack = "A"
    Action_Redirect = "RED"
    Action_Challenge = "CH"
    Action_ChallengeResponse = "CHR"
    Action_Read = "R"
    Action_Create = "C"
    Action_CreateOrRead = "CR"
    Action_CreateAndUpdate = "CU"
    Action_Update = "U"
    Action_Patch = "P"
    Action_Delete = "D"
    Action_Subscribe = "S"
    Action_Unsubscribe = "US"
    Action_Has = "H"
    Action_Snapshot = "SN"
    Action_SubscriptionForPatternFound = "SP"
    Action_SubscriptionForPatternRemoved = "SR"
    Action_SubscriptionHasProvider = "SH"
    Action_Listen = "L"
    Action_Unlisten = "UL"
    Action_ListenAccept = "LA"
    Action_ListenReject = "LR"
    Action_Event = "EVT"
    Action_Request = "REQ"
    Action_Response = "RES"
    Action_Rejection = "REJ"
    Action_PresenceJoin = "PNJ"
    Action_PresenceLeave = "PNL"
    Action_Query = "Q"
    Action_WriteAcknowledgement = "WA"
)
func isActionValid(value string) bool {
    for _, t := range []Action{Action_Error,
        Action_Ping,
        Action_Pong,
        Action_Ack,
        Action_Redirect,
        Action_Challenge,
        Action_ChallengeResponse,
        Action_Read,
        Action_Create,
        Action_CreateOrRead,
        Action_CreateAndUpdate,
        Action_Update,
        Action_Patch,
        Action_Delete,
        Action_Subscribe,
        Action_Unsubscribe,
        Action_Has,
        Action_Snapshot,
        Action_SubscriptionForPatternFound,
        Action_SubscriptionForPatternRemoved,
        Action_SubscriptionHasProvider,
        Action_Listen,
        Action_Unlisten,
        Action_ListenAccept,
        Action_ListenReject,
        Action_Event,
        Action_Request,
        Action_Response,
        Action_Rejection,
        Action_PresenceJoin,
        Action_PresenceLeave,
        Action_Query,
        Action_WriteAcknowledgement} {

        if string(t) == value {
            return true
        }
    }
    return false
}

type ConnectionState string
const (
    ConnectionState_Closed ConnectionState = "CLOSED"
    ConnectionState_AwaitingConnection = "AWAITING_CONNECTION"
    ConnectionState_Challenging = "CHALLENGING"
    ConnectionState_AwaitingAuthentication = "AWAITING_AUTHENTICATION"
    ConnectionState_Authenticating = "AUTHENTICATING"
    ConnectionState_Open = "OPEN"
    ConnectionState_Error = "ERROR"
    ConnectionState_Reconnecting = "RECONNECTING"
)

type GlobalConnectivityState string
const (
    GlobalConnectivityState_Connected GlobalConnectivityState = "CONNECTED"
    GlobalConnectivityState_Disconnected = "DISCONNECTED"
)

type Event string
const (
    Event_UnauthenticatedConnectiontimeout Event = "UNAUTHENTICATED_CONNECTION_TIMEOUT"
    Event_ConnectionError = "CONNECTION_ERROR"
    Event_ConnectionStateChanged = "CONNECTION_STATE_CHANGED"
    Event_AckTimeout = "ACK_TIMEOUT"
    Event_InvalidAuthData = "INVALID_AUTH_DATA"
    Event_ResponseTimeout = "RESPONSE_TIMEOUT"
    Event_CacheRetrievalTimeout = "CACHE_RETRIEVAL_TIMEOUT"
    Event_StorageRetrievalTimeout = "STORAGE_RETRIEVAL_TIMEOUT"
    Event_DeleteTimeout = "DELETE_TIMEOUT"
    Event_UnsolicitedMessage = "UNSOLICITED_MESSAGE"
    Event_MessageParseError = "MESSAGE_PARSE_ERROR"
    Event_VersionExists = "VERSION_EXISTS"
    Event_NotAuthenticated = "NOT_AUTHENTICATED"
    Event_ListenerExists = "LISTENER_EXISTS"
    Event_NotListening = "NOT_LISTENING"
    Event_TooManyAuthAttempts = "TOO_MANY_AUTH_ATTEMPTS"
    Event_IsClosed = "IS_CLOSED"
    Event_RecordNotFound = "RECORD_NOT_FOUND"
    Event_MessageDenied = "MESSAGE_DENIED"
    Event_MultipleSubscriptions = "MULTIPLE_SUBSCRIPTIONS"
)

type Topic string
const (
    Topic_Connection Topic = "C"
    Topic_Auth = "A"
    Topic_Error = "X"
    Topic_Event = "E"
    Topic_Record = "R"
    Topic_RPC = "P"
    Topic_Presence = "U"
)
func isTopicValid(value string) bool {
    for _, t := range []Topic{Topic_Connection,
        Topic_Auth,
        Topic_Error,
        Topic_Event,
        Topic_Record,
        Topic_RPC,
        Topic_Presence} {

        if string(t) == value {
            return true
        }
    }
    return false
}

type Type string
const (
    Type_String Type = "S"
    Type_Object = "O"
    Type_Number = "N"
    Type_Null = "L"
    Type_True = "T"
    Type_False = "F"
    Type_Undefined = "U"
)

type RecordMergeStrategy string
const (
    RecordMergeStrategy_RemoteWins RecordMergeStrategy = "REMOTE_WINS"
    RecordMergeStrategy_LocalWins = "LOCAL_WINS"
)
