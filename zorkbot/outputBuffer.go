package zorkbot

import (
	"bytes"
	"strings"
	"sync"
	"time"
)

type outputBuffer struct {
	msg chan string
	mu  sync.Mutex
	buf bytes.Buffer
}

func newOutputBuffer() *outputBuffer {
	ob := &outputBuffer{
		msg: make(chan string, 1),
	}
	go ob.send()
	return ob
}

func (ob *outputBuffer) WriteString(s string) {
	ob.mu.Lock()
	ob.buf.WriteString(s)
	ob.mu.Unlock()
}

func (ob *outputBuffer) send() {
	tick := time.Tick(time.Second)
	for {
		<-tick

		ob.mu.Lock()
		line := ob.buf.String()
		ob.buf.Reset()
		ob.mu.Unlock()
		if line == "" {
			continue
		}
		line = strings.Replace(line, "\n\n", " - ", -1)
		line = strings.Replace(line, "\n", " - ", -1)
		line = strings.TrimSuffix(line, " - >")
		ob.msg <- line
	}
}
