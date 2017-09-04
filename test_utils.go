package deepstreamio

import "testing"

func ValidateStringsEqual(t *testing.T, actual, expected string)  {
    if actual != expected {
        t.Errorf("Actual: %s didn't match Expected: %s", actual, expected)
    }
}
