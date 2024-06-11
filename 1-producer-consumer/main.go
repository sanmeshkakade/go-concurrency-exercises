//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, tweets chan *Tweet, wg *sync.WaitGroup) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {

			// stream has ended
			close(tweets)
			wg.Done()
			return
		}

		tweets <- tweet
	}
}

func consumer(tweets <-chan *Tweet, wg *sync.WaitGroup) {
	for t := range tweets {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang")
		}
	}
	wg.Done()
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	wg := &sync.WaitGroup{}

	// Producer
	tweets := make(chan *Tweet)
	wg.Add(1)
	// run producer on goroutine
	// return channel to communicate
	// produced value
	// so that we can unblock
	// consumer
	go producer(stream, tweets, wg)

	// Consumer
	wg.Add(1)
	// run consumer goroutine concurrently
	// to optimise the consumption
	go consumer(tweets, wg)

	// wait for both goroutines to finish
	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
