package engine

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type EngineState int32

const (
	Idle EngineState = iota
	Running
	Paused
	Stopped
)

type ScriptEngine struct {
	state        int32
	mu           sync.Mutex
	cond         *sync.Cond
	onResumeFunc func()
	onPauseFunc  func()
	onStopFunc   func()
}

func NewEngine() *ScriptEngine {
	e := &ScriptEngine{}
	e.cond = sync.NewCond(&e.mu)
	return e
}

func (e *ScriptEngine) setState(s EngineState) { atomic.StoreInt32(&e.state, int32(s)) }
func (e *ScriptEngine) GetState() EngineState  { return EngineState(atomic.LoadInt32(&e.state)) }

func (e *ScriptEngine) SetOnStop(fn func()) {
	e.onStopFunc = fn
}

func (e *ScriptEngine) Stop() {
	if e.GetState() == Running || e.GetState() == Paused {
		fmt.Println("Stopped")
		e.setState(Stopped)
		if e.onStopFunc != nil {
			e.onStopFunc()
		}
	}
}

func (e *ScriptEngine) SetOnPause(fn func()) {
	e.onPauseFunc = fn
}

func (e *ScriptEngine) Pause() {
	if e.GetState() == Running {
		fmt.Println("Paused")
		e.setState(Paused)
		if e.onPauseFunc != nil {
			e.onPauseFunc()
		}
	}
}

func (e *ScriptEngine) SetOnResume(fn func()) {
	e.onResumeFunc = fn
}

func (e *ScriptEngine) Resume() {
	if e.GetState() == Paused {
		fmt.Println("Resumed")
		e.setState(Running)
		e.cond.Broadcast()
		if e.onResumeFunc != nil {
			e.onResumeFunc()
		}
	}
}

func (e *ScriptEngine) TogglePause() {
	if e.GetState() == Paused {
		e.Resume()
	} else if e.GetState() == Running {
		e.Pause()
	}
}

type stopSignal struct{}

func (e *ScriptEngine) handleState() {
	e.mu.Lock()
	for e.GetState() == Paused {
		e.cond.Wait()
	}
	e.mu.Unlock()

	if e.GetState() == Stopped {
		panic(stopSignal{})
	}
}

func (e *ScriptEngine) RunFuncWithReturn(fn func() interface{}) interface{} {
	e.GetState()
	return fn()
}

func (e *ScriptEngine) RunFuncNoReturn(fn func()) {
	e.GetState()
	fn()
}

func (e *ScriptEngine) Start(fn func()) {
	if e.GetState() == Idle || e.GetState() == Stopped {
		fmt.Println("Starting macro")
		e.setState(Running)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					if _, ok := r.(stopSignal); ok {
						fmt.Println("✅ Macro halted")
						e.setState(Stopped)
						return
					}
					panic(r)
				}
				e.setState(Idle)
			}()

			fn()
		}()
	}
}

func (e *ScriptEngine) Sleep(d time.Duration) {
	start := time.Now()
	for {
		e.handleState()
		if time.Since(start) >= d {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
