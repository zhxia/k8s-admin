package api

import (
	"fmt"
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	logApi := NewLogApi("default", "")
	out, err := logApi.ExecRemoteCmd("default", "sample-g4wdr", "nginx", []string{"/bin/bash", "-c", "tail -n2 /tmp/aa1.log"})
	if err != nil {
		log.Fatal("error:", err)
	} else {
		fmt.Println(string(out))
	}
}
