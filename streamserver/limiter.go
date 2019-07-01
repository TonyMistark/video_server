package main

import "github.com/gpmgo/gopm/modules/log"

/*
bucket: token[10]
1. request bucket token count -1
2. finish request bucket token count +1

channel: share channel instead of share memory
 */

type ConnLimiter struct {
	concurrentConn int
	bucket chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn:cc,
		bucket:make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {
	if len(cl.bucket) >= cl.concurrentConn{
		log.Print(0, "Reached the rate limitation.")
		return false
	}
	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {
	c := <- cl.bucket
	log.Print(0, "New connection coming: %d", c)
}