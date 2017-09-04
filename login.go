package deepstreamio

type LoginResult struct {
    LoggedIn    bool
    ErrorEvent  Event
    Data        interface{}
}

func getLoginResultSucess(userData interface{}) *LoginResult {
    return &LoginResult{LoggedIn: true, Data: userData}
}

func getLoginResultFailure(errorEvent Event, data interface{}) *LoginResult {
    return &LoginResult{LoggedIn: false, ErrorEvent: errorEvent, Data: data}
}
