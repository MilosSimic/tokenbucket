package main

import "time"

type Message struct {
	Kind  int
	Taken int
}

type BucketOptions struct {
	Capacity int
	Tokens   int
	Rate     time.Duration
}

type TokenBucket struct {
	Capacity int
	Tokens   int
	Rate     time.Duration
	Notify   chan *Message
}
