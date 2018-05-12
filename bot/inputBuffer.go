package zorkbot

import (
	"bytes"
	"sync"
)

type inputBuffer struct {
	msg chan bool
	mu  sync.Mutex
	buf bytes.Buffer
}

func newInputBuffer() *inputBuffer {
	return &inputBuffer{msg: make(chan bool, 1)}
}

func (b *inputBuffer) WriteString(s string) {
	if len(s) == 0 {
		return
	}
	b.mu.Lock()
	b.buf.WriteString(s)
	b.mu.Unlock()
	b.msg <- true
}

func (b *inputBuffer) ReadRune() (rune, int, error) {
	b.mu.Lock()
	len := b.buf.Len()
	b.mu.Unlock()
	if len == 0 {
		<-b.msg
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.ReadRune()
}
