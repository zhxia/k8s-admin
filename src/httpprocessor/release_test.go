package httpprocessor

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second)
	ticker2 := time.NewTicker(time.Microsecond * 500)
	for {
		select {
		case <-ticker.C:
			for {
				time.Sleep(time.Second * 3600)
				fmt.Println("ok")
			}
		case <-ticker2.C:
			fmt.Println("check")
		}
	}
}

func TestTest(t *testing.T) {
	r, err := strconv.Atoi("34RT")
	if err != nil {
		println("@@@", err)
	}
	println("aaaa", r)
	var n int64
	fmt.Println(n)
	fmt.Println(uuid.NewV4().String())
}

func TestFor(t *testing.T) {
	var i, n, total, step int32
	i = 5
	n = 0
	total = 10
	step = 7
	for ; i <= total; i++ {
		n++
		if n%step == 0 {
			fmt.Println("i:", i)
		}
	}
	if n%step != 0 {
		println("total:", total)
	}

}

func TestReleaseScale(t *testing.T) {
	channel := make(chan int32)
	done := make(chan bool)
	go func() {
		i := int32(0)
		for {
			if i == 10 {
				break
			}
			time.Sleep(1 * time.Second)
			channel <- i
			i += 1
		}
		close(channel)
	}()
	go func() {
		for {
			o, ok := <-channel
			fmt.Println(ok)
			if !ok {
				break
			}
			fmt.Println(o)
		}
		done <- true
	}()
	<-done
	fmt.Println("over")
}
