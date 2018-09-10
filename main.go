package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	opt := BucketOptions{10, 0, time.Second}
	bucket := New(ctx, &opt)

	timer := time.NewTimer(time.Second * 20) // Simulate for 20 seconds
	for {
		select {
		case <-timer.C:
			fmt.Println("Drain ", bucket.Drain())
			cancel()
			return
		case msg := <-bucket.Notify:
			switch msg.Kind {
			case BUCKET_FULL:
				fmt.Println("GET: ", bucket.TakeAll())
			case TAKEN_ALL:
				fmt.Println("Taken all")
			case TAKEN_N:
				fmt.Println("Taken n")
			case DRAIN:
				fmt.Println("Drain")
			}
		}
	}
}
