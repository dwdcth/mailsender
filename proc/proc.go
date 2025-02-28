package proc

import (
	"log"
	"sync/atomic"
	"time"
)

// Counter represents a thread-safe counter with QPS calculation
type Counter struct {
	name      string
	count     atomic.Int64
	lastCount int64
	lastTime  time.Time
	qps       float64
}

func NewCounter(name string) *Counter {
	return &Counter{
		name:     name,
		lastTime: time.Now(),
	}
}

func (c *Counter) Incr() {
	c.count.Add(1)
}

func (c *Counter) UpdateQps() {
	now := time.Now()
	current := c.count.Load()
	duration := now.Sub(c.lastTime).Seconds()
	c.qps = float64(current-c.lastCount) / duration
	c.lastCount = current
	c.lastTime = now
}

func (c *Counter) String() string {
	return c.name
}

// counter instances
var (
	HttpRequestCnt = NewCounter("HttpRequestCnt")
	MailSendCnt    = NewCounter("MailSendCnt")
	MailSendOkCnt  = NewCounter("MailSendOkCnt")
	MailSendErrCnt = NewCounter("MailSendErrCnt")
)

func Start() {
	// Update QPS every second
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for range ticker.C {
			HttpRequestCnt.UpdateQps()
			MailSendCnt.UpdateQps()
			MailSendOkCnt.UpdateQps()
			MailSendErrCnt.UpdateQps()
		}
	}()
	log.Println("proc.Start, ok")
}

func GetAll() []interface{} {
	ret := make([]interface{}, 0)
	ret = append(ret, HttpRequestCnt)
	ret = append(ret, MailSendCnt)
	ret = append(ret, MailSendOkCnt)
	ret = append(ret, MailSendErrCnt)
	return ret
}
