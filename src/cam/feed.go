package main

import (
	"cam/ipcam"
	"log"
	"sync"
	"time"
)

type feed struct {
	id   string
	src  *ipcam.Source
	buf  []byte
	last time.Time
	subs map[chan struct{}]bool
	lock sync.Mutex
}

func (f *feed) run() {
	f.subs = make(map[chan struct{}]bool)
	f.src.Run()
}

func (f *feed) sub() (c chan struct{}) {
	c = make(chan struct{})

	f.lock.Lock()
	defer f.lock.Unlock()

	f.subs[c] = true

	return c
}

func (f *feed) unsub(c chan struct{}) {
	f.lock.Lock()
	defer f.lock.Unlock()

	delete(f.subs, c)
}

func (f *feed) getFrame() (buf []byte, ok bool) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.buf != nil {
		return f.buf, true
	}

	return nil, false
}

func (f *feed) sink(id string, buf []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()

	if f.buf == nil {
		log.Println("connected to", id)
	}

	f.buf = append([]byte{}, buf...)
	f.last = time.Now()

	// notify subscribers
	//
	e := struct{}{}

	for sub, _ := range f.subs {
		select {
		case sub <- e:
			// OK
		default:
			// NOP
		}
	}
}
