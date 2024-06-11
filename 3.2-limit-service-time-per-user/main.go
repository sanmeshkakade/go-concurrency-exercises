//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
	Mutex     sync.Mutex
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	done := make(chan bool)

	go func() {
		process()
		done <- true
	}()

	for {
		select {
		case <-done:
			return true
		case <-time.Tick(time.Second):
			u.Mutex.Lock()
			u.TimeUsed++
			if u.TimeUsed > 10 && !u.IsPremium {
				u.Mutex.Unlock()
				return false
			}
			u.Mutex.Unlock()
		}
	}
}

func main() {
	start := time.Now()
	RunMockServer()
	fmt.Printf("Process took %s\n", time.Since(start))
}
