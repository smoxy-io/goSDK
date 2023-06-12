package thread

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

var reloadFuncCallCount atomic.Uint32

func init() {
	reloadFuncCallCount.Store(0)
}

func noopWorkFunc[T any](errs chan<- error, tx chan<- Message[T], stop <-chan bool) {}
func defWorkFunc[T any](errs chan<- error, tx chan<- Message[T], stop <-chan bool) {
	defer func() {
		close(errs)
		close(tx)
	}()

	// loop until stopped
	for {
		// the stop channel will receive a true message when the work func should stop
		s, ok := <-stop

		if !ok || s {
			// stop
			break
		}

		if !s {
			// reload
			reloadFuncCallCount.Add(1)
		}
	}
}

var simpleWorkFunc WorkFunc[bool] = func(errs chan<- error, tx chan<- Message[bool], stop <-chan bool) {
	defer func() {
		close(tx)
		close(errs)
	}()

	startTime := 10 * time.Millisecond
	startTimer := time.NewTimer(startTime)
	stopTimer := time.NewTimer(startTime)

	// wait a bit to simulate actual startup work being done
	<-startTimer.C

	// loop until stopped
	for {
		// the stop channel will receive a true message when the work func should stop
		s, ok := <-stop

		if !ok || s {
			// stop
			break
		}

		if !s {
			// reload
			reloadFuncCallCount.Add(1)
		}
	}

	// wait a bit to simulate actual stop work being done
	<-stopTimer.C
}

var genWorkFunc WorkFunc[int] = func(errs chan<- error, tx chan<- Message[int], stop <-chan bool) {
	defer func() {
		close(errs)
		close(tx)
	}()

	ticker := time.NewTicker(1 * time.Millisecond)
	i := 0

	// loop until stopped
	for {
		select {
		// the stop channel will receive a true message when the work func should stop
		case s, ok := <-stop:
			if !ok || s {
				// stop
				break
			}

			if !s {
				// reload
				reloadFuncCallCount.Add(1)
			}
		case <-ticker.C:
			i++
			errs <- errors.New(strconv.Itoa(i))
			tx <- NewMessage(i, nil)
		}
	}
}

func TestConstants(t *testing.T) {
	assert.Equal(t, uint(10), MinMessageBufferSize)
	assert.Equal(t, Status(0), StatusInit)
	assert.Equal(t, Status(1), StatusStarting)
	assert.Equal(t, Status(2), StatusRunning)
	assert.Equal(t, Status(4), StatusReloading)
	assert.Equal(t, Status(8), StatusStopping)
	assert.Equal(t, Status(16), StatusStopped)
}

func TestNew(t *testing.T) {
	_ = New(noopWorkFunc[int])
	_ = New(noopWorkFunc[string])
}

func TestNewMessage(t *testing.T) {
	msg1 := NewMessage(2, nil)
	msg2 := NewMessage("foo", errors.New("bar"))

	if msg1.Msg != 2 || msg1.Err != nil {
		t.Errorf("msg1 invalid.  got: %v", msg1)
	}

	if msg2.Msg != "foo" || msg2.Err.Error() != "bar" {
		t.Errorf("msg2 invalid.  got: %v", msg1)
	}
}

func TestNewTxChannel(t *testing.T) {
	c1 := NewTxChannel[int](0)
	c2 := NewTxChannel[int](MinMessageBufferSize)
	c3 := NewTxChannel[int](MinMessageBufferSize + 1)
	c4 := NewTxChannel[string](0)

	assert.IsType(t, make(chan Message[int]), c1)
	assert.IsType(t, make(chan Message[int]), c2)
	assert.IsType(t, make(chan Message[int]), c3)
	assert.IsType(t, make(chan Message[string]), c4)

	if cap(c1) != int(MinMessageBufferSize) {
		t.Errorf("invalid tx buffer size. wanted: %v, got %v", MinMessageBufferSize, cap(c1))
	}

	if cap(c2) != int(MinMessageBufferSize) {
		t.Errorf("invalid tx buffer size. wanted: %v, got %v", MinMessageBufferSize, cap(c2))
	}

	if cap(c3) != int(MinMessageBufferSize+1) {
		t.Errorf("invalid tx buffer size. wanted: %v, got %v", MinMessageBufferSize+1, cap(c3))
	}

	if cap(c4) != int(MinMessageBufferSize) {
		t.Errorf("invalid tx buffer size. wanted: %v, got %v", MinMessageBufferSize, cap(c4))
	}
}

