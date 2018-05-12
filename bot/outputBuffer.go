package zorkbot

import (
	"bytes"
	"strings"
	"sync"
	"time"
)

type outputBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func newOutputBuffer() *outputBuffer {
	ob := &outputBuffer{}
	go ob.send()
	return ob
}

func (ob *outputBuffer) WriteString(s string) {
	ob.mu.Lock()
	ob.buf.WriteString(s)
	ob.mu.Unlock()
}

func (ob *outputBuffer) send() {
	tick := time.Tick(2 * time.Second)
	for {
		<-tick

		ob.mu.Lock()
		line := ob.buf.String()
		ob.buf.Reset()
		ob.mu.Unlock()
		if line == "" {
			continue
		}
		line = strings.Replace(line, "\n", " ", -1)
		line = strings.TrimSuffix(line, "  >")
		//send message?
	}
}
