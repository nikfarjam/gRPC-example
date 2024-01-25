package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	multiConsumerProducer(100, 50)
}

func singleConsumerProducer() {
	ch := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go simpleProducer(ch, &wg)
	wg.Add(1)
	go simpleConsumer(ch, &wg)
	wg.Wait()
	close(ch)
}

func simpleProducer(ch chan string, wg *sync.WaitGroup) {
	ch <- "Producer sent message"
	wg.Done()
}

func simpleConsumer(ch chan string, wg *sync.WaitGroup) {
	select {
	case msg := <-ch:
		fmt.Printf("Consumer Received: %s\n", msg)
	case <-time.After(1 * time.Second):
		wg.Done()
	}
	wg.Done()
}

func multiConsumerProducer(producerSize, consumerSize int) {
	ch := make(chan string)
	var wg sync.WaitGroup

	// Start multiple producers
	for i := 0; i < producerSize; i++ {
		wg.Add(1)
		go producer(i, ch, &wg)
	}

	// Start multiple consumers
	for i := 0; i < consumerSize; i++ {
		wg.Add(1)
		go consumer(i, ch, &wg)
	}

	wg.Wait()
	close(ch)
}

func producer(index int, ch chan string, wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		ch <- fmt.Sprintf("Producer %v send %v", index, i)
	}
	wg.Done()
}

func consumer(index int, ch chan string, wg *sync.WaitGroup) {
	done := false
	for !done {
		select {
		case msg, ok := <-ch:
			if !ok {
				done = true
			}
			fmt.Printf("Consumer %v Received: %s\n", index, msg)
		case <-time.After(1 * time.Second):
			wg.Done()
		}

	}
	wg.Done()
}
