package main

import (
	"context"
	"log"
	"time"
)

func New(ctx context.Context, opt *BucketOptions) *TokenBucket {
	bucket := &TokenBucket{
		Capacity: opt.Capacity,
		Tokens:   opt.Tokens,
		Rate:     opt.Rate,
	}

	go func() {
		timer := time.NewTicker(bucket.Rate)
		for {
			select {
			case <-timer.C:
				if bucket.Tokens < bucket.Capacity {
					bucket.Tokens++
					log.Print("Token added ", time.Now())
				} else if bucket.Tokens == bucket.Capacity {
					bucket.flush(BUCKET_FULL, bucket.Capacity, 0, "Bucket full, flush!")
				}

			case <-ctx.Done():
				log.Print(ctx.Err())
				timer.Stop()
				close(bucket.Notify)
				return
			}
		}
	}()

	return bucket
}

func (b *TokenBucket) flush(kind, amount, tokens int, msg string) {
	b.Tokens = tokens
	b.Notify <- &Message{kind, amount}
	log.Print(msg, " ", time.Now())
}

func (b *TokenBucket) TakeAll() bool {
	if b.Tokens == 0 {
		return false
	}

	if b.Tokens == b.Capacity {
		b.flush(TAKEN_ALL, b.Capacity, 0, "Tokens taken, flush!")
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
		b.flush(TAKEN_N, n, curr, "Taken n, flush!")
		return true
	}

	return false
}

func (b *TokenBucket) Drain() bool {
	if b.Tokens == 0 {
		return false
	}

	curr := b.Tokens
	b.flush(DRAIN, curr, 0, "Bucket drain, flush!")
	return true
}
