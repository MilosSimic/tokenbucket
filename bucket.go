package main

import (
	"context"
	"log"
	"time"
)

func New(ctx context.Context, opt *BucketOptions) *TokenBucket {
	b := &TokenBucket{
		Capacity: opt.Capacity,
		Tokens:   opt.Tokens,
		Rate:     opt.Rate,
		Notify:   make(chan *Message),
	}

	go func(bucket *TokenBucket) {
		timer := time.NewTicker(bucket.Rate)
		for {
			select {
			case <-timer.C:
				if bucket.Tokens < bucket.Capacity {
					bucket.Tokens++
					log.Print("Token added")
					bucket.notify(TOKEN_ADDED, 1, "Token added")
				} else if bucket.Tokens == bucket.Capacity {
					bucket.notify(BUCKET_FULL, bucket.Capacity, "Bucket full")
				} else {
					log.Print("Token not added, bucket full")
				}

			case <-ctx.Done():
				log.Print(ctx.Err())
				timer.Stop()
				close(bucket.Notify)
				return
			}
		}
	}(b)

	return b
}

func (b *TokenBucket) notify(kind, amount int, msg string) {
	b.Notify <- &Message{kind, amount}
	log.Print(msg)
}

func (b *TokenBucket) TakeAll() bool {
	if b.Tokens == 0 {
		return false
	}

	if b.Tokens == b.Capacity {
		b.Tokens = 0
		go b.notify(TAKEN_ALL, b.Capacity, "Tokens Taken")
		return true
	}

	return false
}

func (b *TokenBucket) Take(n int) bool {
	if b.Tokens == 0 || b.Tokens < n {
		return false
	}

	if b.Tokens <= b.Capacity {
		curr := b.Tokens - n
		b.Tokens = curr

		go b.notify(TAKEN_ALL, n, "Tokens Taken")
		return true
	}

	return false
}

func (b *TokenBucket) Drain() bool {
	if b.Tokens == 0 {
		return false
	}

	curr := b.Tokens
	b.Tokens = 0
	go b.notify(TAKEN_ALL, curr, "Tokens Drained")
	return true
}
