package ffi

// #include "log.h"
import "C"

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/hsfzxjy/sdxtra/internal/log"
)

func _() {
	// Ensure that the log levels are consistent between C and Go.
	var _ = "0"[C.SD_LOG_DEBUG-log.Debug]
	var _ = "0"[C.SD_LOG_INFO-log.Info]
	var _ = "0"[C.SD_LOG_WARN-log.Warn]
	var _ = "0"[C.SD_LOG_ERROR-log.Error]
}

const endSentinel = ^uintptr(0)

func init() {
	C.sd_set_log_callback((*[0]byte)(C.handleLog), nil)
}

type LogEntry = log.Entry

type EndWaiter struct {
	Signify func()
	Done    <-chan struct{}
}

type logRoute struct {
	owner  any
	doneCh chan<- struct{}
	outCh  chan<- LogEntry
}

type logHandler struct {
	mu     sync.RWMutex
	routes map[C.pthread_t]*logRoute
}

var logHandlerInstance = logHandler{
	routes: make(map[C.pthread_t]*logRoute, 16),
}
var globalLog = make(chan LogEntry, 16)

func GlobalLog() chan LogEntry {
	return globalLog
}

func (lh *logHandler) addRoute(threadId C.pthread_t, owner any, outCh chan<- LogEntry) EndWaiter {
	doneCh := make(chan struct{})
	route := &logRoute{
		owner:  owner,
		doneCh: doneCh,
		outCh:  outCh,
	}
	lh.mu.Lock()
	lh.routes[threadId] = route
	lh.mu.Unlock()
	return EndWaiter{
		Signify: func() {
			C.goHandleLog(0, nil, C.uintptr_t(endSentinel), threadId)
		},
		Done: doneCh,
	}
}

//export goHandleLog
func goHandleLog(level C.sd_log_level_t, text *C.char, data uintptr, threadId C.pthread_t) {
	if data == endSentinel {
		logHandlerInstance.mu.Lock()
		route := logHandlerInstance.routes[threadId]
		delete(logHandlerInstance.routes, threadId)
		logHandlerInstance.mu.Unlock()
		close(route.doneCh)
		return
	}
	logHandlerInstance.mu.RLock()
	defer logHandlerInstance.mu.RUnlock()
	route, ok := logHandlerInstance.routes[threadId]
	if !ok {
		fmt.Printf("WILD [%d] %s", level, C.GoString(text))
		return
	}
	route.outCh <- LogEntry{
		Level:   log.Level(level),
		Message: C.GoString(text),
		Owner:   route.owner,
	}
}

func captureLog(outCh chan<- LogEntry, data any, fn func()) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	endWaiter := logHandlerInstance.addRoute(C.pthread_self(), data, outCh)
	fn()
	endWaiter.Signify()
	<-endWaiter.Done
}

func CaptureLog0[R any](outCh chan<- LogEntry, data any, fn func() R) R {
	var result R
	captureLog(outCh, data, func() {
		result = fn()
	})
	return result
}

func CaptureLog0Err[R any](outCh chan<- LogEntry, data any, fn func() (R, error)) (R, error) {
	var result R
	var err error
	captureLog(outCh, data, func() {
		result, err = fn()
	})
	return result, err
}

func CaptureLog1[A any, R any](outCh chan<- LogEntry, data any, fn func(A) R, a A) R {
	var result R
	captureLog(outCh, data, func() {
		result = fn(a)
	})
	return result
}

func CaptureLog1Err[A any, R any](outCh chan<- LogEntry, data any, fn func(A) (R, error), a A) (R, error) {
	var result R
	var err error
	captureLog(outCh, data, func() {
		result, err = fn(a)
	})
	return result, err
}
