package api

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	mp := map[string]string{"name": "zhxia"}
	if n, ok := mp["name"]; ok {
		fmt.Println("exist:", n)
	} else {
		fmt.Println("not exist!")
	}
	str := "zhxia@aaa"
	arr := strings.SplitN(str, "@", 2)
	println(arr[0], arr[1])
	tm := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(tm)
}

func TestChannel(t *testing.T) {
	taskQueue := make(chan int32)
	finishFlag := make(chan bool)
	go func() {
		i := int32(0)
		for {
			taskQueue <- i
			time.Sleep(1 * time.Second)
			i++
			if i == 10 {
				break
			}
		}
		close(taskQueue)
	}()

	go func() {
		for v := range taskQueue {
			fmt.Println(v)
		}
		finishFlag <- true
		close(finishFlag)
	}()
	fmt.Println("waiting...")
	<-finishFlag
	fmt.Println("finish!")
}