func TestThread_SetBufferSize(t *testing.T) {
	thread := New(noopWorkFunc[int])

	assert.Equal(t, thread.txBufferSize, uint(0))

	thread.SetBufferSize(MinMessageBufferSize)

	assert.Equal(t, thread.txBufferSize, MinMessageBufferSize)

	thread.SetBufferSize(MinMessageBufferSize - 1)

	assert.Equal(t, thread.txBufferSize, MinMessageBufferSize-1)

	thread.SetBufferSize(MinMessageBufferSize + 1)

	assert.Equal(t, thread.txBufferSize, MinMessageBufferSize+1)
}

func TestThread_Start(t *testing.T) {
	thread := New(defWorkFunc[int])

	if err := thread.Start(); err != nil {
		t.Errorf("Thread[T].Start() returned error: %v", err)
	}

	status := thread.GetStatus()
	assert.Equal(t, StatusRunning, status)

	_ = thread.Stop()
}

func TestThread_Stop(t *testing.T) {
	thread := New(simpleWorkFunc)

	if err := thread.Start(); err != nil {
		t.Errorf("Thread[T].Start() returned error: %v", err)
		return
	}

	status := thread.GetStatus()
	assert.Equal(t, StatusRunning, status)

	// simpleWorkFunc takes a minimum of 10ms to "start"
	time.Sleep(15 * time.Millisecond)

	if err := thread.Stop(); err != nil {
		t.Errorf("thread[T].Stop() returned error: %v", err)
	}

	status = thread.GetStatus()
	assert.Equal(t, StatusStopping, status)

	// wait a bit and confirm the thread has stopped
	// simpleWorkFunc takes a minimum of 10ms to stop
	time.Sleep(15 * time.Millisecond)

	status = thread.GetStatus()
	assert.Equal(t, StatusStopped, status)
}

func TestThread_Reload(t *testing.T) {
	thread := New(simpleWorkFunc)

	if err := thread.Start(); err != nil {
		t.Errorf("Thread[T].Start() returned error: %v", err)
		return
	}

	status := thread.GetStatus()
	assert.Equal(t, StatusRunning, status)

	// simpleWorkFunc takes a minimum of 10ms to "start"
	time.Sleep(15 * time.Millisecond)

	err := thread.Reload()

	if err != nil {
		t.Errorf("thread[T].Reload() returned error: %v", err)
	}

	status = thread.GetStatus()
	assert.Equal(t, StatusReloading, status)

	// wait for reload to finish
	time.Sleep(5 * time.Millisecond)

	status = thread.GetStatus()
	assert.Equal(t, StatusRunning, status)

	assert.Equal(t, uint32(1), reloadFuncCallCount.Load())

	if err := thread.Stop(); err != nil {
		t.Errorf("thread[T].Stop() returned error: %v", err)
	}

	// wait a bit and confirm the thread has stopped
	// simpleWorkFunc takes a minimum of 10ms to stop
	time.Sleep(15 * time.Millisecond)

	status = thread.GetStatus()
	assert.Equal(t, StatusStopped, status)

	reloadFuncCallCount.Swap(0)
}

func TestThread_Subscribe(t *testing.T) {
	thread := New(genWorkFunc)

	if err := thread.Start(); err != nil {
		t.Errorf("Threat[T].Start() returned error: %v", err)
	}

	msgs := thread.Subscribe()

	m, ok := <-msgs

	assert.Equal(t, true, ok)
	assert.Equal(t, NewMessage(1, nil), m)

	_ = thread.Stop()
}

func TestThread_SubscribeErrs(t *testing.T) {
	thread := New(genWorkFunc)

	if err := thread.Start(); err != nil {
		t.Errorf("Threat[T].Start() returned error: %v", err)
	}

	errs := thread.SubscribeErrs()

	e, ok := <-errs

	assert.Equal(t, true, ok)
	assert.Equal(t, e.Error(), "1")

	_ = thread.Stop()
}

func TestThread_Wait(t *testing.T) {
	thread := New(simpleWorkFunc)

	if err := thread.Start(); err != nil {
		t.Errorf("Thread[T].Start() returned error: %v", err)
	}

	if err := thread.Stop(); err != nil {
		t.Errorf("Thread[T].Stop() returned error: %v", err)
	}

	thread.Wait()

	status := thread.GetStatus()
	assert.Equal(t, StatusStopped, status)
}
