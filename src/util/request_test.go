package util

import (
    "fmt"
    "testing"
)

func TestShow(t *testing.T) {
    mp := map[string]string{
        "name": "zhxia",
    }
    if name, ok := mp["name1"]; !ok {
        fmt.Println("not exist")
    } else {
        fmt.Println(name)
    }
}
