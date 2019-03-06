package xsync

import (
	"sync"
	"sync/atomic"
)

// MustOnce is an object that will perform exactly one action until the action succeed.
// All failed run for the action will have the opportunity to run the action again.
type MustOnce struct {
	m    sync.Mutex
	done int32
}

// Do calls the function f if and only if f is not being called successfully
// for this instance of MustOnce. In other words, given
// 	var mu MustOnce
// if mu.Do(f) is called multiple times, only until the first invocation to f that is
// succeeded will be executed. After the first successful invocation of f
// Do will not invoke f even if f has a different value in each invocation.
// A new instance of MustOnce is required for different function to execute successfully.
//
// Do is intended for initialization that must be run exactly once successfully. f must
// return an error in case of any failed invocation. it may be necessary to use
// a function literal to capture the arguments to a function to be invoked by Do.
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// If f panics, Do considers it to have returned successfully; future calls of Do return
// without calling f.
//
// This implementation is an adaption from https://golang.org/pkg/sync/#Once. It only varies
// by the successful invocation of f.
//
// This implementation is an fork of
// https://github.com/appscode/go/blob/master/sync with changing the Signatures.
func (o *MustOnce) Do(f func() error) {
	if atomic.LoadInt32(&o.done) == 1 {
		return
	}

	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreInt32(&o.done, 1)
		err := f()
		if err != nil {
			// action have failed, unset the done value so it can be rerun
			atomic.StoreInt32(&o.done, -1)
		}
	}
}
