package deepstreamio

import "time"

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func getIntOrDefault(x, d int) int {
    if x == 0 {         // i.e. nil
        return d
    }
    return x
}

func getStringOrDefault(x, d string) string {
    if len(x) == 0 {    // i.e. nil
        return d
    }
    return x
}

func getDurationMillis(t int) time.Duration {
    return time.Duration(t) * time.Millisecond
}

func scheduleFuncAfterMillis(f func(), millis int) *time.Timer {
    d := getDurationMillis(millis)
    timer := time.NewTimer(d)
    go func() {
        <- timer.C
        f()
    }()

    return timer
}
