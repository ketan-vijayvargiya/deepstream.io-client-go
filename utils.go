package deepstreamio

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
