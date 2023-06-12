package thread

import (
	"errors"
	"sync"
	"sync/atomic"
)

const (
	MinMessageBufferSize uint = 10
)

type Message[T any] struct {
	Err error
	Msg T
}

type WorkFunc[T any] func(errs chan<- error, tx chan<- Message[T], stop <-chan bool)

func NewMessage[T any](msg T, err error) Message[T] {
	return Message[T]{Msg: msg, Err: err}
}

func NewTxChannel[T any](size uint) chan Message[T] {
	if size < MinMessageBufferSize {
		size = MinMessageBufferSize
	}

	return make(chan Message[T], size)
}

type Status int

func (s Status) String() string {
	switch s {
	case StatusInit:
		return "init"
	case StatusStopped:
		return "stopped"
	case StatusStarting:
		return "starting"
	case StatusRunning:
		return "running"
	case StatusReloading:
		return "reloading"
	case StatusStopping:
		return "stopping"
	}

	return ""
}

type AtomicStatus struct {
	atomic.Uint32
}

// String thread safe implementation of the Stringer interface
func (s *AtomicStatus) String() string {
	return Status(s.Load()).String()
}

const (
	StatusInit     Status = 0
	StatusStarting Status = 1 << (iota - 1)
	StatusRunning
	StatusReloading
	StatusStopping
	StatusStopped
)

type statusChange struct {
	New Status
	Old Status
}

type Thread[T any] struct {
	tx           chan Message[T]
	errs         chan error
	statusChan   chan statusChange
	txBufferSize uint
	workFn       WorkFunc[T]
	status       *AtomicStatus
	wg           *sync.WaitGroup
}

func New[T Thread[M], M any](workFunc WorkFunc[M]) *T {
	t := &T{
		statusChan:   make(chan statusChange, 1),
		tx:           nil,
		errs:         nil,
		txBufferSize: 0, // 0 causes MinMessageBufferSize to be used
		workFn:       workFunc,
		status:       &AtomicStatus{},
		wg:           &sync.WaitGroup{},
	}

	return t
}

func (t *Thread[T]) SetBufferSize(size uint) {
	t.txBufferSize = size
}

// GetStatus thread safe method to get the current status of the thread
func (t *Thread[T]) GetStatus() Status {
	return Status(t.status.Load())
}

// setStatus thread safe internal method for setting the thread status
func (t *Thread[T]) setStatus(status Status) {
	old := Status(t.status.Swap(uint32(status)))

	t.statusChan <- statusChange{
		New: status,
		Old: old,
	}
}

func (t *Thread[T]) Reload() error {
	status := t.GetStatus()

	if status != StatusRunning {
		return errors.New("cannot reload when thread status is: " + status.String())
	}

	t.setStatus(StatusReloading)

	return nil
}

func (t *Thread[T]) Start() error {
	status := t.GetStatus()

	if status == StatusStarting {
		return nil
	}

	if status != StatusStopped && status != StatusInit {
		return errors.New("cannot start when thread status is " + status.String())
	}

	t.setStatus(StatusStarting)

	t.tx = NewTxChannel[T](t.txBufferSize)
	t.errs = make(chan error, 1)

	// run the thread
	t.run()

	return nil
}

func (t *Thread[T]) Stop() error {
	status := t.GetStatus()

	if status == StatusStopped || status == StatusInit {
		// already stopped
		return nil
	}

	if status != StatusRunning && status != StatusReloading {
		return errors.New("cannot stop when thread status is: " + status.String())
	}

	if status == StatusRunning {
		t.setStatus(StatusStopping)
	} else { // case when status == StatusReloading
		// wait for status to change
		s, ok := <-t.statusChan

		if !ok {
			// channel closed.  this is a hard failure
			return errors.New("thread failed to stop and is unrecoverable. reason: unknown")
		}

		if s.New != StatusRunning {
			// reload finished, but no longer in the correct status to run Stop()
			return errors.New("thread failed to stop while waiting for reload. status changed: " + s.Old.String() + " --> " + s.New.String())
		}

		t.setStatus(StatusStopping)
	}

	return nil
}

func (t *Thread[T]) Wait() {
	t.wg.Wait()
}

func (t *Thread[T]) Subscribe() <-chan Message[T] {
	return t.tx
}

func (t *Thread[T]) SubscribeErrs() <-chan error {
	return t.errs
}

// run starts t.workFn in a new go routine and starts go routines for thread communication
func (t *Thread[T]) run() {
	// create new channels for communicating with t.workFn go routine
	errChan := make(chan error, 1)
	txChan := NewTxChannel[T](t.txBufferSize)
	stopChan := make(chan bool, 1)

	// bridge errors from the worker with the calling thread
	t.wg.Add(1)
	go func() {
		defer func() {
			close(t.errs)
			t.wg.Done()
		}()

		for {
			err, ok := <-errChan

			if !ok {
				// channel closed. exit
				return
			}

			t.errs <- err
		}
	}()

	// bridge messages from the worker with the calling thread
	t.wg.Add(1)
	go func() {
		defer func() {
			close(t.tx)
			t.wg.Done()
		}()

		for {
			msg, ok := <-txChan

			if !ok {
				// channel closed. exit
				return
			}

			t.tx <- msg
		}
	}()

	// bridge stop and reload requests from the calling thread with the worker
	t.wg.Add(1)
	go func() {
		defer func() {
			close(stopChan)
			t.wg.Done()
		}()

		for {
			select {
			case s, ok := <-t.statusChan:
				if !ok {
					// channel closed. this SHOULD NEVER happen. abort immediately
					t.errs <- errors.New("aborting thread. unrecoverable")
					return
				}

				switch s.New {
				case StatusStopping:
					// tell the worker to stop (the worker will still get signaled when the channel closes)
					if len(stopChan) < cap(stopChan) {
						stopChan <- true
					}

					t.setStatus(StatusStopped)

					return
				case StatusReloading:
					// tell the worker to reload (only reason this might happen is a backlog of reload messages, no need to block on another one)
					if len(stopChan) < cap(stopChan) {
						stopChan <- false
					}

					t.setStatus(StatusRunning)
				case StatusStopped:
					// thread has stopped (shouldn't get here, but just in case)
					return
				}
			}
		}
	}()

	// run the worker
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.workFn(errChan, txChan, stopChan)
	}()

	// startup complete. set status to running
	t.setStatus(StatusRunning)
}
