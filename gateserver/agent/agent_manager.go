package agent

import (
	"fmt"
	"nuanv3/shared/timer"
	"sync"
	"time"
)

type AgentManager struct {
}

var (
	RequestCount     int
	LastRequestCount int
	CurrentQPS       float64
	MaxQPS           float64
	qpsLock          sync.RWMutex
)

const (
	GLOBAL_TIMER = 5
)

func InitAgentManager() {
	RequestCount = 0
	LastRequestCount = 0
	CurrentQPS = 0
	MaxQPS = 0

	go RequestWatch()
}

func InrcRequestCount() {
	qpsLock.Lock()
	defer qpsLock.Unlock()

	RequestCount++
}

func SetRequestCount(num int) {
	qpsLock.Lock()
	defer qpsLock.Unlock()

	RequestCount = num
}

func RequestWatch() {
	globalTimer := make(chan int32, 1)
	timer.Add(-2, time.Now().Unix()+GLOBAL_TIMER, globalTimer)

	for {
		select {
		case <-globalTimer:
			CurrentQPS = float64(RequestCount-LastRequestCount) / float64(GLOBAL_TIMER)
			if CurrentQPS > MaxQPS {
				MaxQPS = CurrentQPS
			}
			LastRequestCount = RequestCount

			fmt.Println("CurrentQPS", CurrentQPS, "MaxQPS", MaxQPS)
			timer.Add(-2, time.Now().Unix()+GLOBAL_TIMER, globalTimer)
		}
	}
}
