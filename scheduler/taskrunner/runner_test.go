package taskrunner

import (
	"errors"
	"log"
	"testing"
	"time"
)

func TestRunner(t *testing.T)  {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i;
			log.Printf("Dispatch sent: %d", i)
		}
		return nil
	}
	e := func(dc dataChan) error {
		forloop:
			for {
				select {
				case d :=<- dc:
					log.Printf("Executor received: %v", d)
				default:
					break forloop
				}
			}
		return errors.New("Executor")
	}
	runner := NewRunner(30, false, d, e)
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}
